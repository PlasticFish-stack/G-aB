package model

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProdType struct {
	Global
	Name        string            `gorm:"size:255;not null;unique" json:"name"`
	Description string            `gorm:"size:255" json:"description"`
	Sort        uint              `gorm:"column:type_sort;default:0" json:"sort"`
	ParentId    uint              `gorm:"default:0;not null" json:"parentId"`
	Tax         float64           `gorm:"default:0;not null" json:"tax"`
	Children    []ProdType        `gorm:"-" json:"children,omitempty"`
	Field       []ProdTypeField   `gorm:"foreignKey:TypeId" json:"fields"`
	Formulas    []ProdTypeFormula `gorm:"foreignKey:TypeId" json:"formulas"`
	// Products    []Product            `gorm:"foreignKey:ProductTypeId" json:"formulas"`
}

type ProdTypeField struct {
	Global
	Name     string `gorm:"size:255;not null;" json:"name"`
	NickName string `gorm:"size:255;not null;" json:"nickname"`
	TypeId   uint
}

type ProdTypeFormula struct {
	Global
	Name     string `gorm:"size:255;not null;" json:"name"`
	NickName string `gorm:"size:255;not null;" json:"nickname"`
	Formula  string `gorm:"size:255;not null;" json:"formula"`
	TypeId   uint
}

func SearchTreeProductType(db *gorm.DB) ([]*ProdType, error) {
	var productTypes []ProdType
	if err := db.Order("type_sort").Preload(clause.Associations).Find(&productTypes).Error; err != nil {
		return nil, fmt.Errorf("查询产品类别失败: %v", err)
	}
	treeMapProductTypes := make(map[uint]*ProdType)
	for i := range productTypes {
		productTypes[i].Children = []ProdType{}
		treeMapProductTypes[productTypes[i].Id] = &productTypes[i]
	}
	var treeProductTypes []*ProdType
	for i := range productTypes {
		productType := treeMapProductTypes[productTypes[i].Id]
		if productType.ParentId == 0 {
			treeProductTypes = append(treeProductTypes, productType)
		} else {
			treeMapProductTypes[productType.ParentId].Children = append(treeMapProductTypes[productType.ParentId].Children, *productType)
		}
	}
	for _, parentProductType := range treeProductTypes {
		sort.Slice(parentProductType.Children, func(i, j int) bool {
			return parentProductType.Children[i].Id < parentProductType.Children[j].Id
		})
	}
	return treeProductTypes, nil
}

func (p *ProdType) Search(db *gorm.DB) (*ProdType, error) {
	var prodType ProdType
	if err := db.First(&prodType, p.Id).Error; err != nil {
		return nil, fmt.Errorf("查询产品类别失败: %v", err)
	}
	return &prodType, nil
}

func (p *ProdType) Add(db *gorm.DB) error {
	if err := db.Create(p).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("产品类别名称已存在: %v", err)
		}
		return fmt.Errorf("新建产品类别失败: %v", err)
	}
	return nil
}

func (p *ProdType) Update(db *gorm.DB) error {
	var resultProductType ProdType
	if err := db.Where(p.Id).First(&resultProductType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("未查询到该产品类别: %v", err)
		}
		return fmt.Errorf("查询产品类别失败: %v", err)
	}
	if err := db.Model(&resultProductType).Updates(&p).Error; err != nil {
		return fmt.Errorf("更新产品类别失败,请检查: %v", err)
	}

	if err := db.Unscoped().Model(&resultProductType).Association("Field").Unscoped().Replace(&p.Field); err != nil {
		fmt.Printf("%+v\n", err)
	}

	if err := db.Unscoped().Model(&resultProductType).Association("Formulas").Unscoped().Replace(&p.Formulas); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return nil
}

func (p *ProdType) Delete(db *gorm.DB) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Commit()
		}
	}()
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	p.Name = p.Name + "_is_deleted" + currentTime
	if err := tx.Updates(&p).Error; err != nil {
		return fmt.Errorf("删除产品类别失败,请检查: %v", err)
	}
	if err := tx.Delete(&p).Error; err != nil {
		return fmt.Errorf("删除产品类别失败,请检查: %v", err)
	}
	return nil
}

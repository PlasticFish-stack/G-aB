package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Product struct {
	Global
	ItemNumber     string        `gorm:"size:255;not null;unique;comment:货号" json:"itemName"`
	BrandId        uint          `gorm:"comment:品牌" json:"brandId"`
	Sku            string        `gorm:"size:255;comment:sku" json:"sku"`
	Spu            string        `gorm:"size:255;comment:spu" json:"spu"`
	Quantity       uint64        `gorm:"default:1;not null;comment:数量" json:"quantity"`
	Specifications string        `gorm:"type:text;comment:规格" json:"specifications"`
	Barcode        string        `gorm:"size:255;comment:条形码" json:"barcode"`
	Customscode    string        `gorm:"size:255;comment:海关编码" json:"customscode"`
	Description    string        `gorm:"size:255;comment:描述" json:"description"`
	Color          string        `gorm:"size:255;comment:颜色" json:"color"`
	DwPrice        float64       `gorm:"comment:得物价格" json:"dwPrice"`
	ProductTypeId  uint          `gorm:"comment:产品类型id"`
	ProductCost    []ProductCost `gorm:"foreignKey:ProductID"`
}

type ProductCost struct {
	Global
	ProductID  uint    //product外键
	ExcelLogId uint    //operationlog表格外键
	Cost       float64 //成本价
	CurrencyId uint    //货币id
}

func SearchAllProduct(db *gorm.DB, page int, pageSize int) ([]Product, error) {
	var product []Product
	if err := db.Preload(clause.Associations).Find(&product).Error; err != nil {
		return nil, fmt.Errorf("搜索出错: %v", product)
	}
	return product, nil
}

func (p *Product) Add(db *gorm.DB) error {
	var resultProd Product
	if err := db.Where("item_nubmer = ?", p.ItemNumber).First(&resultProd).Error; err != nil {
		if err := db.Create(p).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return fmt.Errorf("产品名称已存在: %v", err)
			}
			return fmt.Errorf("新建产品失败: %v", err)
		}
	}
	return fmt.Errorf("产品已存在")
}

func (p *Product) Update(db *gorm.DB) error {
	var resultProd Product
	if err := db.Where("item_nubmer = ?", p.ItemNumber).First(&resultProd).Error; err != nil {
		if err := db.Model(&resultProd).Updates(&p).Error; err != nil {
			return fmt.Errorf("更新产品失败,请检查: %v", err)
		}
	}
	return nil
}

func (p *Product) Search(db *gorm.DB) (*Product, error) {
	var prod *Product
	if err := db.First(&prod, p.Id).Error; err != nil {
		return nil, fmt.Errorf("查询错误: %v", err)
	}
	return prod, nil
}

func (p *Product) Delete(db *gorm.DB) error {
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
	p.ItemNumber = p.ItemNumber + "_is_deleted" + currentTime
	if err := tx.Updates(&p).Error; err != nil {
		return fmt.Errorf("删除产品失败,请检查: %v", err)
	}
	if err := tx.Delete(&p).Error; err != nil {
		return fmt.Errorf("删除产品失败,请检查: %v", err)
	}
	return nil
}

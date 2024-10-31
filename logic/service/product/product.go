package product

import (
	"errors"
	"fmt"
	"project/logic"
	"project/logic/model/product"
	"project/logic/service/tool"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (serviceProduct *ServiceProductGroup) SearchProductPages(limits tool.RequestLimits) (products []*product.Product, formatLimits *tool.ResponseLimits, err error) {
	var total int64
	offset, err := limits.GetOffset()
	if err != nil {
		return nil, nil, err
	}
	err = logic.Gorm.
		Preload(clause.Associations).
		Count(&total).
		Offset(offset).
		Limit(limits.PageSize).
		Find(&products).Error
	if err != nil {
		return nil, nil, fmt.Errorf("查询产品失败: %v", err)
	}
	formatLimits = tool.NewLimits(total, limits.PageSize, limits.PageNum)
	return
}

func (serviceProduct *ServiceProductGroup) SearchProduct(pid uint) (*product.Product, error) {
	var product *product.Product
	if err := logic.Gorm.Preload(clause.Associations).First(&product, pid).Error; err != nil {
		return nil, fmt.Errorf("查询产品类别失败: %v", err)
	}
	return product, nil
}

func (serviceProduct *ServiceProductGroup) AddProduct(p product.Product) error {
	tx := logic.Gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	var product product.Product
	fields, err := serviceProduct.SearchFields(*tx, p.TypeId)
	if err != nil {
		tx.Rollback()
		return err
	}
	var jsons = make(map[string]string)
	for _, field := range fields {
		jsons[field.Name] = ""
	}

	if err := tx.Where("item_nubmer = ?", p.ItemNumber).First(&product).Error; err != nil {
		if err := tx.Create(p).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				tx.Rollback()
				return fmt.Errorf("产品名称已存在: %v", err)
			}
			tx.Rollback()
			return fmt.Errorf("新建产品失败: %v", err)
		}
	}
	return tx.Commit().Error
}

func (serviceProduct *ServiceProductGroup) UpdateProduct(p product.Product) error {
	var product product.Product
	if err := logic.Gorm.Where("item_nubmer = ?", p.ItemNumber).First(&product).Error; err != nil {
		if err := logic.Gorm.Model(&product).Updates(&p).Error; err != nil {
			return fmt.Errorf("更新产品失败,请检查: %v", err)
		}
	}
	return nil
}

func (serviceProduct *ServiceProductGroup) DeleteProduct(p product.Product) error {
	tx := logic.Gorm.Begin()
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

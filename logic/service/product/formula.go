package product

import (
	"errors"
	"fmt"
	"project/logic"
	"project/logic/model/product"

	"gorm.io/gorm"
)

func (serviceProduct *ServiceProductGroup) SearchFormulas(db gorm.DB, typeId uint) ([]product.Formula, error) {
	var Formula []product.Formula
	err := db.Where("type_id = ?", typeId).Order("name").Find(&Formula).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("未查询到相关Formula: %v", err)
		}
		return nil, fmt.Errorf("查询Formula失败: %v", err)
	}
	return Formula, nil
}

func (serviceProduct *ServiceProductGroup) SearchFormula(name string, typeId uint) (*product.Formula, error) {
	var Formula product.Formula
	err := logic.Gorm.Where("type_id = ? AND name = ?", typeId, name).First(&Formula).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("未查询到相关Formula: %v", err)
		}
		return nil, fmt.Errorf("查询Formula失败: %v", err)
	}
	return &Formula, nil
}

func (serviceProduct *ServiceProductGroup) AddFormula(Formula []product.Formula) error {
	for _, Formula := range Formula {
		findFormula, _ := ServiceProductGroupApp.SearchFormula(Formula.Name, Formula.TypeId)
		if findFormula != nil {
			return fmt.Errorf("此Type: %v的这条Formula已经创建: %v", findFormula.TypeId, findFormula.Name)
		}
	}
	err := logic.Gorm.Create(&Formula).Error
	if err != nil {
		return fmt.Errorf("插入Formula失败: %v", err)
	}
	return nil
}

func (serviceProduct *ServiceProductGroup) UpdateFormula(Formula []product.Formula) error {
	err := logic.Gorm.Updates(Formula).Error
	if err != nil {
		return fmt.Errorf("更新Formula失败: %v", err)
	}
	return nil
}

func (serviceProduct *ServiceProductGroup) DeleteFormulas(db gorm.DB, Formula []product.Formula) error {
	err := db.Delete(Formula).Error
	if err != nil {
		return fmt.Errorf("更新Formula失败: %v", err)
	}
	return nil
}

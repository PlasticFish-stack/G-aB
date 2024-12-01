package product

import (
	"errors"
	"fmt"
	"project/logic"
	"project/logic/model/product"

	"gorm.io/gorm"
)

func (serviceProduct *ServiceProductGroup) SearchFields(db gorm.DB, typeId uint) ([]product.Field, error) {
	var fields []product.Field
	err := db.Where("type_id = ?", typeId).Order("name").Find(&fields).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("未查询到相关Fields: %v", err)
		}
		return nil, fmt.Errorf("查询Fields失败: %v", err)
	}
	return fields, nil
}

func (serviceProduct *ServiceProductGroup) SearchField(name string, typeId uint) (*product.Field, error) {
	var field product.Field
	err := logic.Gorm.Where("type_id = ? AND name = ?", typeId, name).First(&field).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("未查询到相关Fields: %v", err)
		}
		return nil, fmt.Errorf("查询Fields失败: %v", err)
	}
	return &field, nil
}

func (serviceProduct *ServiceProductGroup) AddFields(fields []product.Field) error {
	for _, field := range fields {
		findField, _ := ServiceProductGroupApp.SearchField(field.Name, field.TypeId)
		if findField != nil {
			return fmt.Errorf("此Type: %v的这条Field已经创建: %v", findField.TypeId, findField.Name)
		}
	}
	err := logic.Gorm.Create(&fields).Error
	if err != nil {
		return fmt.Errorf("插入Fields失败: %v", err)
	}
	return nil
}

func (serviceProduct *ServiceProductGroup) UpdateFields(fields []product.Field) error {
	err := logic.Gorm.Updates(fields).Error
	if err != nil {
		return fmt.Errorf("更新Fields失败: %v", err)
	}
	return nil
}

func (serviceProduct *ServiceProductGroup) DeleteFields(db gorm.DB, fields []product.Field) error {
	err := db.Delete(fields).Error
	if err != nil {
		return fmt.Errorf("删除Fields失败: %v", err)
	}
	return nil
}

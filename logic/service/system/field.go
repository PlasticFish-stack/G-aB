package system

import (
	"fmt"
	"project/logic"
	"project/logic/model"
	"time"

	"gorm.io/gorm"
)

func (s *ServiceSystemGroup) SearchFieldGroup() ([]model.Field, error) {
	var fieldGroup []model.Field
	err := logic.Gorm.Model(&model.Field{}).Find(&fieldGroup).Error
	if err != nil {
		return nil, fmt.Errorf("搜索Field失败: %v", err)
	}
	return fieldGroup, nil
}

func (s *ServiceSystemGroup) SearchField(apiId uint, fieldName string) (*model.Field, error) {
	var field *model.Field
	err := logic.Gorm.Model(&model.Field{}).Where("field_name = ? AND parent_api_id = ?", fieldName, apiId).First(&field).Error
	if err != nil {
		return nil, fmt.Errorf("查询Fields失败: %v", err)
	}
	return field, nil
}

func (s *ServiceSystemGroup) AddFields(apiId uint, fields []model.Field) error {
	for _, field := range fields {
		findField, _ := s.SearchField(apiId, field.FieldName)
		if findField != nil {
			return fmt.Errorf("此api: %v的这条Field已经创建: %v", apiId, findField.FieldName)
		}
	}
	err := logic.Gorm.Create(&fields).Error
	if err != nil {
		return fmt.Errorf("新建fields失败: %v", err)
	}
	return nil
}

func (s *ServiceSystemGroup) UpdateField(requestField model.Field) error {
	var field *model.Field
	_, err := ServiceSystemGroupApp.SearchField(requestField.ParentApiId, requestField.FieldName)
	if err != nil {
		return err
	}
	err = logic.Gorm.Model(&field).Where("id = ?", field.Id).Omit("parent_api_id").Updates(&requestField).Error
	if err != nil {
		return fmt.Errorf("更新Field失败: %v", err)
	}
	return nil
}

func (s *ServiceSystemGroup) DeleteField(db *gorm.DB, fieldId uint) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Commit()
		}
	}()
	var field *model.Field
	err := tx.Where("id = ?", fieldId).First(&field).Error
	if err != nil {
		return fmt.Errorf("查询不到该field,请检查: %v", err)
	}
	err = tx.Model(&model.Field{}).Preload("Role").Find(&field).Error
	if err != nil {
		return fmt.Errorf("检查角色字段绑定失败: %v", err)
	}
	if len(field.Role) != 0 {
		return fmt.Errorf("该字段仍有角色使用,请检查并解绑")
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	field.FieldName = field.FieldName + "_is_deleted" + currentTime
	field.FieldDescription = field.FieldDescription + "_is_deleted" + currentTime
	if err := tx.Updates(&field).Error; err != nil {
		return fmt.Errorf("删除字段失败,请检查: %v", err)
	}
	if err := tx.Delete(&field).Error; err != nil {
		return fmt.Errorf("删除字段失败,请检查: %v", err)
	}
	return nil
}

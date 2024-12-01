package system

import (
	"fmt"
	"project/logic"
	"project/logic/model"
	"time"

	"gorm.io/gorm"
)

func (s *ServiceSystemGroup) SearchApiGroup() ([]model.Api, error) {
	var apiGroup []model.Api
	if err := logic.Gorm.Find(&apiGroup).Error; err != nil {
		return nil, fmt.Errorf("搜索Api列表失败: %v", err)
	}
	return apiGroup, nil
}

func (s *ServiceSystemGroup) SearchApi(menuId uint, apiName string) (*model.Field, error) {
	var api *model.Field
	err := logic.Gorm.Model(&model.Field{}).Where("field_name = ? AND parent_menu_id = ?", apiName, menuId).First(&api).Error
	if err != nil {
		return nil, fmt.Errorf("查询Api失败: %v", err)
	}
	return api, nil
}

func (s *ServiceSystemGroup) AddApi(api model.Api) error {
	findApi, _ := ServiceSystemGroupApp.SearchApi(api.ParentMenuId, api.ApiName)
	if findApi != nil {
		return fmt.Errorf("此Menu: %v的这条api已经创建: %v", api.ParentMenuId, api.ApiName)
	}
	err := logic.Gorm.Create(&api).Error
	if err != nil {
		return fmt.Errorf("创建api失败: %v", err)
	}
	return nil
}

func (s *ServiceSystemGroup) UpdateApi(requestApi model.Api) error {
	var api *model.Api
	_, err := ServiceSystemGroupApp.SearchApi(requestApi.ParentMenuId, requestApi.ApiName)
	if err != nil {
		return err
	}
	err = logic.Gorm.Model(&api).Where("id = ?", api.Id).Omit("parent_menu_id").Updates(&requestApi).Error
	if err != nil {
		return fmt.Errorf("更新api失败: %v", err)
	}
	return nil
}

func (s *ServiceSystemGroup) DeleteApi(db *gorm.DB, apiId uint) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Commit()
		}
	}()
	var api *model.Api
	err := tx.Where("id = ?", apiId).First(&api).Error
	if err != nil {
		return fmt.Errorf("查询不到该field,请检查: %v", err)
	}
	err = tx.Model(&model.Field{}).Preload("Role").Find(&apiId).Error
	if err != nil {
		return fmt.Errorf("检查角色字段绑定失败: %v", err)
	}
	if len(api.Role) != 0 {
		return fmt.Errorf("该字段仍有角色使用,请检查并解绑")
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	api.ApiName = api.ApiName + "_is_deleted" + currentTime
	api.ApiDescrption = api.ApiDescrption + "_is_deleted" + currentTime
	if err := tx.Updates(&api).Error; err != nil {
		return fmt.Errorf("删除字段失败,请检查: %v", err)
	}
	if err := tx.Delete(&api).Error; err != nil {
		return fmt.Errorf("删除字段失败,请检查: %v", err)
	}
	return nil
}

// func (s *ServiceSystemGroup) GetBindRoleApi(roleId uint) ([]model.Api, error) {
// 	var apiGroup []model.Api
// 	var roleCurrency model.Role
// 	if err := logic.Gorm.Model(&model.Role{}).Where("id = ?", roleId).First(&roleCurrency).Error; err != nil {

// 	}
// 	if err := logic.Gorm.Model(&roles).Association("Api").Find(&apiGroup); err != nil {
// 		return nil, fmt.Errorf("获取角色绑定菜单数组失败: %v", err)
// 	}
// 	return menus, nil
// }

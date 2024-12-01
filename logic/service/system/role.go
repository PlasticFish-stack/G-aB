package system

import (
	"errors"
	"fmt"
	"project/logic"
	"project/logic/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ServiceSystemGroup) SearchRole(db *gorm.DB, id uint) (*model.Role, error) {
	var role *model.Role
	if err := db.Model(&model.Role{}).Preload(clause.Associations).Where("id = ?", id).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("查询角色失败,未找到该角色: %v", err)
		}
		return nil, fmt.Errorf("查询角色失败: %v", err)
	}
	return role, nil
}

func (s *ServiceSystemGroup) SearchRoleGroup(db *gorm.DB) ([]model.Role, error) {
	var roleGroup []model.Role
	if err := db.Model(&model.Role{}).Find(&roleGroup).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("查询角色失败,未找到角色: %v", err)
		}
		return nil, fmt.Errorf("查询角色失败: %v", err)
	}
	return roleGroup, nil
}

func (s *ServiceSystemGroup) AddRole(db *gorm.DB, name string, ident string, description string) error {
	role := &model.Role{
		Name:        name,
		Identifier:  ident,
		Description: description,
	}
	if err := db.Model(&model.Role{}).Create(role).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("角色名称已存在: %v", err)
		}
		return fmt.Errorf("新建角色失败: %v", err)
	}
	return nil
}

func (s *ServiceSystemGroup) UpdateRole(db *gorm.DB, id uint, name string, identifier string, description string, status bool) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	var updateGroup = &model.Role{
		Identifier:  identifier,
		Name:        name,
		Description: description,
	}
	if err := tx.Model(&model.Role{}).Where("id = ?", id).Updates(&updateGroup).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新角色失败: %v", err)
	}
	if err := tx.Model(&model.Role{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新角色状态失败: %v", err)
	}
	return tx.Commit().Error
}

func (s *ServiceSystemGroup) DeleteRole(db *gorm.DB, id uint) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Commit()
		}
	}()
	var role *model.Role
	var apiField []*model.Api
	var menu []*model.Menu
	if role, err = s.SearchRole(tx, id); err != nil {
		tx.Rollback()
		return err
	}
	if menu, err = s.GetRoleBindMenu(tx, id); err != nil || len(menu) == 0 {
		tx.Rollback()
		return err
	}
	if apiField, err = s.GetRoleBindApiField(tx, id); err != nil || len(apiField) == 0 {
		tx.Rollback()
		return err
	}
	if err = tx.Delete(&role).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *ServiceSystemGroup) RoleGetAll() ([]model.Role, error) {
	var roles []model.Role
	var err error
	if roles, err = model.SearchRole(logic.Gorm); err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *ServiceSystemGroup) BindMenuToRole(db *gorm.DB, id uint, menuGroup []uint) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	fmt.Println(id, menuGroup)
	var role *model.Role
	var bindMenuGroup []model.Menu
	if role, err = s.SearchRole(tx, id); err != nil {
		tx.Rollback()
		return err
	}
	for _, id := range menuGroup {
		bindMenuGroup = append(bindMenuGroup, model.Menu{Global: model.Global{Id: id}})
	}
	if err := db.Model(&role).Association("Menu").Replace(bindMenuGroup); err != nil {
		tx.Rollback()
		return fmt.Errorf("该角色绑定菜单失败: %v", err)
	}
	return nil
}
func (s *ServiceSystemGroup) BindApiAndFieldToRole(db *gorm.DB, id uint, apiGroup map[uint][]uint) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var role *model.Role
	var bindApiGroup []model.Api
	var bindFieldGroup []model.Field
	var bindMenuGroup []*model.Menu
	var bindApiGroupSearch []uint
	var bindFieldGroupSearch []uint
	for apiId, fieldGroup := range apiGroup {
		bindApiGroupSearch = append(bindApiGroupSearch, apiId)
		bindFieldGroupSearch = append(bindFieldGroupSearch, fieldGroup...)
	}
	if role, err = s.SearchRole(tx, id); err != nil {
		tx.Rollback()
		return err
	}

	if bindMenuGroup, err = s.GetRoleBindMenu(tx, role.Id); err != nil {
		tx.Rollback()
		return err
	}
	var menuHashMap = make(map[uint]model.Menu)
	for _, menu := range bindMenuGroup {
		menuHashMap[menu.Id] = *menu
	}
	if err := tx.Model(&model.Api{}).Find(&bindApiGroup, bindApiGroupSearch).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("查询Api失败: %v", err)
	}
	var apiHashMap = make(map[uint]model.Api)
	for _, api := range bindApiGroup {
		apiHashMap[api.Id] = api
	}
	if err := tx.Model(&model.Field{}).Find(&bindFieldGroup, bindFieldGroupSearch).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("查询field失败: %v", err)
	}
	var fieldHashMap = make(map[uint]model.Field)
	for _, field := range bindFieldGroup {
		fieldHashMap[field.Id] = field
	}
	for apiId, fieldGroup := range apiGroup {
		if _, ok := menuHashMap[apiHashMap[apiId].ParentMenuId]; !ok {
			tx.Rollback()
			return fmt.Errorf("该角色未绑定对应菜单,不能绑定api")
		}
		for _, fieldId := range fieldGroup {
			if _, ok := apiHashMap[fieldHashMap[fieldId].ParentApiId]; !ok {
				tx.Rollback()
				return fmt.Errorf("该字段不属于对应api,无法绑定")
			}
		}
	}
	if err := tx.Model(&role).Association("Api").Replace(bindApiGroup); err != nil {
		tx.Rollback()
		return fmt.Errorf("该角色绑定api失败: %v", err)
	}
	if err := tx.Model(&role).Association("Field").Replace(bindFieldGroup); err != nil {
		tx.Rollback()
		return fmt.Errorf("该角色绑定field失败: %v", err)
	}
	return tx.Commit().Error
}

func (s *ServiceSystemGroup) GetRoleBindMenu(db *gorm.DB, roleId uint) ([]*model.Menu, error) {
	var menuGroup []*model.Menu
	role := &model.Role{
		Global: model.Global{Id: roleId},
	}
	if err := db.Model(&role).Association("Menu").Find(&menuGroup); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("该角色没有绑定任何菜单: %v", err)
		}
		return nil, fmt.Errorf("获取角色绑定菜单失败: %v", err)
	}
	return menuGroup, nil
}

func (s *ServiceSystemGroup) GetRoleBindApiField(db *gorm.DB, roleId uint) ([]*model.Api, error) {
	var apiGroup []*model.Api
	var fieldGroup []model.Field
	role := &model.Role{
		Global: model.Global{Id: roleId},
	}
	if err := db.Model(&role).Association("Api").Find(&apiGroup); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apiGroup, nil
		}
		return nil, fmt.Errorf("获取角色绑定Api失败: %v", err)
	}
	if err := db.Model(&role).Association("Field").Find(&fieldGroup); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apiGroup, nil
		}
		return nil, fmt.Errorf("获取角色绑定Api的字段失败: %v", err)
	}
	var apiHashMap = make(map[uint]*model.Api)
	for _, api := range apiGroup {
		apiHashMap[api.Id] = api
	}
	for _, field := range fieldGroup {
		if api, ok := apiHashMap[field.ParentApiId]; ok {
			api.Fields = append(api.Fields, field)
		}
	}
	return apiGroup, nil
}

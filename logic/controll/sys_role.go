package controll

import (
	"fmt"
	"project/logic"
	"project/logic/model"
)

func RoleAdd(name string, ident string, description string) error {
	role := &model.Role{
		Name:        name,
		Identifier:  ident,
		Description: description,
	}
	err := role.Add(logic.Gorm)
	if err != nil {
		return err
	}
	return nil
}

func RoleUpdate(id uint, name string, identifier string, description string, status bool) error {
	updateInfo := model.RoleUpdate{
		Global:      model.Global{Id: id},
		Name:        &name,
		Identifier:  &identifier,
		Description: &description,
		Status:      &status,
	}
	if err := updateInfo.Update(logic.Gorm); err != nil {
		return err
	}
	return nil
}

func RoleDelete(id uint, name string) error {
	var Role = model.Role{
		Global: model.Global{Id: id},
	}
	resultUser, err := Role.Search(logic.Gorm)
	if err != nil {
		return err
	}
	if name != resultUser.Name {
		return fmt.Errorf("角色名与id不匹配")
	}
	if err := resultUser.Delete(logic.Gorm); err != nil {
		return err
	}
	return nil
}

func RoleGetAll() ([]model.Role, error) {
	var roles []model.Role
	var err error
	if roles, err = model.SearchRole(logic.Gorm); err != nil {
		return nil, err
	}
	return roles, nil
}

func BindRoleMenu(roleId uint, menuId []uint) error {
	role := &model.Role{
		Global: model.Global{Id: roleId},
	}
	err := role.Bind(logic.Gorm, menuId)
	if err != nil {
		return err
	}
	return nil
}

//	func UnBindRoleMenu(roleId uint, menuId []uint) error {
//		role := &model.Role{
//			Id: roleId,
//		}
//		err := role.UnBind(logic.Gorm, menuId)
//		if err != nil {
//			return err
//		}
//		return nil
//	}
func GetBindRoleMenu(roleId uint) ([]model.Menu, error) {
	role := &model.Role{
		Global: model.Global{Id: roleId},
	}
	menus, err := role.GetBind(logic.Gorm)
	if err != nil {
		return nil, err
	}
	return menus, nil
}

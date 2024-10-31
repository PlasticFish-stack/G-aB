package controll

import (
	"fmt"
	"project/logic"
	"project/logic/model"
)

func MenuAdd(name string, ident string, description string, component string, path string, icon string, sort uint, parentid uint) error {
	menu := &model.Menu{
		Name:        name,
		Identifier:  ident,
		Description: description,
		Component:   component,
		Path:        path,
		Icon:        icon,
		Sort:        sort,
		ParentId:    parentid,
	}
	err := menu.Add(logic.Gorm)
	if err != nil {
		return err
	}
	return nil
}

func MenuUpdate(id uint, name string, ident string, description string, component string, path string, icon string, sort uint, parentid uint, status bool) error {
	updateInfo := model.MenuUpdate{
		Global:      model.Global{Id: id},
		Name:        &name,
		Identifier:  &ident,
		Description: &description,
		Component:   &component,
		Path:        &path,
		Icon:        &icon,
		Sort:        &sort,
		ParentId:    &parentid,
		Status:      &status,
	}
	if err := updateInfo.Update(logic.Gorm); err != nil {
		return err
	}
	return nil
}

func MenuDelete(id uint, name string) error {
	var Menu = model.Menu{
		Global: model.Global{Id: id},
	}
	resultMenu, err := Menu.Search(logic.Gorm)
	if err != nil {
		return err
	}
	if name != resultMenu.Name {
		return fmt.Errorf("菜单与id不匹配")
	}
	if err := resultMenu.Delete(logic.Gorm); err != nil {
		return err
	}
	return nil
}

func MenuGetAll() ([]*model.Menu, error) {
	var menus []*model.Menu
	var err error
	if menus, err = model.SearchTreeMenu(logic.Gorm); err != nil {
		return nil, err
	}
	return menus, nil
}

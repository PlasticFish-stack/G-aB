package system

import "project/logic/model"

type Role struct {
	model.Global
	Name        string  `gorm:"size:255;not null;unique" json:"name"`
	Identifier  string  `gorm:"size:255;not null;unique" json:"identifier"`
	Description string  `gorm:"size:255" json:"description"`
	Status      bool    `gorm:"default:true" json:"status"`
	Menu        []*Menu `gorm:"many2many:role_bind_menu" json:"-"`
}

// type RoleUpdate struct {
// 	model.Global
// 	Name        *string `json:"name"`
// 	Identifier  *string `json:"identifier"`
// 	Description *string `json:"description"`
// 	Status      *bool   `json:"status"`
// }

// func SearchRole(db *gorm.DB) ([]Role, error) {
// 	var roles []Role
// 	if err := db.Find(&roles).Error; err != nil {
// 		return nil, fmt.Errorf("查询角色失败: %v", err)
// 	}
// 	return roles, nil
// }

// func (role *Role) Search(db *gorm.DB) (*Role, error) {
// 	var resultRole Role
// 	if err := db.Find(&resultRole, role.Id).Error; err != nil {
// 		return nil, fmt.Errorf("查询角色失败: %v", err)
// 	}
// 	return &resultRole, nil
// }

// func (role *Role) Add(db *gorm.DB) error {
// 	if err := db.Create(role).Error; err != nil {
// 		if errors.Is(err, gorm.ErrDuplicatedKey) {
// 			return fmt.Errorf("角色名称已存在: %v", err)
// 		}
// 		return fmt.Errorf("新建角色失败: %v", err)
// 	}
// 	return nil
// }

// func (role *Role) Delete(db *gorm.DB) error {
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			panic(r)
// 		} else {
// 			tx.Commit()
// 		}
// 	}()
// 	currentTime := time.Now().Format("2006-01-02 15:04:05")
// 	role.Name = role.Name + "_is_deleted" + currentTime
// 	role.Identifier = role.Identifier + "_is_deleted" + currentTime
// 	if err := tx.Updates(&role).Error; err != nil {
// 		return fmt.Errorf("删除角色失败,请检查: %v", err)
// 	}
// 	if err := tx.Delete(&role).Error; err != nil {
// 		return fmt.Errorf("删除角色失败,请检查: %v", err)
// 	}
// 	return nil
// }

// func (role *RoleUpdate) Update(db *gorm.DB) error {
// 	var resultRole Role
// 	if err := db.First(&resultRole, role.Id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return fmt.Errorf("未查询到该角色: %v", err)
// 		}
// 		return fmt.Errorf("查询角色失败: %v", err)
// 	}
// 	if err := db.Model(&resultRole).Updates(role).Error; err != nil {
// 		return fmt.Errorf("更新角色失败,请检查: %v", err)
// 	}
// 	return nil
// }

// func (role *Role) Bind(db *gorm.DB, MenuIds []uint) error {
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			panic(r)
// 		}
// 	}()
// 	if err := tx.First(&role, role.Id).Error; err != nil {
// 		return fmt.Errorf("该角色不存在: %v", err)
// 	}
// 	var menus []Menu
// 	for _, id := range MenuIds {
// 		menus = append(menus, Menu{Global: Global{Id: id}})
// 	}
// 	// if err := db.Model(&role).Association("RoleBindMenu").Clear().Error; err != nil {
// 	// 	return fmt.Errorf("绑定菜单失败-清空菜单步骤: %v", err)
// 	// }
// 	if err := db.Model(&role).Association("Menu").Replace(menus); err != nil {
// 		return fmt.Errorf("绑定菜单失败: %v", err)
// 	}
// 	return nil
// }

// func (role *Role) GetBind(db *gorm.DB) ([]Menu, error) {
// 	var menus []Menu
// 	roles, err := role.Search(db)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := db.Model(&roles).Association("Menu").Find(&menus); err != nil {
// 		return nil, fmt.Errorf("获取角色绑定菜单数组失败: %v", err)
// 	}
// 	return menus, nil
// }

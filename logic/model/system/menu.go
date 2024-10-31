package system

import "project/logic/model"

type Menu struct {
	model.Global
	Name        string  `gorm:"size:255;not null;unique" json:"name"`
	Description string  `gorm:"size:255" json:"description"`
	Identifier  string  `gorm:"size:255;not null;unique" json:"identifier"`
	Component   string  `gorm:"size:255;default:'/null';not null" json:"component"`
	Path        string  `gorm:"size:255;default:'/null';not null" json:"path"`
	Icon        string  `gorm:"size:255" json:"icon"`
	Sort        uint    `gorm:"column:menus_sort;default:0" json:"sort"`
	ParentId    uint    `gorm:"default:0;not null" json:"parentId"`
	Status      bool    `gorm:"default:true" json:"status"`
	Children    []Menu  `gorm:"-" json:"children,omitempty"`
	Role        []*Role `gorm:"many2many:role_bind_menu" json:"-"`
}

// type MenuUpdate struct {
// 	model.Global
// 	Name        *string `json:"name"`
// 	Identifier  *string `json:"identifier"`
// 	Description *string `json:"description"`
// 	Status      *bool   `json:"status"`
// 	Component   *string `json:"component"`
// 	Path        *string `json:"path"`
// 	Icon        *string `json:"icon"`
// 	Sort        *uint   `json:"sort"`
// 	ParentId    *uint   `json:"parentId"`
// }

// func SearchFlatMenu(db *gorm.DB) ([]Menu, error) {
// 	var menus []Menu
// 	if err := db.Find(&menus).Error; err != nil {
// 		return nil, fmt.Errorf("查询菜单失败: %v", err)
// 	}
// 	return menus, nil
// }

// func SearchTreeMenu(db *gorm.DB) ([]*Menu, error) {
// 	var menus []Menu
// 	if err := db.Order("menus_sort").Find(&menus).Error; err != nil {
// 		return nil, fmt.Errorf("查询菜单失败: %v", err)
// 	}
// 	treeMapMenus := make(map[uint]*Menu)
// 	for i := range menus {
// 		menus[i].Children = []Menu{}
// 		treeMapMenus[menus[i].Id] = &menus[i]
// 	}
// 	var treeMenus []*Menu
// 	for i := range menus {
// 		menu := treeMapMenus[menus[i].Id]
// 		if menu.ParentId == 0 {
// 			treeMenus = append(treeMenus, menu)
// 		} else {
// 			treeMapMenus[menu.ParentId].Children = append(treeMapMenus[menu.ParentId].Children, *menu)
// 		}
// 	}
// 	for _, parentMenu := range treeMenus {
// 		sort.Slice(parentMenu.Children, func(i, j int) bool {
// 			return parentMenu.Children[i].Id < parentMenu.Children[j].Id
// 		})
// 	}
// 	return treeMenus, nil
// }

// func (menu *Menu) Search(db *gorm.DB) (*Menu, error) {
// 	var resultMenu Menu
// 	if err := db.Find(&resultMenu, menu.Id).Error; err != nil {
// 		return nil, fmt.Errorf("查询菜单失败: %v", err)
// 	}
// 	return &resultMenu, nil
// }

// func (menu *Menu) Add(db *gorm.DB) error {
// 	if err := db.Create(menu).Error; err != nil {
// 		if errors.Is(err, gorm.ErrDuplicatedKey) {
// 			return fmt.Errorf("菜单名称已存在: %v", err)
// 		}
// 		return fmt.Errorf("新建菜单失败: %v", err)
// 	}
// 	return nil
// }

// func (menu *Menu) Delete(db *gorm.DB) error {
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			panic(r)
// 		} else {
// 			tx.Commit()
// 		}
// 	}()
// 	var count int64
// 	if err := tx.Model(&Menu{}).Where("parent_id = ?", menu.Id).Count(&count).Error; err != nil {
// 		return fmt.Errorf("查询菜单是否有子菜单失败: %v", err)
// 	}
// 	if count > 0 {
// 		return fmt.Errorf("该菜单存在子菜单")
// 	}
// 	if err := tx.Model(menu).Preload("Role").Find(menu).Error; err != nil {
// 		return fmt.Errorf("检查角色菜单绑定失败: %v", err)
// 	}
// 	if len(menu.Role) != 0 {
// 		return fmt.Errorf("该菜单仍有角色使用,请检查并解绑")
// 	}
// 	currentTime := time.Now().Format("2006-01-02 15:04:05")
// 	menu.Name = menu.Name + "_is_deleted" + currentTime
// 	menu.Identifier = menu.Identifier + "_is_deleted" + currentTime
// 	if err := tx.Updates(&menu).Error; err != nil {
// 		return fmt.Errorf("删除菜单失败,请检查: %v", err)
// 	}
// 	if err := tx.Delete(&menu).Error; err != nil {
// 		return fmt.Errorf("删除菜单失败,请检查: %v", err)
// 	}
// 	return nil
// }

// func (menu *MenuUpdate) Update(db *gorm.DB) error {
// 	var resultMenu Menu
// 	if err := db.First(&resultMenu, menu.Id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return fmt.Errorf("未查询到该菜单: %v", err)
// 		}
// 		return fmt.Errorf("查询菜单失败: %v", err)
// 	}
// 	if err := db.Model(&resultMenu).Updates(menu).Error; err != nil {
// 		return fmt.Errorf("更新菜单失败,请检查: %v", err)
// 	}
// 	return nil
// }

//	func (menu *Menu) GetChild(db *gorm.DB) ([]uint, error) {
//		var childMenus []Menu
//		if err := db.Where("parent_id = ?", menu.Id).Find(&childMenus).Error; err != nil {
//			return nil, err
//		}
//		var childMenuIds []uint
//		for _, childMenu := range childMenus {
//			childMenuIds = append(childMenuIds, childMenu.Id)
//			subChildMenuIds, err := childMenu.GetChild(db)
//			if err != nil {
//				return nil, err
//			}
//			childMenuIds = append(childMenuIds, subChildMenuIds...)
//		}
//		return childMenuIds, nil
//	}
// func (menu *Menu) GetChild(db *gorm.DB) ([]uint, error) {
// 	var childMenuIds []uint
// 	stack := []uint{menu.Id} // 使用栈来替代递归

// 	for len(stack) > 0 {
// 		currentId := stack[len(stack)-1]
// 		stack = stack[:len(stack)-1]

// 		var childMenus []Menu
// 		if err := db.Where("parent_id = ?", currentId).Find(&childMenus).Error; err != nil {
// 			return nil, fmt.Errorf("查询子菜单失败: %v", err)
// 		}

// 		for _, childMenu := range childMenus {
// 			childMenuIds = append(childMenuIds, childMenu.Id)
// 			stack = append(stack, childMenu.Id) // 将子菜单入栈
// 		}
// 	}

// 	return childMenuIds, nil
// }

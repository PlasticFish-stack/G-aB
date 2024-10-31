package controll

import (
	"log"
	"project/logic"
	"project/logic/model"
	"sort"
)

type Meta struct {
	Title      string   `json:"title"`
	Icon       string   `json:"icon"`
	Rank       uint     `json:"rank,omitempty"`
	Auth       []string `json:"auths"`
	Field      []string `json:"fields"`
	Role       []string `json:"roles"`
	ShowParent bool     `json:"showParent,omitempty"`
	// ShowParent bool     `json:"showParent,omitempty"`
}
type PermissionRouter struct {
	Id        uint                `json:"id"`
	Name      string              `json:"name"`
	Path      string              `json:"path"`
	Meta      Meta                `json:"meta"`
	Component *string             `json:"component,omitempty"`
	Children  []*PermissionRouter `json:"children,omitempty"`
}

func GetRoutes(userId uint) ([]*PermissionRouter, error) {
	// 获取绑定角色
	user := &model.User{
		Global: model.Global{Id: userId},
	}
	roles, err := user.GetBind(logic.Gorm)
	if err != nil {
		return nil, err
	}
	allMenus := make(map[uint]*model.Menu)

	// allMenus := make([]*model.Menu, 0)
	// menuMap := make(map[uint]*model.Menu)
	roleNamesMap := make(map[uint][]string) // 用于存储每个菜单的角色名称

	// 遍历角色获取菜单
	for _, role := range roles {
		getMenus, err := role.GetBind(logic.Gorm)
		if err != nil {
			return nil, err
		}
		for i := range getMenus {
			allMenus[getMenus[i].Id] = &getMenus[i]
			roleNamesMap[getMenus[i].Id] = append(roleNamesMap[getMenus[i].Id], role.Identifier)
		}
		// for _, menu := range getMenus {
		// 	menuCopy := menu
		// 	roleIdentifier := role.Identifier

		// 	if existingMenu, exists := menuMap[menuCopy.Id]; exists {
		// 		// 检查角色名称是否已存在
		// 		if !roleExists(roleNamesMap[existingMenu.Id], roleIdentifier) {
		// 			roleNamesMap[existingMenu.Id] = append(roleNamesMap[existingMenu.Id], roleIdentifier)
		// 		}
		// 	} else {
		// 		// 初始化菜单的角色列表并添加角色名称
		// 		roleNamesMap[menuCopy.Id] = []string{roleIdentifier}
		// 		allMenus = append(allMenus, menuCopy)
		// 		menuMap[menuCopy.Id] = allMenus[len(allMenus)-1]
		// 	}
		// }
	}
	// 局部变量来存储最终的路由
	var router []*PermissionRouter

	// 创建一个 map 来存储所有 PermissionRouter，以便快速访问
	routerMap := make(map[uint]*PermissionRouter)

	// 初始化所有菜单
	for _, p := range allMenus {
		permissions := &PermissionRouter{
			Id:   p.Id,
			Path: p.Path,
			Name: p.Identifier,
			Meta: Meta{
				Title: p.Name,
				Icon:  p.Icon,
				Role:  roleNamesMap[p.Id], // 使用辅助 map 存储的角色名称
				Rank:  p.Sort,
			},
		}
		routerMap[permissions.Id] = permissions

		if p.ParentId == 0 {
			router = append(router, permissions)
		}
	}

	// 构建父子关系
	for _, p := range allMenus {
		if p.ParentId != 0 {
			if parent, exists := routerMap[p.ParentId]; exists {
				permissions := routerMap[p.Id]
				permissions.Component = &p.Component
				parent.Children = append(parent.Children, permissions)
			} else {
				log.Printf("Parent menu not found for child menu ID: %d\n", p.Id)
			}
		}
	}

	// 设置显示父节点标志
	for _, p := range router {
		if p.Children != nil && len(p.Children) == 1 {
			p.Children[0].Meta.ShowParent = true
		}
	}
	for _, parentMenu := range router {
		sort.Slice(parentMenu.Children, func(i, j int) bool {
			return parentMenu.Children[i].Id < parentMenu.Children[j].Id
		})
	}
	return router, nil
}

// func roleExists(roles []string, roleName string) bool {
// 	for _, v := range roles {
// 		if v == roleName {
// 			return true
// 		}
// 	}
// 	return false
// }

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

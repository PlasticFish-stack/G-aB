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

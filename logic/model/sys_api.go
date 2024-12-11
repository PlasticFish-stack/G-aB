package model

type Api struct {
	Global
	ApiName       string  `gorm:"size:255" json:"name"`
	ApiPath       string  `gorm:"size:255" json:"path"`
	ApiDescrption string  `gorm:"size:255" json:"description"`
	ApiMethod     string  `gorm:"size:255" json:"method"`
	ApiType       string  `gorm:"size:255" json:"type"`
	ParentMenuId  uint    `gorm:"default:0;not null" json:"parentId"`
	Role          []*Role `gorm:"many2many:role_bind_api" json:"-"`
	Fields        []Field `gorm:"-" json:"fields,omitempty"`
}

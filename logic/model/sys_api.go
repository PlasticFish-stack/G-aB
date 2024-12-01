package model

type Api struct {
	Global
	ApiName       string
	ApiPath       string
	ApiDescrption string
	ApiMethod     string
	ApiType       string
	ParentMenuId  uint
	Role          []*Role `gorm:"many2many:role_bind_api" json:"-"`
	Fields        []Field `gorm:"-" json:"fields,omitempty"`
}

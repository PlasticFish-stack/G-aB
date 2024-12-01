package model

type Field struct {
	Global
	FieldName        string
	FieldDescription string
	ParentApiId      uint
	AllowFields      string
	Role             []*Role `gorm:"many2many:role_bind_field" json:"-"`
}

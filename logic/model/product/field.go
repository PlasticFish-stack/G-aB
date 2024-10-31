package product

import "project/logic/model"

type Field struct {
	model.Global
	Name     string `gorm:"size:255" json:"name"`
	NickName string `gorm:"size:255" json:"nickname"`
	TypeId   uint   `gorm:"index"`
}

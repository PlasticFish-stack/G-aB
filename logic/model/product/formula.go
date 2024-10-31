package product

import "project/logic/model"

type Formula struct {
	model.Global
	Name     string `gorm:"size:255" json:"name"`
	NickName string `gorm:"size:255" json:"nickname"`
	Formula  string `gorm:"size:255" json:"formula"`
	TypeId   uint
}

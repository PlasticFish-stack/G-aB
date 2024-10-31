package product

import (
	"project/logic/model"
)

type Type struct {
	model.Global
	Name        string    `gorm:"size:255" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	Sort        uint      `gorm:"column:type_sort;default:0" json:"sort"`
	ParentId    uint      `gorm:"default:0" json:"parentId"`
	Tax         float64   `gorm:"default:0" json:"tax"`
	Children    []Type    `gorm:"-" json:"children,omitempty"`
	Fields      []Field   `gorm:"foreignKey:TypeId" json:"fields"`
	Formulas    []Formula `gorm:"foreignKey:TypeId" json:"formulas"`
	Product     []Product `gorm:"foreignKey:TypeId" json:"products"`
}

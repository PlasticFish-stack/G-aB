package product

import (
	"project/logic/model"
)

type Brand struct {
	model.Global
	Name        string    `gorm:"unique" json:"name"`
	Description string    `json:"description"`
	Products    []Product `gorm:"foreignKey:BrandId"`
}

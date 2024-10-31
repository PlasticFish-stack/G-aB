package excel

import (
	"project/logic/model"
	"project/logic/model/product"
)

type ExcelLog struct {
	model.Global
	FileName     string         `gorm:"size:255;comment:文件名" json:"fileName"`
	Description  string         `gorm:"size:255;comment:描述" json:"description"`
	ProductCosts []product.Cost `gorm:"foreignKey:ExcelLogId" json:"-"`
}

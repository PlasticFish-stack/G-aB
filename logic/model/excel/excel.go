package excel

import (
	"project/logic/model"
	"project/logic/model/product"
)

// Excel在数据库表内的字段
type ExcelLog struct {
	model.Global
	FileName     string         `gorm:"size:255;comment:文件名" json:"fileName"`
	Description  string         `gorm:"size:255;comment:描述" json:"description"`
	ProductCosts []product.Cost `gorm:"foreignKey:ExcelLogId" json:"costs"`
}

// 导出时从前端接受到json字段
type ExportType struct {
	CurrencyName   string `gorm:"-" json:"currencyName"`
	ProductTypeId  uint   `gorm:"-" json:"productTypeId"`
	ProductBrandId uint   `gorm:"-" json:"productBrandId"`
}

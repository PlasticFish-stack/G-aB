package product

import "project/logic/model"

type Cost struct {
	model.Global
	ProductID    uint    `gorm:"comment:产品id" json:"-"` //product外键
	ExcelLogId   uint    `gorm:"comment:Excel的id" json:"excelId"`
	ItemNumber   string  `gorm:"comment:产品货号" json:"itemNumber"`
	ExcelName    string  `gorm:"comment:Excel的Name" json:"excelName"`       //operationlog表格外键
	Cost         float64 `gorm:"comment:成本价" json:"cost"`                   //成本价
	CurrencyName string  `gorm:"size:255;comment:货币名称" json:"currencyName"` //货币名称外键
	CurrencyCost float64 `gorm:"comment:当时汇率" json:"currencyCost"`
	DwPrice      float64 `gorm:"comment:得物价格" json:"dwPrice"`
	DwSales      float64 `gorm:"comment:得物销售额" json:"dwSale"`
	IsAuto       bool    `gorm:"comment:是否被人为修改过" json:"isAuto"`
}

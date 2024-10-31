package product

import "project/logic/model"

type Cost struct {
	model.Global
	ProductID    uint    `gorm:"comment:产品id" json:"-"`                     //product外键
	ExcelLogId   uint    `gorm:"comment:Excel的id" json:"-"`                 //operationlog表格外键
	Cost         float64 `gorm:"comment:成本价" json:"cost"`                   //成本价
	CurrencyName string  `gorm:"size:255;comment:货币名称" json:"currencyName"` //货币名称外键
}

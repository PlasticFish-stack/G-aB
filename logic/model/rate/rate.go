package rate

import (
	"project/logic/model"
	"project/logic/model/product"
	"time"
)

type Rate struct {
	model.Global
	CurrencyName  string         `gorm:"size:255;not null;unique" json:"currencyName"`
	DescriptionEn string         `gorm:"size:255" json:"descriptionEn"`
	DescriptionCn string         `gorm:"size:255" json:"descriptionCn"`
	Country       string         `gorm:"size:255" json:"country"`
	Organization  string         `gorm:"size:255" json:"organization"`
	Cost          float64        `gorm:"not null" json:"cost"`
	CountryIcon   string         `gorm:"size:255" json:"countryIcon"`
	UpdateTime    time.Time      `gorm:"api_update_time" json:"apiUpdateTime"`
	Sort          uint           `json:"sort"`
	ProductCost   []product.Cost `gorm:"foreignKey:CurrencyName;references:CurrencyName" json:"-"`
}

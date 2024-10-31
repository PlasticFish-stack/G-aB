package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Rate struct {
	Global
	CurrencyName  string        `gorm:"size:255;not null;unique" json:"currencyName"`
	DescriptionEn string        `gorm:"size:255" json:"descriptionEn"`
	DescriptionCn string        `gorm:"size:255" json:"descriptionCn"`
	Cost          float64       `gorm:"not null" json:"cost"`
	CountryIcon   string        `gorm:"size:255" json:"countryIcon"`
	UpdateTime    time.Time     `gorm:"api_update_time" json:"apiUpdateTime"`
	Sort          uint          `json:"sort"`
	ProductCost   []ProductCost `gorm:"foreignKey:CurrencyId" json:"-"`
}

func RateSearchAndPush(db *gorm.DB, rates []Rate) ([]Rate, error) {
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "currency_name"}},
		DoUpdates: clause.AssignmentColumns([]string{"currency_name", "cost", "update_time", "updated_at"}),
	}).Create(&rates).Error
	if err != nil {
		return nil, err
	}
	return rates, nil
}

func RateSearch(db *gorm.DB) ([]Rate, error) {
	var rates []Rate
	if err := db.Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}

func (r *Rate) RateUpdate(db *gorm.DB) error {
	var rate Rate
	if err := db.Where("currency_name = ?", &r.CurrencyName).First(&rate).Error; err != nil {
		return fmt.Errorf("查询不到该币种: %v", err)
	}
	if err := db.Model(&rate).Updates(map[string]interface{}{
		"description_en": r.DescriptionEn,
		"description_cn": r.DescriptionCn,
		"country_icon":   r.CountryIcon,
		"sort":           r.Sort,
	}).Error; err != nil {
		return fmt.Errorf("更改到该币种信息失败: %v", err)
	}
	return nil
}

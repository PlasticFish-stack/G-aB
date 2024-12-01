package rate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"project/logic"
	"project/logic/model/rate"
	"sort"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RateResponse struct {
	Result             string             `json:"result"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeLastUpdateUtc  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	TimeNextUpdateUtc  string             `json:"time_next_update_utc"`
	BaseCode           string             `json:"base_code"`
	ConversionRates    map[string]float64 `json:"conversion_rates"`
}

var token = "6e4911c07c8f13f11da1962b"

// var wg sync.WaitGroup

func (r *ServiceRateGroup) RataApiUpdate() ([]rate.Rate, error) {
	httpAddress := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%v/latest/CNY", token)
	resp, err := http.Get(httpAddress)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败,请通知管理员测试token是否过期: %v", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("解析接口数据失败: %v", err)
	}
	var RateResponse RateResponse
	if err := json.Unmarshal(body, &RateResponse); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}
	var rates []rate.Rate
	for k, v := range RateResponse.ConversionRates {
		rate := rate.Rate{
			CurrencyName: k,
			Cost:         v,
			UpdateTime:   time.Unix(RateResponse.TimeLastUpdateUnix, 0),
		}
		rates = append(rates, rate)
	}
	_, err = RateSearchAndPush(logic.Gorm, rates)
	if err != nil {
		return nil, err
	}
	sort.Slice(rates, func(i, j int) bool {
		// 根据 "name" 字段的字典顺序排序
		return rates[i].Id < rates[j].Id
	})
	return rates, nil
}

func (r *ServiceRateGroup) RateUpdate(updateRate rate.Rate) error {
	var rates *rate.Rate
	if err := logic.Gorm.Where("currency_name = ?", &updateRate.CurrencyName).First(&rates).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("查询不到该币种: %v", err)
		}
	}
	if err := logic.Gorm.Model(&rate.Rate{}).Where("currency_name = ?", rates.CurrencyName).Updates(map[string]interface{}{
		"description_en": updateRate.DescriptionEn,
		"description_cn": updateRate.DescriptionCn,
		"country_icon":   updateRate.CountryIcon,
		"country":        updateRate.Country,
		"organization":   updateRate.Organization,
		"sort":           updateRate.Sort,
	}).Error; err != nil {
		return fmt.Errorf("更改该币种信息失败: %v", err)
	}
	return nil
}

func (r *ServiceRateGroup) RateGet() ([]rate.Rate, error) {
	var rates []rate.Rate
	err := logic.Gorm.Order("id ASC").Find(&rates).Error
	if err != nil {
		return nil, err
	}
	return rates, nil
}

func (r *ServiceRateGroup) RateGetName(name string) (*rate.Rate, error) {
	var rate *rate.Rate
	err := logic.Gorm.Where("currency_name = ? ", name).First(&rate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("查询不到该币种: %v", err)
		}
	}
	return rate, nil
}

func RateSearchAndPush(db *gorm.DB, rates []rate.Rate) ([]rate.Rate, error) {
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "currency_name"}},
		DoUpdates: clause.AssignmentColumns([]string{"currency_name", "cost", "update_time", "updated_at"}),
	}).Create(&rates).Error
	if err != nil {
		return nil, err
	}
	return rates, nil
}

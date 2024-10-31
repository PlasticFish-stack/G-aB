package controll

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"project/logic"
	"project/logic/model"
	"sort"
	"time"
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

func RataApiUpdate() ([]model.Rate, error) {
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
	// fmt.Println(RateResponse)
	var rates []model.Rate
	for k, v := range RateResponse.ConversionRates {
		rate := model.Rate{
			CurrencyName: k,
			Cost:         v,
			UpdateTime:   time.Unix(RateResponse.TimeLastUpdateUnix, 0),
		}
		rates = append(rates, rate)
	}
	// for i := range rates {
	_, err = model.RateSearchAndPush(logic.Gorm, rates)
	if err != nil {
		return nil, err
	}
	// }
	sort.Slice(rates, func(i, j int) bool {
		// 根据 "name" 字段的字典顺序排序
		return rates[i].Id < rates[j].Id
	})
	return rates, nil
}

func RateUpdate(name string, descEn string, descCn string, icon string, sort uint) error {
	rate := &model.Rate{
		CurrencyName:  name,
		DescriptionEn: descEn,
		DescriptionCn: descCn,
		CountryIcon:   icon,
		Sort:          sort,
	}
	err := rate.RateUpdate(logic.Gorm)
	if err != nil {
		return err
	}
	return nil
}

func RateGet() ([]model.Rate, error) {
	rates, err := model.RateSearch(logic.Gorm)
	if err != nil {
		return nil, err
	}
	return rates, nil
}

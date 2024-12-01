package excel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"project/logic"
	"project/logic/model/excel"
	prod "project/logic/model/product"
	"project/logic/service/product"
	"project/logic/service/rate"
	"project/logic/service/tool"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgtype"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Item struct {
	BrandName      string            `json:"brandName"`
	BrandId        string            `json:"brandId"`
	TypeId         string            `json:"typeId"`
	TypeName       string            `json:"typeName"`
	ItemNumber     string            `json:"itemNumber"`
	SKU            string            `json:"sku"`
	SPU            string            `json:"spu"`
	Quantity       string            `json:"quantity"`
	Specifications string            `json:"specifications"`
	Barcode        string            `json:"barcode"`
	CustomsCode    string            `json:"customsCode"`
	Description    string            `json:"description"`
	Cost           string            `json:"cost"`
	Color          string            `json:"color"`
	Image          string            `json:"image"`
	Options        map[string]string `json:"options"`
}

type ExcelQuery struct {
	StartCreateTime *time.Time `json:"startCreateTime"`
	EndCreateTime   *time.Time `json:"endCreateTime"`
	StartUpdateTime *time.Time `json:"startUpdateTime"`
	EndUpdateTime   *time.Time `json:"endUpdateTime"`
}

func (q *ExcelQuery) BuildQuery() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if q.StartCreateTime != nil {
			db = db.Where("created_at >= ?", q.StartCreateTime)
		}

		if q.EndCreateTime != nil {
			db = db.Where("created_at <= ?", q.EndCreateTime)
		}

		if q.StartUpdateTime != nil {
			db = db.Where("updated_at >= ?", q.StartUpdateTime)
		}

		if q.EndUpdateTime != nil {
			db = db.Where("updated_at <= ?", q.EndUpdateTime)
		}
		return db
	}
}

var dir, _ = os.Getwd()
var excelFilePath = filepath.Join(dir, "public", "moban.xlsx")

func OpenExcelFile(filePath string) (*excelize.File, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	return f, nil
}

func SearchExcelPages(limits tool.RequestLimits, query *ExcelQuery) ([]*excel.ExcelLog, *tool.ResponseLimits, error) {
	var total int64
	offset, err := limits.GetOffset()
	if err != nil {
		return nil, nil, err
	}
	var excels []*excel.ExcelLog
	if err := logic.Gorm.
		Model(&excel.ExcelLog{}).
		Scopes(query.BuildQuery()).
		Order("created_at DESC").
		Count(&total).
		Offset(offset).
		Limit(limits.PageSize).
		Find(&excels).Error; err != nil {
		return nil, nil, fmt.Errorf("搜索excel页错误 :%v", err)
	}
	formatLimits := tool.NewLimits(total, limits.PageSize, limits.PageNum)
	return excels, formatLimits, nil
}

func SearchExcelCost(id uint) ([]prod.Cost, error) {
	var costs []prod.Cost
	if err := logic.Gorm.
		Model(&prod.Cost{}).
		Where("excel_log_id = ?", id).
		Order("created_at ASC").
		Find(&costs).Error; err != nil {
		return nil, fmt.Errorf("搜索COST失败 :%v", err)
	}
	return costs, nil
}

func (e *ServiceExcelGroup) ExportExcel(export excel.ExportType) ([]byte, error) {
	productBrand, err := product.ServiceProductGroupApp.SearchBrandId(export.ProductBrandId)
	if err != nil {
		return nil, err
	}
	productType, err := product.ServiceProductGroupApp.SearchType(export.ProductTypeId)
	if err != nil {
		return nil, err
	}
	rate, err := rate.ServiceRateGroupApp.RateGetName(export.CurrencyName)
	if err != nil {
		return nil, err
	}
	f, err := OpenExcelFile(excelFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	configTypeId := strconv.FormatUint(uint64(productType.Id), 10)
	configBrandId := strconv.FormatUint(uint64(productBrand.Id), 10)
	configRateCost := strconv.FormatFloat(rate.Cost, 'f', -1, 64)
	configParam := []interface{}{"Sheet1", configTypeId, productType.Name, configBrandId, productBrand.Name, rate.CurrencyName, configRateCost}
	productTypeParam := []string{}
	productTypeNickName := []string{}
	if productType.Fields != nil {
		for _, v := range productType.Fields {
			productTypeNickName = append(productTypeNickName, v.NickName)
			productTypeParam = append(productTypeParam, v.Name)
		}
	}
	// configSheetIndex, _ := f.GetSheetIndex("Config")
	err = f.SetSheetRow("Config", "A2", &configParam)
	if err != nil {
		return nil, fmt.Errorf("初始化模板表格失败: %v", err)
	}
	err = f.SetSheetRow("Sheet1", "L1", &productTypeParam)
	if err != nil {
		return nil, fmt.Errorf("初始化模板表格失败: %v", err)
	}
	err = f.SetSheetRow("Sheet1", "L2", &productTypeNickName)
	if err != nil {
		return nil, fmt.Errorf("初始化模板表格失败: %v", err)
	}
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (e *ServiceExcelGroup) CheckExcel(excel *excelize.File) ([]*Item, map[string]string, error) {
	defer func() {
		if err := excel.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	books := excel.GetSheetMap()

	configParam := map[string]string{
		"bookName":         "",
		"productTypeId":    "",
		"productTypeName":  "",
		"productBrandId":   "",
		"productBrandName": "",
		"currencyName":     "",
		"rate":             "",
	}
	rows, err := excel.GetRows("Config")
	if err != nil {
		return nil, nil, err
	}
	for i, key := range rows[0] {
		if i < len(rows[1]) { // 确保第二行有相同数量的列
			configParam[key] = rows[1][i]
		}
	}
	haveBook := false
	for _, name := range books {
		if name == configParam["bookName"] {
			haveBook = true
			break
		}
	}
	if !haveBook {
		return nil, nil, fmt.Errorf("config表内定义的数据表名未找到")
	}

	productTypeId, err := strconv.ParseUint(configParam["productTypeId"], 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("查询config表格内TypeId列是否不为数字: %v", err)
	}
	productType, err := product.ServiceProductGroupApp.SearchType(uint(productTypeId))
	if err != nil {
		return nil, nil, err
	}
	if productType.Name != configParam["productTypeName"] {
		return nil, nil, fmt.Errorf("查询config表格内TypeName列是否有误,id为%v的name应该为%v", productType.Id, productType.Name)
	}

	productBrandId, err := strconv.ParseUint(configParam["productBrandId"], 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("查询config表格内BrandId列是否不为数字: %v", err)
	}
	productBrand, err := product.ServiceProductGroupApp.SearchBrandId(uint(productBrandId))
	if err != nil {
		return nil, nil, err
	}
	if productBrand.Name != configParam["productBrandName"] {
		return nil, nil, fmt.Errorf("查询config表格内BrandName列是否有误,id为%v的name应该为%v", productBrand.Id, productBrand.Name)
	}
	_, err = rate.ServiceRateGroupApp.RateGetName(configParam["currencyName"])
	if err != nil {
		return nil, nil, err
	}
	var headerValue = map[string]interface{}{
		"Itemnumber":     "",
		"SKU":            "",
		"SPU":            "",
		"Quantity":       "",
		"Specifications": "",
		"Barcode":        "",
		"Customscode":    "",
		"Description":    "",
		"Color":          "",
		"Cost":           "",
		"Image":          "",
		"Options":        map[string]string{},
	}
	sheetRows, err := excel.GetRows(configParam["bookName"])
	if err != nil {
		return nil, nil, fmt.Errorf("config表内定义的数据表名未找到")
	}
	var vailOptions = make(map[string]int)
	for _, field := range productType.Fields {
		vailOptions[field.Name] = 0
	}
	for _, rows := range sheetRows[:1] {
		for index, row := range rows {
			if index > 10 {
				if _, exists := vailOptions[row]; exists {
					headerValue["Options"].(map[string]string)[row] = ""
				} else {
					return nil, nil, fmt.Errorf("type不存在%v值,请检查", row)
				}
			} else {
				if _, exists := headerValue[row]; exists {
					headerValue[row] = "" // 只有在不存在时才添加
				} else {
					return nil, nil, fmt.Errorf("表格内不存在%v值,请检查", row)
				}
			}
		}
	}
	var bodyData = []*Item{}
	for _, row := range sheetRows[2:] {
		rowData := &Item{
			Options: make(map[string]string),
		}
		for i, cell := range row {
			if i < len(sheetRows[0]) {
				columnName := sheetRows[0][i]
				if i > 10 {
					rowData.Options[columnName] = cell
				} else {
					switch columnName { // 使用 switch 来填充结构体字段
					case "Itemnumber":
						if cell == "" {
							return nil, nil, fmt.Errorf("货号不能为空")
						}
						rowData.ItemNumber = strings.TrimSpace(cell)
					case "SKU":
						rowData.SKU = strings.TrimSpace(cell)
					case "SPU":
						rowData.SPU = strings.TrimSpace(cell)
					case "Quantity":
						rowData.Quantity = strings.TrimSpace(cell)
					case "Specifications":
						rowData.Specifications = strings.TrimSpace(cell)
					case "Barcode":
						rowData.Barcode = strings.TrimSpace(cell)
					case "Customscode":
						rowData.CustomsCode = strings.TrimSpace(cell)
					case "Description":
						rowData.Description = strings.TrimSpace(cell)
					case "Color":
						rowData.Color = strings.TrimSpace(cell)
					case "Cost":
						rowData.Cost = strings.TrimSpace(cell)
					case "Image":
						rowData.Image = strings.TrimSpace(cell)
					}

				}
			}
		}
		rowData.BrandId = configParam["productBrandId"]
		rowData.BrandName = configParam["productBrandName"]
		rowData.TypeId = configParam["productTypeId"]
		rowData.TypeName = configParam["productTypeName"]
		bodyData = append(bodyData, rowData)
	}
	// for _, v := range bodyData {
	// 	b, _ := json.Marshal(v)
	// 	var out bytes.Buffer
	// 	err = json.Indent(&out, b, "", "    ")
	// 	fmt.Println(out.String())
	// }
	// fmt.Println(bodyData)
	return bodyData, configParam, nil
}
func validateImportParams(fileName, description string, importExcel []*Item, configParam map[string]string) error {
	if fileName == "" {
		return errors.New("文件名不能为空")
	}
	if description == "" {
		return errors.New("描述不能为空")
	}
	if len(importExcel) == 0 {
		return errors.New("导入数据不能为空")
	}
	if configParam == nil || configParam["currencyName"] == "" {
		return errors.New("货币信息不能为空")
	}
	return nil
}
func (e *ServiceExcelGroup) ImportExcel(fileName string, description string, importExcel []*Item, configParam map[string]string) error {
	tx := logic.Gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := validateImportParams(fileName, description, importExcel, configParam); err != nil {
		tx.Rollback()
		return err
	}
	currencyPrice, err := rate.ServiceRateGroupApp.RateGet()
	if err != nil {
		tx.Rollback()
		return err
	}
	var currencyPrices = make(map[string]float64)
	for _, v := range currencyPrice {
		currencyPrices[v.CurrencyName] = v.Cost
	}

	excelLog := &excel.ExcelLog{
		FileName:    fileName,
		Description: description,
	}
	if err := tx.Where("file_name = ?", fileName).First(&excelLog).Error; err == nil {
		tx.Rollback()
		return errors.New("已存在同名excel文件,不允许重复导入")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return fmt.Errorf("检查文件名时发生错误: %v", err)
	}
	if err := tx.Model(&excel.ExcelLog{}).Create(&excelLog).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建excelLog失败: %v", err)
	}
	var logCosts []prod.Cost
	for _, row := range importExcel {
		var currencyProduct = &prod.Product{
			Barcode:        row.Barcode,
			BrandId:        parseUint(row.BrandId),
			TypeId:         parseUint(row.TypeId),
			TypeName:       row.TypeName,
			Color:          row.Color,
			ItemNumber:     row.ItemNumber,
			Quantity:       uint64(parseUint(row.Quantity)),
			Sku:            row.SKU,
			Spu:            row.SPU,
			Specifications: row.Specifications,
			Customscode:    row.CustomsCode,
			Description:    row.Description,
		}
		optionsJSON, err := json.Marshal(row.Options)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("JSON序列化错误: %v", err)
		}
		var jsonbValue pgtype.JSONB
		if err = jsonbValue.Set(string(optionsJSON)); err != nil {
			tx.Rollback()
			return fmt.Errorf("设置JSONB错误: %v", err)
		}
		currencyProduct.Options = jsonbValue
		var logCost prod.Cost
		result := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "item_number"}}, // 冲突检查的唯一字段
			DoUpdates: clause.AssignmentColumns([]string{
				"barcode", "brand_id", "type_id",
				"color", "quantity", "sku", "spu", "specifications",
				"customscode", "description", "options",
				// 列出所有需要更新的字段
			}),
		}).Create(&currencyProduct)
		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("产品批量修改错误: %v", result.Error)
		}
		// if err != nil {
		// 	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 创建新产品

		// 处理Options

		// 	} else {
		// 		return err
		// 	}
		// }
		err = tx.Where("item_number = ?", row.ItemNumber).First(&currencyProduct).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("搜索产品失败,请检查: %v", err)
		}
		costPrice, err := strconv.ParseFloat(row.Cost, 64)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("转义Cost失败: %v", err)
		}

		logCost = prod.Cost{
			ExcelLogId:   excelLog.Id,
			ExcelName:    fileName,
			IsAuto:       true,
			ProductID:    currencyProduct.Id,
			Cost:         costPrice,
			CurrencyName: configParam["currencyName"],
			CurrencyCost: currencyPrices[configParam["currencyName"]],
		}
		logCosts = append(logCosts, logCost)
	}
	err = tx.Create(&logCosts).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("添加成本价失败: %v", err)
	}
	// for _, v := range logCosts {
	// 	b, _ := json.Marshal(v)
	// 	var out bytes.Buffer
	// 	err = json.Indent(&out, b, "", "    ")
	// 	if err != nil {
	// 		tx.Rollback()
	// 	}
	// 	fmt.Println(out.String())
	// }

	return tx.Commit().Error
}
func parseUint(s string) uint {
	val, _ := strconv.ParseUint(s, 10, 64)
	return uint(val)
}

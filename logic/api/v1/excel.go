package V1

import (
	"fmt"
	"net/http"
	"project/logic/model/excel"
	excels "project/logic/service/excel"
	"project/logic/service/tool"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ExcelApi struct{}
type Data struct {
	FileName    string            `json:"fileName"`
	Descrption  string            `json:"description"`
	ConfigParam map[string]string `json:"configParam"`
	ExcelData   []excels.Item     `json:"excelData"`
}

func (e *ExcelApi) ExcelExport(c *gin.Context) {
	responseBody := &Response{Success: true}
	var export excel.ExportType
	err := c.ShouldBindJSON(&export)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	excelData, err := ExcelService.ExportExcel(export)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\"moban.xlsx\"")
	// responseBody.Data = "导出模板成功"
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
}

func (e *ExcelApi) ExcelCheck(c *gin.Context) {
	responseBody := &Response{Success: true}
	// var export excel.ExportType
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法获取文件"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法打开文件"})
		return
	}
	defer file.Close()
	f, err := excelize.OpenReader(file)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	importExcel, configParam, err := ExcelService.CheckExcel(f)
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = map[string]interface{}{
		"excelData":   importExcel,
		"configParam": configParam,
	}
	c.JSON(http.StatusOK, responseBody)
}

func (e *ExcelApi) ExcelImport(c *gin.Context) {
	responseBody := &Response{Success: true}
	var rawData Data
	if err := c.ShouldBindJSON(&rawData); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusBadRequest, responseBody)
		return
	}
	configParam := rawData.ConfigParam
	var items []*excels.Item
	for _, excelItem := range rawData.ExcelData {
		// 假设 product.Product 结构体与 ExcelData 结构体一致或可以转换
		item := excels.Item{
			BrandName:      excelItem.BrandName,
			BrandId:        excelItem.BrandId,
			TypeId:         excelItem.TypeId,
			TypeName:       excelItem.TypeName,
			ItemNumber:     excelItem.ItemNumber,
			SKU:            excelItem.SKU,
			SPU:            excelItem.SPU,
			Quantity:       excelItem.Quantity,
			Specifications: excelItem.Specifications,
			Barcode:        excelItem.Barcode,
			CustomsCode:    excelItem.CustomsCode,
			Description:    excelItem.Description,
			Cost:           excelItem.Cost,
			Color:          excelItem.Color,
			Options:        excelItem.Options,
		}
		items = append(items, &item)
	}
	err := ExcelService.ImportExcel(rawData.FileName, rawData.Descrption, items, configParam)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "更新成功"
	c.JSON(http.StatusOK, responseBody)
}

func (p *ExcelApi) SearchExcel(c *gin.Context) {
	responseBody := &Response{Success: true}
	var limit tool.RequestLimits
	var err error
	// 时间处理的通用方法
	parseTime := func(timeStr string) (*time.Time, error) {
		if timeStr == "" {
			return nil, nil
		}
		parsedTime, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return nil, err
		}
		return &parsedTime, nil
	}

	// StartCreateTime 处理
	var startCreateTime *time.Time
	startCreateTime, err = parseTime(c.Query("startCreateTime"))
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusBadRequest, responseBody)
		return
	}

	// EndCreateTime 处理
	var endCreateTime *time.Time
	endCreateTime, err = parseTime(c.Query("endCreateTime"))
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusBadRequest, responseBody)
		return
	}

	// StartUpdateTime 处理
	var startUpdateTime *time.Time
	startUpdateTime, err = parseTime(c.Query("startUpdateTime"))
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusBadRequest, responseBody)
		return
	}

	// EndUpdateTime 处理
	var endUpdateTime *time.Time
	endUpdateTime, err = parseTime(c.Query("endUpdateTime"))
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusBadRequest, responseBody)
		return
	}

	// 构建查询参数
	var params = &excels.ExcelQuery{
		StartCreateTime: startCreateTime,
		EndCreateTime:   endCreateTime,
		StartUpdateTime: startUpdateTime,
		EndUpdateTime:   endUpdateTime,
	}

	// 分页处理
	pageSize := c.Query("pageSize")
	pageNum := c.Query("pageNum")

	sPageNum, err := strconv.ParseInt(pageNum, 10, 64)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusBadRequest, responseBody)
		return
	}

	sPageSize, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusBadRequest, responseBody)
		return
	}

	limit.PageNum = int(sPageNum)
	limit.PageSize = int(sPageSize)

	// 查询
	excels, limits, err := excels.SearchExcelPages(limit, params)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusInternalServerError, responseBody)
		return
	}

	responseBody.Data = map[string]interface{}{
		"data":   excels,
		"limits": limits,
	}
	c.JSON(http.StatusOK, responseBody)
}

func (p *ExcelApi) SearchExcelCosts(c *gin.Context) {
	responseBody := &Response{Success: true}
	paramsExcelID := c.Query("excelId")
	var excelId uint
	if paramsExcelID != "" {
		if sexcelID, err := strconv.ParseUint(strings.TrimSpace(paramsExcelID), 10, 64); err != nil {
			isErr(err, responseBody)
			c.JSON(http.StatusBadRequest, responseBody)
			return
		} else {
			excelId = uint(sexcelID)
		}
	}
	rates, err := excels.SearchExcelCost(excelId)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = rates
	c.JSON(http.StatusOK, responseBody)
}

// func (e *ExcelApi) DeleteProd(c *gin.Context) {
// 	responseBody := &Response{Success: true}
// 	var prod product.Product
// 	err := c.ShouldBindJSON(&prod)
// 	if err != nil {
// 		isErr(err, responseBody)
// 		c.JSON(http.StatusNotFound, responseBody)
// 		return
// 	}
// 	err = productService.DeleteProduct(prod)
// 	if err != nil {
// 		isErr(err, responseBody)
// 		c.JSON(http.StatusNotFound, responseBody)
// 		return
// 	}
// 	responseBody.Data = "删除成功"
// 	c.JSON(http.StatusOK, responseBody)
// }

// func (e *ExcelApi) SearchProd(c *gin.Context) {
// 	responseBody := &Response{Success: true}
// 	var limit tool.RequestLimits
// 	pageSize := c.Query("pageSize")
// 	pageNum := c.Query("pageNum")
// 	sPageNum, err := strconv.ParseInt(pageNum, 10, 64)
// 	if err != nil {
// 		isErr(err, responseBody)
// 		c.JSON(http.StatusNotFound, responseBody)
// 		return
// 	}
// 	sPageSize, err := strconv.ParseInt(pageSize, 10, 64)
// 	if err != nil {
// 		isErr(err, responseBody)
// 		c.JSON(http.StatusNotFound, responseBody)
// 		return
// 	}
// 	limit.PageNum = int(sPageNum)
// 	limit.PageSize = int(sPageSize)
// 	products, limits, err := productService.SearchProductPages(limit)
// 	if err != nil {
// 		isErr(err, responseBody)
// 		c.JSON(http.StatusNotFound, responseBody)
// 		return
// 	}
// 	responseBody.Data = map[string]interface{}{
// 		"data":   products,
// 		"limits": limits,
// 	}
// 	c.JSON(http.StatusOK, responseBody)
// }

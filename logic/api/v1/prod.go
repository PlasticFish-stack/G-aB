package V1

import (
	"net/http"
	"project/logic"
	"project/logic/model/product"
	prod "project/logic/service/product"
	"project/logic/service/tool"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductApi struct{}

func (p *ProductApi) AddProd(c *gin.Context) {
	responseBody := &Response{Success: true}
	var prod product.Product
	err := c.ShouldBindJSON(&prod)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = productService.AddProduct(logic.Gorm, prod)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "添加成功"
	c.JSON(http.StatusOK, responseBody)
}

func (p *ProductApi) UpdateProd(c *gin.Context) {
	responseBody := &Response{Success: true}
	var prod product.Product
	err := c.ShouldBindJSON(&prod)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = productService.UpdateProduct(prod)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "更新成功"
	c.JSON(http.StatusOK, responseBody)
}

func (p *ProductApi) DeleteProd(c *gin.Context) {
	responseBody := &Response{Success: true}
	var prod product.Product
	err := c.ShouldBindJSON(&prod)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = productService.DeleteProduct(prod)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "删除成功"
	c.JSON(http.StatusOK, responseBody)
}

func (p *ProductApi) SearchProd(c *gin.Context) {
	responseBody := &Response{Success: true}
	var limit tool.RequestLimits
	var err error

	// BrandId 处理
	paramsBrandId := c.Query("brandId")
	var brandId *uint
	if paramsBrandId != "" {
		if sbrandId, err := strconv.ParseUint(strings.TrimSpace(paramsBrandId), 10, 64); err != nil {
			isErr(err, responseBody)
			c.JSON(http.StatusBadRequest, responseBody)
			return
		} else {
			uintVal := uint(sbrandId)
			brandId = &uintVal
		}
	}

	// TypeId 处理
	paramsTypeID := c.Query("typeId")
	var typeId *uint
	if paramsTypeID != "" {
		if stypeID, err := strconv.ParseUint(strings.TrimSpace(paramsTypeID), 10, 64); err != nil {
			isErr(err, responseBody)
			c.JSON(http.StatusBadRequest, responseBody)
			return
		} else {
			uintVal := uint(stypeID)
			typeId = &uintVal
		}
	}
	paramsItemNumber := c.Query("itemNumber")
	var itemNumber *string
	if paramsItemNumber != "" {
		itemNumber = &paramsItemNumber
	}
	paramsSku := c.Query("sku")
	var sku *string
	if paramsSku != "" {
		sku = &paramsSku
	}
	paramsSpu := c.Query("spu")
	var spu *string
	if paramsSpu != "" {
		spu = &paramsSpu
	}
	paramsBarcode := c.Query("barcode")
	var barcode *string
	if paramsBarcode != "" {
		barcode = &paramsBarcode
	}
	paramsCustomscode := c.Query("customscode")
	var customscode *string
	if paramsCustomscode != "" {
		customscode = &paramsCustomscode
	}
	paramsSpecifications := c.Query("specifications")
	var specifications *string
	if paramsSpecifications != "" {
		specifications = &paramsSpecifications
	}
	// Color 处理
	paramsColor := c.Query("color")
	var color *string
	if paramsColor != "" {
		color = &paramsColor
	}

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
	var params = &prod.ProductQuery{
		BrandID:         brandId,
		TypeID:          typeId,
		ItemNumber:      itemNumber,
		SKU:             sku,
		SPU:             spu,
		Barcode:         barcode,
		Customscode:     customscode,
		Specifications:  specifications,
		Color:           color,
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
	products, limits, err := productService.SearchProductPages(limit, params)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusInternalServerError, responseBody)
		return
	}

	responseBody.Data = map[string]interface{}{
		"data":   products,
		"limits": limits,
	}
	c.JSON(http.StatusOK, responseBody)
}
func (p *ProductApi) SearchOneProd(c *gin.Context) {
	responseBody := &Response{Success: true}
	paramsProductId := c.Query("productId")
	var productId uint
	if paramsProductId != "" {
		if sproductId, err := strconv.ParseUint(strings.TrimSpace(paramsProductId), 10, 64); err != nil {
			isErr(err, responseBody)
			c.JSON(http.StatusBadRequest, responseBody)
			return
		} else {
			productId = uint(sproductId)
		}
	}
	products, err := productService.SearchProduct(productId)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusInternalServerError, responseBody)
		return
	}

	responseBody.Data = map[string]interface{}{
		"data": products,
	}
	c.JSON(http.StatusOK, responseBody)
}

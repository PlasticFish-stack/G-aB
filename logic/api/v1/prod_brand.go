package V1

import (
	"net/http"
	"project/logic"
	"project/logic/model/product"

	"github.com/gin-gonic/gin"
)

func ProductBrandGetAll(c *gin.Context) {
	responseBody := &Response{Success: true}
	result, err := productService.SearchBrands()
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = result
	c.JSON(http.StatusOK, responseBody)
}

func ProductBrandAdd(c *gin.Context) {
	responseBody := &Response{Success: true}
	var prodBrand []product.Brand
	err := c.ShouldBindJSON(&prodBrand)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = productService.AddBrands(prodBrand)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "添加成功"
	c.JSON(http.StatusOK, responseBody)
}

func ProductBrandUpdate(c *gin.Context) {
	responseBody := &Response{Success: true}
	var prodBrand product.Brand
	err := c.ShouldBindJSON(&prodBrand)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = productService.UpdateBrands(prodBrand)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "更新成功"
	c.JSON(http.StatusOK, responseBody)
}

func ProductBrandDelete(c *gin.Context) {
	responseBody := &Response{Success: true}
	var prodBrand []product.Brand
	err := c.ShouldBindJSON(&prodBrand)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = productService.DeleteBrands(*logic.Gorm, prodBrand)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "删除成功"
	c.JSON(http.StatusOK, responseBody)
}

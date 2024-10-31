package V1

import (
	"net/http"
	"project/logic/controll"
	"project/logic/model"

	"github.com/gin-gonic/gin"
)

func ProductBrandGetAll(c *gin.Context) {
	responseBody := &Response{Success: true}
	result, err := controll.ProductBrandGetAll()
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = result
	c.JSON(http.StatusOK, responseBody)
}

func ProductBrandAdd(c *gin.Context) {
	responseBody := &Response{Success: true}
	var productBrand model.ProductBrand
	err := c.ShouldBindJSON(&productBrand)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = controll.ProductBrandAdd(
		productBrand.Name,
		productBrand.Description)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "添加成功"
	c.JSON(http.StatusOK, responseBody)
}

func ProductBrandUpdate(c *gin.Context) {
	responseBody := &Response{Success: true}
	var productBrand model.ProductBrand
	err := c.ShouldBindJSON(&productBrand)
	err = controll.ProductBrandUpdate(
		productBrand.Id,
		productBrand.Name,
		productBrand.Description)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "更新成功"
	c.JSON(http.StatusOK, responseBody)
}

func ProductBrandDelete(c *gin.Context) {
	responseBody := &Response{Success: true}
	var prodBrand model.ProductBrand
	err := c.ShouldBindJSON(&prodBrand)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = controll.ProductBrandDelete(prodBrand.Id, prodBrand.Name)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "删除成功"
	c.JSON(http.StatusOK, responseBody)
}

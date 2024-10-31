package V1

import (
	"net/http"
	"project/logic/controll"
	"project/logic/model"
	"project/logic/model/product"
	"project/logic/service/tool"

	"github.com/gin-gonic/gin"
)

type ProductTypeApi struct{}

func (p *ProductTypeApi) GetListProdType(c *gin.Context) {
	responseBody := &Response{Success: true}
	var limits tool.RequestLimits
	err := c.ShouldBindJSON(&limits)
	typeList, reslimits, err := productService.SearchTypeTree(limits)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = map[string]interface{}{
		"data":   typeList,
		"limits": reslimits,
	}
	c.JSON(http.StatusOK, responseBody)
}

func (p *ProductTypeApi) AddProdType(c *gin.Context) {
	responseBody := &Response{Success: true}
	var prodType product.Type
	err := c.ShouldBindJSON(&prodType)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = productService.AddType(prodType)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "添加成功"
	c.JSON(http.StatusOK, responseBody)
}

func (p *ProductTypeApi) UpdateProdType(c *gin.Context) {
	responseBody := &Response{Success: true}
	var prodType product.Type
	err := c.ShouldBindJSON(&prodType)
	err = controll.ProductTypeUpdate(
		productType.Id,
		productType.Name,
		productType.Description,
		productType.Sort,
		productType.ParentId,
		productType.Tax,
		productType.Field,
		productType.Formulas)
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

func ProductTypeDelete(c *gin.Context) {
	responseBody := &Response{Success: true}
	var prodType model.ProdType
	err := c.ShouldBindJSON(&prodType)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = controll.ProductTypeDelete(prodType.Id, prodType.Name)
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

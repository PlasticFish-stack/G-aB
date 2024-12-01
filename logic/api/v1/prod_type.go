package V1

import (
	"fmt"
	"net/http"
	"project/logic/model/product"

	"github.com/gin-gonic/gin"
)

type ProductTypeApi struct{}

func (p *ProductTypeApi) GetProdTypeList(c *gin.Context) {
	responseBody := &Response{Success: true}
	// var limits tool.RequestLimits
	// err := c.ShouldBindJSON(&limits)
	// typeList, reslimits, err := productService.SearchTypeTree(limits)
	typeList, err := productService.SearchTypeTree()
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = map[string]interface{}{
		"data": typeList,
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
	fmt.Println(prodType, "type")
	err = productService.UpdateType(prodType)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "更新成功"
	c.JSON(http.StatusOK, responseBody)
}

// func (p *ProductTypeApi) DeleteProdType(c *gin.Context) {
// 	responseBody := &Response{Success: true}
// 	var prodType product.Type
// 	err := c.ShouldBindJSON(&prodType)
// 	if err != nil {
// 		isErr(err, responseBody)
// 		c.JSON(http.StatusNotFound, responseBody)
// 		return
// 	}
// 	err = controll.ProductTypeDelete(prodType.Id, prodType.Name)
// 	if err != nil {
// 		isErr(err, responseBody)
// 		c.JSON(http.StatusNotFound, responseBody)
// 		return
// 	}
// 	responseBody.Data = "删除成功"
// 	c.JSON(http.StatusOK, responseBody)
// }

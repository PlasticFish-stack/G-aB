package V1

import (
	"net/http"
	"project/logic"
	"project/logic/model"

	"github.com/gin-gonic/gin"
)

type requesetFieldGroup struct {
	Fields []model.Field `json:"fields"`
	ApiId  uint          `json:"apiId"`
}

type FieldApi struct{}

func (f *FieldApi) FieldAdd(c *gin.Context) {
	responseBody := &Response{Success: true}
	var requeset requesetFieldGroup
	if err := c.ShouldBindJSON(&requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := systemService.AddFields(requeset.ApiId, requeset.Fields); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "添加成功"
	c.JSON(http.StatusOK, responseBody)
}

func (f *FieldApi) FieldUpdate(c *gin.Context) {
	responseBody := &Response{Success: true}
	var requeset model.Field
	if err := c.ShouldBindJSON(&requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := systemService.UpdateField(requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "更新成功"
	c.JSON(http.StatusOK, responseBody)
}

func (f *FieldApi) FieldDelete(c *gin.Context) {
	responseBody := &Response{Success: true}
	var requeset model.Field
	if err := c.ShouldBindJSON(&requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := systemService.DeleteField(logic.Gorm, requeset.Id); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "删除成功"
	c.JSON(http.StatusOK, responseBody)
}

func (f *FieldApi) FieldGetAll(c *gin.Context) {
	responseBody := &Response{Success: true}
	result, err := systemService.SearchFieldGroup()
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = result
	c.JSON(http.StatusOK, responseBody)
}

package V1

import (
	"net/http"
	"project/logic"
	"project/logic/model"

	"github.com/gin-gonic/gin"
)

//	type requesetApi struct {
//		Api    model.Api `json:"api"`
//		MenuId uint      `json:"menuId"`
//	}
type ApiApi struct{}

func (a *ApiApi) ApiAdd(c *gin.Context) {
	responseBody := &Response{Success: true}
	var requeset model.Api
	if err := c.ShouldBindJSON(&requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := systemService.AddApi(requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "添加成功"
	c.JSON(http.StatusOK, responseBody)
}

func (a *ApiApi) ApiUpdate(c *gin.Context) {
	responseBody := &Response{Success: true}
	var requeset model.Api
	if err := c.ShouldBindJSON(&requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := systemService.UpdateApi(requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "添加成功"
	c.JSON(http.StatusOK, responseBody)
}

func (a *ApiApi) ApiDelete(c *gin.Context) {
	responseBody := &Response{Success: true}
	var requeset model.Api
	if err := c.ShouldBindJSON(&requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := systemService.DeleteApi(logic.Gorm, requeset.Id); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "删除成功"
	c.JSON(http.StatusOK, responseBody)
}

func (a *ApiApi) ApiGetAll(c *gin.Context) {
	responseBody := &Response{Success: true}
	result, err := systemService.SearchApiGroup()
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = result
	c.JSON(http.StatusOK, responseBody)
}

package V1

import (
	"net/http"
	"project/logic/controll"
	"project/logic/model"

	"github.com/gin-gonic/gin"
)

func MenusGetAll(c *gin.Context) {
	responseBody := &Response{Success: true}
	result, err := controll.MenuGetAll()
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

func MenuAdd(c *gin.Context) {
	responseBody := &Response{Success: true}
	var menu model.Menu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = controll.MenuAdd(
		menu.Name,
		menu.Identifier,
		menu.Description,
		menu.Component,
		menu.Path,
		menu.Icon,
		menu.Sort,
		menu.ParentId,
	)
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
func MenuUpdate(c *gin.Context) {
	responseBody := &Response{Success: true}
	var menu model.Menu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	// jMenuId := c.PostForm("id")
	// jMenuName := c.PostForm("name")
	// jIdentifier := c.PostForm("identifier")
	// jDescription := c.PostForm("description")
	// jComponent := c.PostForm("component")
	// jPath := c.PostForm("path")
	// jIcon := c.PostForm("icon")
	// jSort := c.PostForm("sort")
	// jParentId := c.PostForm("parent_id")
	// jStatus := c.PostForm("status")
	// id, err := strconv.Atoi(jMenuId)
	// sort, err := strconv.Atoi(jSort)
	// parentid, err := strconv.Atoi(jParentId)
	// status, err := strconv.ParseBool(jStatus)
	err = controll.MenuUpdate(
		menu.Id,
		menu.Name,
		menu.Identifier,
		menu.Description,
		menu.Component,
		menu.Path,
		menu.Icon,
		menu.Sort,
		menu.ParentId,
		menu.Status,
	)
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

func MenuDelete(c *gin.Context) {
	responseBody := &Response{Success: true}
	var menu model.Menu
	err := c.ShouldBindJSON(&menu)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = controll.MenuDelete(menu.Id, menu.Name)
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

package V1

import (
	"fmt"
	"net/http"
	"project/logic/controll"
	"project/logic/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RoleAdd(c *gin.Context) {
	responseBody := &Response{Success: true}
	var role model.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = controll.RoleAdd(role.Name, role.Identifier, role.Description)
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

func RoleUpdate(c *gin.Context) {
	responseBody := &Response{Success: true}
	var role model.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = controll.RoleUpdate(role.Id, role.Name, role.Identifier, role.Description, role.Status)
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

func RoleDelete(c *gin.Context) {
	responseBody := &Response{Success: true}
	var role model.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = controll.RoleDelete(role.Id, role.Name)
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

func RoleBindMenu(c *gin.Context) {
	responseBody := &Response{Success: true}
	var bind struct {
		RoleId uint   `json:"roleId"`
		MenuId []uint `json:"menuId"`
	}
	err := c.ShouldBindJSON(&bind)
	fmt.Println(bind)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := controll.BindRoleMenu(bind.RoleId, bind.MenuId); err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusInternalServerError, responseBody)
		return
	}
	responseBody.Data = "绑定成功"
	c.JSON(http.StatusOK, responseBody)
}

// func RoleUnBindMenu(c *gin.Context) {
// 	responseBody := &Response{Success: true}
// 	var bind struct {
// 		RoleId uint   `json:"role_id"`
// 		MenuId []uint `json:"menu_id"`
// 	}
// 	err := c.ShouldBindJSON(&bind)
// 	if err != nil {
// 		responseBody.Success = false
// 		responseBody.Data = map[string]interface{}{
// 			"error": err.Error(),
// 		}
// 		c.JSON(http.StatusNotFound, responseBody)
// 		return
// 	}
// 	if err := controll.UnBindRoleMenu(bind.RoleId, bind.MenuId); err != nil {
// 		responseBody.Success = false
// 		responseBody.Data = map[string]interface{}{
// 			"error": err.Error(),
// 		}
// 		c.JSON(http.StatusInternalServerError, responseBody)
// 		return
// 	}
// 	responseBody.Data = "解绑成功"
// 	c.JSON(http.StatusOK, responseBody)
// }

func RoleGetAll(c *gin.Context) {
	responseBody := &Response{Success: true}
	result, err := controll.RoleGetAll()
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

func GetRoleBindMenu(c *gin.Context) {
	responseBody := &Response{Success: true}
	froleId := c.Query("roleId")
	roleId, err := strconv.Atoi(froleId)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	result, err := controll.GetBindRoleMenu(uint(roleId))
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

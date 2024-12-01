package V1

import (
	"net/http"
	"project/logic"
	"project/logic/controll"
	"project/logic/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

type requesetRoleBindMenu struct {
	RoleId  uint   `json:"roleId"`
	MenuIds []uint `json:"menuIds"`
}

type requesetRoleBindApiAndField struct {
	RoleId uint            `json:"roleId"`
	ApiIds map[uint][]uint `json:"apiIds"`
}

type RoleApi struct{}

func (r *RoleApi) Add(c *gin.Context) {
	responseBody := &Response{Success: true}
	var requeset model.Role
	if err := c.ShouldBindJSON(&requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := systemService.AddRole(
		logic.Gorm,
		requeset.Name,
		requeset.Identifier,
		requeset.Description,
	); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "添加成功"
	c.JSON(http.StatusOK, responseBody)
}

func (r *RoleApi) Update(c *gin.Context) {
	responseBody := &Response{Success: true}
	var requeset model.Role
	if err := c.ShouldBindJSON(&requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := systemService.UpdateRole(
		logic.Gorm,
		requeset.Id,
		requeset.Name,
		requeset.Identifier,
		requeset.Description,
		requeset.Status,
	); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "更新成功"
	c.JSON(http.StatusOK, responseBody)
}

func (r *RoleApi) Delete(c *gin.Context) {
	responseBody := &Response{Success: true}
	var requeset requesetRoleBindMenu
	if err := c.ShouldBindJSON(&requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := controll.RoleDelete(requeset.RoleId); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "删除成功"
	c.JSON(http.StatusOK, responseBody)
}

func (r *RoleApi) BindMenu(c *gin.Context) {
	responseBody := &Response{Success: true}
	var requeset requesetRoleBindMenu
	if err := c.ShouldBindJSON(&requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := systemService.BindMenuToRole(logic.Gorm, requeset.RoleId, requeset.MenuIds); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "绑定成功"
	c.JSON(http.StatusOK, responseBody)
}

func (r *RoleApi) BindApiField(c *gin.Context) {
	responseBody := &Response{Success: true}
	var requeset requesetRoleBindApiAndField
	if err := c.ShouldBindJSON(&requeset); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := systemService.BindApiAndFieldToRole(logic.Gorm, requeset.RoleId, requeset.ApiIds); err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "绑定成功"
	c.JSON(http.StatusOK, responseBody)
}

func (r *RoleApi) GetGroup(c *gin.Context) {
	responseBody := &Response{Success: true}
	result, err := controll.RoleGetAll()
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = result
	c.JSON(http.StatusOK, responseBody)
}

func (r *RoleApi) GetBindMenu(c *gin.Context) {
	responseBody := &Response{Success: true}
	queryRoleId := c.Query("roleId")
	roleId, err := strconv.Atoi(queryRoleId)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	result, err := systemService.GetRoleBindMenu(logic.Gorm, uint(roleId))
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = result
	c.JSON(http.StatusOK, responseBody)
}

func (r *RoleApi) GetBindApiField(c *gin.Context) {
	responseBody := &Response{Success: true}
	queryRoleId := c.Query("roleId")
	roleId, err := strconv.Atoi(queryRoleId)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	result, err := systemService.GetRoleBindApiField(logic.Gorm, uint(roleId))
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = result
	c.JSON(http.StatusOK, responseBody)
}

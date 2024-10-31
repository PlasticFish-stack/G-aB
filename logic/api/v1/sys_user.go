package V1

import (
	"fmt"
	"net/http"
	"project/logic/controll"
	"project/logic/model"
	"project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	// Duration string      `json:"duration"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func Login(c *gin.Context) {
	responseBody := &Response{Success: true}
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体格式不正确"})
		return
	}
	ip := c.ClientIP()
	result, err := controll.Login(request.Username, request.Password, ip)
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

func Register(c *gin.Context) {
	responseBody := &Response{Success: true}
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	// jUsername := c.PostForm("username")
	// jPassword := c.PostForm("password")
	// jPicture := c.PostForm("picture")
	// jNickname := c.PostForm("nickname")
	hashPassword, _ := utils.HashExec(user.Password)
	err = controll.Register(user.Name, hashPassword, user.Picture, user.Nickname)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "注册成功"
	c.JSON(http.StatusOK, responseBody)
}

func Refresh(c *gin.Context) {
	responseBody := &Response{Success: true}
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体格式不正确"})
		return
	}
	result, err := controll.Refresh(request.RefreshToken)
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

func UserGetAll(c *gin.Context) {
	responseBody := &Response{Success: true}
	result, err := controll.UserGetAll()
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

func UserBindRole(c *gin.Context) {
	responseBody := &Response{Success: true}
	var bind struct {
		UserId uint   `json:"userId"`
		RoleId []uint `json:"roleId"`
	}
	err := c.ShouldBindJSON(&bind)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := controll.BindUserRole(bind.UserId, bind.RoleId); err != nil {
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

func GetUserBindRole(c *gin.Context) {
	responseBody := &Response{Success: true}
	fuserId := c.Query("userId")
	userId, err := strconv.Atoi(fuserId)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	result, err := controll.GetBindUserRole(uint(userId))
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

func RoutesGet(c *gin.Context) {
	responseBody := &Response{Success: true}
	cUserid, _ := c.Get("userid")
	userid, ok := cUserid.(uint)
	if !ok {
		responseBody.Success = false
		responseBody.Data = fmt.Errorf("userId转换uint格式失败,请检查是否为number格式")
		c.JSON(http.StatusInternalServerError, responseBody)
		return
	}
	menus, err := controll.GetRoutes(userid)
	responseBody.Data = menus
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusInternalServerError, responseBody)
		return
	}
	c.JSON(http.StatusOK, responseBody)
}

func UserDelete(c *gin.Context) {
	responseBody := &Response{Success: true}
	var user struct {
		Id   uint   `json:"id"`
		Name string `json:"name"`
	}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := controll.UserDelete(user.Id, user.Name); err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusInternalServerError, responseBody)
		return
	}
	responseBody.Data = "删除成功"
	c.JSON(http.StatusOK, responseBody)
}

func UserUpdate(c *gin.Context) {
	responseBody := &Response{Success: true}
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	if err := controll.UserUpdate(user.Id, user.Nickname, user.Picture, user.Status); err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusInternalServerError, responseBody)
		return
	}
	responseBody.Data = "更新成功"
	c.JSON(http.StatusOK, responseBody)
}

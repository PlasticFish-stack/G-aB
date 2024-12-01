package V1

import (
	"net/http"
	"project/logic/model/rate"

	"github.com/gin-gonic/gin"
)

type RateApi struct{}

func (r *RateApi) RateGet(c *gin.Context) {
	responseBody := &Response{Success: true}
	rates, err := RateService.RateGet()
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = rates
	c.JSON(http.StatusOK, responseBody)
}

func (r *RateApi) RateApiUpdate(c *gin.Context) {
	responseBody := &Response{Success: true}
	_, err := RateService.RataApiUpdate()
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "更新成功"
	c.JSON(http.StatusOK, responseBody)
}

func (r *RateApi) RateUpdate(c *gin.Context) {
	responseBody := &Response{Success: true}
	var rate rate.Rate
	err := c.ShouldBindJSON(&rate)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = RateService.RateUpdate(rate)
	if err != nil {
		isErr(err, responseBody)
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = "更新成功"
	c.JSON(http.StatusOK, responseBody)
}

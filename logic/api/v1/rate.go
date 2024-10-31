package V1

import (
	"net/http"
	"project/logic/controll"
	"project/logic/model"

	"github.com/gin-gonic/gin"
)

func RateGet(c *gin.Context) {
	responseBody := &Response{Success: true}
	rates, err := controll.RateGet()
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	responseBody.Data = rates
	c.JSON(http.StatusOK, responseBody)
}

func RateApiUpdate(c *gin.Context) {
	responseBody := &Response{Success: true}
	_, err := controll.RataApiUpdate()
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

func RateUpdate(c *gin.Context) {
	responseBody := &Response{Success: true}
	var rate model.Rate
	err := c.ShouldBindJSON(&rate)
	if err != nil {
		responseBody.Success = false
		responseBody.Data = map[string]interface{}{
			"error": err.Error(),
		}
		c.JSON(http.StatusNotFound, responseBody)
		return
	}
	err = controll.RateUpdate(rate.CurrencyName, rate.DescriptionEn, rate.DescriptionCn, rate.CountryIcon, rate.Sort)
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

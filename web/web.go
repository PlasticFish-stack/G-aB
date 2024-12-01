package web

import (
	"net/http"
	V1 "project/logic/api/v1"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Response struct {
	// Duration string      `json:"duration"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

var api = V1.ApiGroup{}

func basicRouter(router *gin.Engine) {
	private := router.Group("/")
	{
		private.POST("/login", V1.Login)
		private.GET("/get-async-routes", V1.RoutesGet)
		private.POST("/refresh", V1.Refresh)
	}
}

func userRouter(router *gin.Engine) {
	private := router.Group("/user")
	{
		private.GET("/getall", V1.UserGetAll)
		private.POST("/bind", V1.UserBindRole)
		private.POST("/delete", V1.UserDelete)
		private.POST("/add", V1.Register)
		private.GET("/getbind", V1.GetUserBindRole)
		private.POST("/update", V1.UserUpdate)
	}
}
func roleRouter(router *gin.Engine) {
	private := router.Group("/role")
	{
		private.POST("/add", roleApi.Add)
		private.POST("/update", roleApi.Update)
		private.DELETE("/delete", roleApi.Delete)
		private.POST("/bind-menu", roleApi.BindMenu)
		private.POST("/bind-api-field", roleApi.BindApiField)
		private.GET("/get-all", roleApi.GetGroup)
		private.GET("/get-bind-menu", roleApi.GetBindMenu)
		private.GET("/get-bind-api-field", roleApi.GetBindMenu)
	}
}
func apiRouter(router *gin.Engine) {
	private := router.Group("/api")
	{
		private.POST("/add", api.Add)
		private.POST("/update", roleApi.Update)
		private.DELETE("/delete", roleApi.Delete)
		private.POST("/bind-menu", roleApi.BindMenu)
		private.POST("/bind-api-field", roleApi.BindApiField)
		private.GET("/get-all", roleApi.GetGroup)
		private.GET("/get-bind-menu", roleApi.GetBindMenu)
		private.GET("/get-bind-api-field", roleApi.GetBindMenu)
	}
}
func fieldRouter(router *gin.Engine) {
	private := router.Group("/api")
	{
		private.POST("/add", roleApi.Add)
		private.POST("/update", roleApi.Update)
		private.DELETE("/delete", roleApi.Delete)
		private.POST("/bind-menu", roleApi.BindMenu)
		private.POST("/bind-api-field", roleApi.BindApiField)
		private.GET("/get-all", roleApi.GetGroup)
		private.GET("/get-bind-menu", roleApi.GetBindMenu)
		private.GET("/get-bind-api-field", roleApi.GetBindMenu)
	}
}
func menuRouter(router *gin.Engine) {
	private := router.Group("/menu")
	{
		private.GET("/getall", V1.MenusGetAll)
		private.POST("/add", V1.MenuAdd)
		private.POST("/update", V1.MenuUpdate)
		private.POST("/delete", V1.MenuDelete)
	}
}
func productRouter(router *gin.Engine) {
	private := router.Group("/product/product")
	{
		private.POST("/add", productApi.AddProd)
		private.POST("/update", productApi.UpdateProd)
		private.POST("/delete", productApi.DeleteProd)
		private.GET("/getlimits", productApi.SearchProd)
		private.GET("/get", productApi.SearchOneProd)
	}

}
func productTypeRouter(router *gin.Engine) {
	private := router.Group("/product/type")
	{
		private.POST("/add", productTypeApi.AddProdType)
		private.POST("/update", productTypeApi.UpdateProdType)
		// private.POST("/delete", api.DeleteProdType)
		private.GET("/getall", productTypeApi.GetProdTypeList)
	}
}
func productBrandRouter(router *gin.Engine) {
	private := router.Group("/product/brand")
	{
		private.POST("/add", V1.ProductBrandAdd)
		private.POST("/update", V1.ProductBrandUpdate)
		private.POST("/delete", V1.ProductBrandDelete)
		private.GET("/getall", V1.ProductBrandGetAll)
	}
}
func RateRouter(router *gin.Engine) {
	private := router.Group("/rate")
	{
		private.GET("/get", rateApi.RateGet)
		private.POST("/getapi", rateApi.RateApiUpdate)
		private.POST("/update", rateApi.RateUpdate)
	}
}

func ExcelRouter(router *gin.Engine) {
	private := router.Group("/excel")
	{
		private.POST("/export", api.ExcelExport)
		private.POST("/check", api.ExcelCheck)
		private.POST("/import", api.ExcelImport)
		private.GET("/getlimits", api.SearchExcel)
		private.GET("/getcostpages", api.SearchExcelCosts)
	}
}
func Start() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(JWTMiddleware())
	basicRouter(router)
	roleRouter(router)
	userRouter(router)
	menuRouter(router)
	productTypeRouter(router)
	productBrandRouter(router)
	productRouter(router)
	RateRouter(router)
	ExcelRouter(router)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found"})
	})
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Run(":8080")
}

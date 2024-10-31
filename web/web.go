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
		private.POST("/add", V1.RoleAdd)
		private.POST("/update", V1.RoleUpdate)
		private.POST("/delete", V1.RoleDelete)
		private.GET("/getall", V1.RoleGetAll)
		private.POST("/bind", V1.RoleBindMenu)
		private.GET("/getbind", V1.GetRoleBindMenu)
		// private.POST("/unbind", V1.RoleUnBindMenu)
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

func productTypeRouter(router *gin.Engine) {
	private := router.Group("/product/type")
	{
		private.POST("/add", api.AddProdType)
		private.POST("/update", V1.ProductTypeUpdate)
		private.POST("/delete", V1.ProductTypeDelete)
		private.GET("/getall", api.GetListProdType)
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
		private.GET("/get", V1.RateGet)
		private.POST("/getapi", V1.RateApiUpdate)
		private.POST("/update", V1.RateUpdate)
	}
}

//	func productInformationRouter(router *gin.Engine) {
//		private := router.Group("/product/information")
//		{
//			private.POST("/add", V1.AddRole)
//			private.POST("/update", V1.UpdateRole)
//			private.POST("/delete", V1.DeleteRole)
//			private.GET("/getall", V1.GetAllRole)
//		}
//	}
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
	RateRouter(router)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Page not found"})
	})
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.Run(":8080")
}

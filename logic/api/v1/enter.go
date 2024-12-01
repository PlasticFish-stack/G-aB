package V1

import (
	"project/logic/service"
)

type ApiGroup struct {
	ProductTypeApi
	ProductApi
	ExcelApi
	RateApi
	RoleApi
	ApiApi
	FieldApi
}

var ApiGroupApp = new(ApiGroup)

var (
	productService = service.ServiceGroupApp.ProductServiceGroup
	ExcelService   = service.ServiceGroupApp.ExcelServiceGroup
	RateService    = service.ServiceGroupApp.RateServiceGroup
	systemService  = service.ServiceGroupApp.SystemServiceGroup
)

func isErr(err error, body *Response) {
	if err != nil {
		body.Success = false
		body.Data = map[string]interface{}{
			"error": err.Error(),
		}
		return
	}
}

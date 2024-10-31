package V1

import (
	"project/logic/service/product"
)

type ApiGroup struct {
	ProductTypeApi
}

var (
	productService = product.ServiceProductGroupApp
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

package service

import "project/logic/service/product"

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	ProductServiceGroup product.ServiceProductGroup
}

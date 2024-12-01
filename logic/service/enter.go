package service

import (
	"project/logic/service/excel"
	"project/logic/service/product"
	"project/logic/service/rate"
	"project/logic/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	ProductServiceGroup product.ServiceProductGroup
	ExcelServiceGroup   excel.ServiceExcelGroup
	RateServiceGroup    rate.ServiceRateGroup
	SystemServiceGroup  system.ServiceSystemGroup
}

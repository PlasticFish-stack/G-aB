package main

import (
	"log"
	_ "project/console"
	"project/console/logger"
	"project/logic"
	"project/logic/model"
	"project/logic/model/excel"
	"project/logic/model/product"
	"project/logic/model/rate"
	"project/web"
)

func main() {
	// utils.CloseDB()
	logger.InitLogger(false)
	logic.GetConn()
	err := logic.Gorm.AutoMigrate(
		&model.Field{},
		&model.User{},
		&model.Role{},
		&model.Menu{},
		// &model.ProductBrand{},

		// &model.Product{},
		// &model.ProductCost{},
		&rate.Rate{},
		&product.Brand{},

		&product.Type{},
		&product.Field{},
		&excel.ExcelLog{},

		&product.Cost{},
		&product.Formula{},
		&product.Product{},
		// &model.ProdType{},
		// &model.ProdTypeField{},
		// &model.ProdTypeFormula{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	web.Start()
	defer logger.Sync()
}

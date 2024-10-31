package main

import (
	"log"
	"project/logic"
	"project/logic/model"
	"project/web"
)

func main() {
	// utils.CloseDB()
	logic.GetConn()
	err := logic.Gorm.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.ProductBrand{},
		// &model.ExcelLog{},
		&model.Product{},
		&model.ProductCost{},
		&model.Rate{},
		&model.ProdType{},
		&model.ProdTypeField{},
		&model.ProdTypeFormula{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	web.Start()

}

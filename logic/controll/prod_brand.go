package controll

import (
	"fmt"
	"project/logic"
	"project/logic/model"
)

func ProductBrandGetAll() ([]model.ProductBrand, error) {
	var productBrand []model.ProductBrand
	var err error
	if productBrand, err = model.SearchProductBrand(logic.Gorm); err != nil {
		return nil, err
	}
	return productBrand, nil
}

func ProductBrandAdd(name string, description string) error {
	productBrand := &model.ProductBrand{
		Name:        name,
		Description: description,
	}
	err := productBrand.Add(logic.Gorm)
	if err != nil {
		return err
	}
	return nil
}

func ProductBrandUpdate(id uint, name string, description string) error {
	productBrand := &model.ProductBrand{
		Global:      model.Global{Id: id},
		Name:        name,
		Description: description,
	}
	if err := productBrand.Update(logic.Gorm); err != nil {
		return err
	}
	return nil
}

func ProductBrandDelete(id uint, name string) error {
	var productBrand = model.ProductBrand{
		Global: model.Global{Id: id},
	}
	resultBrand, err := productBrand.Search(logic.Gorm)
	if err != nil {
		return err
	}
	if name != resultBrand.Name {
		return fmt.Errorf("品牌与id不匹配")
	}
	if err := resultBrand.Delete(logic.Gorm); err != nil {
		return err
	}
	return nil
}

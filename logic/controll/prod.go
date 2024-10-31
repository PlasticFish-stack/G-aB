package controll

import (
	"fmt"
	"project/logic"
	"project/logic/model"
)

func ProdAdd(itemnumber string, brandId uint, sku string, spu string, quantity uint64, specifications string, barcode string, customscode string, description string, color string) error {
	prod := &model.Product{
		ItemNumber:     itemnumber,
		BrandId:        brandId,
		Sku:            sku,
		Spu:            spu,
		Quantity:       quantity,
		Specifications: specifications,
		Barcode:        barcode,
		Customscode:    customscode,
		Description:    description,
		Color:          color,
	}
	err := prod.Add(logic.Gorm)
	if err != nil {
		return err
	}
	return nil
}

func ProdUpdate(id uint, itemnumber string, brandId uint, sku string, spu string, quantity uint64, specifications string, barcode string, customscode string, description string, color string) error {
	updateInfo := model.Product{
		Global:         model.Global{Id: id},
		ItemNumber:     itemnumber,
		BrandId:        brandId,
		Sku:            sku,
		Spu:            spu,
		Quantity:       quantity,
		Specifications: specifications,
		Barcode:        barcode,
		Customscode:    customscode,
		Description:    description,
		Color:          color,
	}
	if err := updateInfo.Update(logic.Gorm); err != nil {
		return err
	}
	return nil
}

func ProdDelete(id uint, itemnumber string) error {
	var prod = model.Product{
		Global: model.Global{Id: id},
	}
	resultProd, err := prod.Search(logic.Gorm)
	if err != nil {
		return err
	}
	if itemnumber != resultProd.ItemNumber {
		return fmt.Errorf("菜单与id不匹配")
	}
	if err := resultProd.Delete(logic.Gorm); err != nil {
		return err
	}
	return nil
}

func ProdGetAll() ([]model.Product, error) {
	var products []model.Product
	var err error
	if products, err = model.SearchAllProduct(logic.Gorm); err != nil {
		return nil, err
	}
	return products, nil
}

package product

import (
	"errors"
	"fmt"
	"project/logic"
	"project/logic/model/product"

	"gorm.io/gorm"
)

func (serviceProduct *ServiceProductGroup) SearchBrands(db gorm.DB) ([]product.Brand, error) {
	var brands []product.Brand
	err := db.Order("name").Find(&brands).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("未查询到相关Fields: %v", err)
		}
		return nil, fmt.Errorf("查询Fields失败: %v", err)
	}
	return brands, nil
}

func (serviceProduct *ServiceProductGroup) SearchBrand(name string) (*product.Brand, error) {
	var brand product.Brand
	err := logic.Gorm.Where("name = ?", name).First(&brand).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("未查询到相关Fields: %v", err)
		}
		return nil, fmt.Errorf("查询Fields失败: %v", err)
	}
	return &brand, nil
}

func (serviceProduct *ServiceProductGroup) AddBrands(brands []product.Brand) error {
	for _, brand := range brands {
		findBrand, _ := ServiceProductGroupApp.SearchBrand(brand.Name)
		if findBrand != nil {
			return fmt.Errorf("此Type: %v的这条Brand已经创建: %v", findBrand.Name)
		}
	}
	err := logic.Gorm.Create(&brands).Error
	if err != nil {
		return fmt.Errorf("插入Brands失败: %v", err)
	}
	return nil
}

func (serviceProduct *ServiceProductGroup) UpdateBrands(brands []product.Brand) error {
	err := logic.Gorm.Updates(brands).Error
	if err != nil {
		return fmt.Errorf("更新Brands失败: %v", err)
	}
	return nil
}

func (serviceProduct *ServiceProductGroup) DeleteBrands(db gorm.DB, brands []product.Brand) error {
	err := db.Delete(brands).Error
	if err != nil {
		return fmt.Errorf("更新Brands失败: %v", err)
	}
	return nil
}

package product

import (
	"errors"
	"fmt"
	"project/logic"
	"project/logic/model/product"
	"project/logic/service/tool"

	"gorm.io/gorm"
)

func (serviceProduct *ServiceProductGroup) SearchCost(productId uint, limitCount uint) ([]product.Cost, error) {
	var costs []product.Cost
	err := logic.Gorm.Order("created_at desc").Limit(int(limitCount)).Find(&costs).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("未查询到相关Costs: %v", err)
		}
		return nil, fmt.Errorf("查询Costs失败: %v", err)
	}
	return costs, nil
}

func (serviceProduct *ServiceProductGroup) SearchCosts(limits tool.RequestLimits) (costs []*product.Cost, formatLimits *tool.ResponseLimits, err error) {
	var total int64
	offset, err := limits.GetOffset()
	if err != nil {
		return nil, nil, err
	}
	err = logic.Gorm.
		Model(&product.Cost{}).
		Count(&total).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limits.PageSize).
		Find(&costs).Error
	if err != nil {
		return nil, nil, fmt.Errorf("查询成本价表失败: %v", err)
	}
	return
}

func (serviceProduct *ServiceProductGroup) UpdateCostsDwMsg(ExcelLogId uint, costID uint, newDwPrice float64, newDwSales float64) error {
	return logic.Gorm.Transaction(func(tx *gorm.DB) error {
		var cost product.Cost
		if err := tx.First(&cost, costID).Error; err != nil {
			return fmt.Errorf("获取cost失败: %w", err)
		}
		if err := tx.Model(&product.Cost{}).Where("id = ?", costID).Updates(product.Cost{
			DwPrice: newDwPrice,
			DwSales: newDwSales,
		}).Error; err != nil {
			return fmt.Errorf("更新cost失败: %w", err)
		}
		if err := tx.Model(&product.Product{}).Where("id = ?", cost.ProductID).Updates(product.Product{
			DwPrice: newDwPrice,
			DwSales: newDwSales,
		}).Error; err != nil {
			return fmt.Errorf("更新Product的得物相关属性失败: %w", err)
		}
		return nil
	})
}

func (serviceProduct *ServiceProductGroup) AddCosts(costs []*product.Cost) error {
	err := logic.Gorm.Create(&costs).Error
	if err != nil {
		return fmt.Errorf("增加成本价失败,请检查: %v", err)
	}
	return nil
}

// func (serviceProduct *ServiceProductGroup) SearchBrandId(id uint) (*product.Brand, error) {
// 	var brand product.Brand
// 	err := logic.Gorm.First(&brand, id).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, fmt.Errorf("未查询到相关Brand: %v", err)
// 		}
// 		return nil, fmt.Errorf("查询Brand失败: %v", err)
// 	}
// 	return &brand, nil
// }

// func (serviceProduct *ServiceProductGroup) AddBrands(brands []product.Brand) error {
// 	for _, brand := range brands {
// 		findBrand, _ := ServiceProductGroupApp.SearchBrand(brand.Name)
// 		if findBrand != nil {
// 			return fmt.Errorf("此Type: %v的这条Brand已经创建: %v", findBrand.Name)
// 		}
// 	}
// 	err := logic.Gorm.Create(&brands).Error
// 	if err != nil {
// 		return fmt.Errorf("插入Brands失败: %v", err)
// 	}
// 	return nil
// }

// func (serviceProduct *ServiceProductGroup) UpdateBrands(brands product.Brand) error {
// 	var brand product.Brand
// 	err := logic.Gorm.First(&brand, brands.Id).Error
// 	if err != nil {
// 		return fmt.Errorf("没有查询到该Brand: %v", err)
// 	}

// 	findBrand, err := ServiceProductGroupApp.SearchBrand(brands.Name)
// 	if err == nil && findBrand.Id != brands.Id {
// 		return fmt.Errorf("该Brand名称已存在")
// 	}
// 	err = logic.Gorm.Updates(&brands).Error
// 	if err != nil {
// 		return fmt.Errorf("更新Brand失败: %v", err)
// 	}
// 	return nil
// }

// func (serviceProduct *ServiceProductGroup) DeleteBrands(db gorm.DB, brands []product.Brand) error {
// 	err := db.Delete(brands).Error
// 	if err != nil {
// 		return fmt.Errorf("更新Brands失败: %v", err)
// 	}
// 	return nil
// }

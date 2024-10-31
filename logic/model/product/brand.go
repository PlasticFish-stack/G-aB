package product

import (
	"project/logic/model"
)

type Brand struct {
	model.Global
	Name        string    `gorm:"unique" json:"name"`
	Description string    `json:"description"`
	Products    []Product `gorm:"foreignKey:BrandId"`
}

// func SearchProductBrand(db *gorm.DB) ([]ProductBrand, error) {
// 	var brands []ProductBrand
// 	if err := db.Find(&brands).Error; err != nil {
// 		return nil, fmt.Errorf("查询产品品牌失败: %v", err)
// 	}
// 	return brands, nil
// }

// func (p *ProductBrand) Search(db *gorm.DB) (*ProductBrand, error) {
// 	var prodBrand ProductBrand
// 	if err := db.First(&prodBrand, p.Id).Error; err != nil {
// 		return nil, fmt.Errorf("查询产品品牌失败: %v", err)
// 	}
// 	return &prodBrand, nil
// }

// func (p *ProductBrand) Add(db *gorm.DB) error {
// 	if err := db.Create(p).Error; err != nil {
// 		if errors.Is(err, gorm.ErrDuplicatedKey) {
// 			return fmt.Errorf("产品品牌名称已存在: %v", err)
// 		}
// 		return fmt.Errorf("新建产品品牌失败: %v", err)
// 	}
// 	return nil
// }

// func (p *ProductBrand) Update(db *gorm.DB) error {
// 	var resultProductBrand ProductBrand
// 	if err := db.Where(p.Id).First(&resultProductBrand).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return fmt.Errorf("未查询到该产品品牌: %v", err)
// 		}
// 		return fmt.Errorf("查询产品品牌失败: %v", err)
// 	}
// 	if err := db.Model(&resultProductBrand).Updates(&p).Error; err != nil {
// 		return fmt.Errorf("更新产品品牌失败,请检查: %v", err)
// 	}
// 	return nil
// }

// func (p *ProductBrand) Delete(db *gorm.DB) error {
// 	tx := db.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			panic(r)
// 		} else {
// 			tx.Commit()
// 		}
// 	}()
// 	currentTime := time.Now().Format("2006-01-02 15:04:05")
// 	p.Name = p.Name + "_is_deleted" + currentTime
// 	if err := tx.Updates(&p).Error; err != nil {
// 		return fmt.Errorf("删除产品品牌失败,请检查: %v", err)
// 	}
// 	if err := tx.Delete(&p).Error; err != nil {
// 		return fmt.Errorf("删除产品品牌失败,请检查: %v", err)
// 	}
// 	return nil
// }

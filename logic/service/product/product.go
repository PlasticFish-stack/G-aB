package product

import (
	"encoding/json"
	"fmt"
	"project/logic"
	"project/logic/model/product"
	"project/logic/service/tool"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductQuery struct {
	ItemNumber      *string    `json:"itemNumber"`
	BrandID         *uint      `json:"brandId"`
	TypeID          *uint      `json:"typeId"`
	SKU             *string    `json:"sku"`
	SPU             *string    `json:"spu"`
	Specifications  *string    `json:"specifications"`
	Barcode         *string    `json:"barcode"`
	Customscode     *string    ` json:"customscode"`
	Color           *string    `json:"color"`
	StartCreateTime *time.Time `json:"startCreateTime"`
	EndCreateTime   *time.Time `json:"endCreateTime"`
	StartUpdateTime *time.Time `json:"startUpdateTime"`
	EndUpdateTime   *time.Time `json:"endUpdateTime"`
}

func (q *ProductQuery) BuildQuery() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if q.BrandID != nil {
			db = db.Where("brand_id = ?", q.BrandID)
		}

		if q.TypeID != nil {
			db = db.Where("type_id = ?", q.TypeID)
		}
		if q.ItemNumber != nil {
			db = db.Where("item_number LIKE ?", "%"+*q.ItemNumber+"%")
		}
		if q.SKU != nil {
			db = db.Where("sku LIKE ?", "%"+*q.SKU+"%")
		}
		if q.SPU != nil {
			db = db.Where("spu LIKE ?", "%"+*q.SPU+"%")
		}
		if q.Barcode != nil {
			db = db.Where("barcode LIKE ?", "%"+*q.Barcode+"%")
		}
		if q.Customscode != nil {
			db = db.Where("customscode LIKE ?", "%"+*q.Customscode+"%")
		}
		if q.Color != nil {
			db = db.Where("color LIKE ?", "%"+*q.Color+"%")
		}
		if q.Specifications != nil {
			db = db.Where("specifications LIKE ?", "%"+*q.Specifications+"%")
		}

		if q.StartCreateTime != nil {
			db = db.Where("created_at >= ?", q.StartCreateTime)
		}

		if q.EndCreateTime != nil {
			db = db.Where("created_at <= ?", q.EndCreateTime)
		}

		if q.StartUpdateTime != nil {
			db = db.Where("updated_at >= ?", q.StartUpdateTime)
		}

		if q.EndUpdateTime != nil {
			db = db.Where("updated_at <= ?", q.EndUpdateTime)
		}

		return db
	}
}
func (serviceProduct *ServiceProductGroup) SearchProductPages(limits tool.RequestLimits, query *ProductQuery) (products []*product.Product, formatLimits *tool.ResponseLimits, err error) {
	var total int64
	offset, err := limits.GetOffset()
	if err != nil {
		return nil, nil, err
	}
	err = logic.Gorm.
		Model(&product.Product{}).
		Scopes(query.BuildQuery()).
		Order("id ASC").
		Preload(clause.Associations).
		Count(&total).
		Offset(offset).
		Limit(limits.PageSize).
		Find(&products).Error
	if err != nil {
		return nil, nil, fmt.Errorf("查询产品失败: %v", err)
	}
	formatLimits = tool.NewLimits(total, limits.PageSize, limits.PageNum)
	return
}

func (serviceProduct *ServiceProductGroup) SearchProduct(pid uint) (*product.Product, error) {
	var product *product.Product
	if err := logic.Gorm.Preload(clause.Associations).First(&product, pid).Error; err != nil {
		return nil, fmt.Errorf("查询产品失败: %v", err)
	}
	return product, nil
}

func (serviceProduct *ServiceProductGroup) AddProduct(db *gorm.DB, p product.Product) error {
	tx := logic.Gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	var product product.Product
	fields, err := serviceProduct.SearchFields(*tx, p.TypeId)
	if err != nil {
		tx.Rollback()
		return err
	}
	var requestOptions = make(map[string]string)
	if err = json.Unmarshal(p.Options.Bytes, &requestOptions); err != nil {
		return fmt.Errorf("格式化Options错误,请检查json字段是否为字符串类型: %v", err)
	}
	for optionKey := range requestOptions {
		var isInField = false
		for _, field := range fields {
			if optionKey == field.Name {
				isInField = true
				break
			}
		}
		if !isInField {
			return fmt.Errorf("Options值错误,请检查是否为对应产品类型字段")
		}
	}
	if err := tx.Where("item_number = ?", p.ItemNumber).First(&product).Error; err != nil {
		if err := tx.Create(&p).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("新建产品失败: %v", err)
		}
	} else {
		tx.Rollback()
		return fmt.Errorf("已存在相同货号产品")
	}
	return tx.Commit().Error
}

func (serviceProduct *ServiceProductGroup) UpdateProduct(p product.Product) error {
	tx := logic.Gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	fields, err := serviceProduct.SearchFields(*tx, p.TypeId)
	if err != nil {
		tx.Rollback()
		return err
	}
	var requestOptions = make(map[string]string)
	if err = json.Unmarshal(p.Options.Bytes, &requestOptions); err != nil {
		return fmt.Errorf("格式化Options错误,请检查json字段是否为字符串类型: %v", err)
	}
	for optionKey := range requestOptions {
		var isInField = false
		for _, field := range fields {
			if optionKey == field.Name {
				isInField = true
				break
			}
		}
		if !isInField {
			return fmt.Errorf("Options值错误,请检查是否为对应产品类型字段")
		}
	}
	var product product.Product
	if err := logic.Gorm.Where("item_number = ?", p.ItemNumber).First(&product).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("未查询到该货号所属产品: %v", err)
	}
	if err := logic.Gorm.Model(&product).Updates(&p).Error; err != nil {
		return fmt.Errorf("更新产品失败,请检查: %v", err)
	}
	return nil
}

func (serviceProduct *ServiceProductGroup) DeleteProduct(p product.Product) error {
	tx := logic.Gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Commit()
		}
	}()
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	p.ItemNumber = p.ItemNumber + "_is_deleted" + currentTime
	if err := tx.Select("ItemNumber").Updates(&p).Error; err != nil {
		return fmt.Errorf("删除产品失败,请检查: %v", err)
	}
	if err := tx.Delete(&p).Error; err != nil {
		return fmt.Errorf("删除产品失败,请检查: %v", err)
	}
	return nil
}

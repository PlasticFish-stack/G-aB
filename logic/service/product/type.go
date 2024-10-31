package product

import (
	"errors"
	"fmt"
	"project/logic"
	"project/logic/model/product"
	"project/logic/service/tool"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ExpandTypeGroup(parentId uint, typeParent *product.Type) (err error) {
	var childrenType []*product.Type
	err = logic.Gorm.Where("parent_id = ?", parentId).Find(&childrenType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("查询Type子项错误: %v", err)
	}
	for _, children := range childrenType {
		err := ExpandTypeGroup(children.Id, children)
		if err != nil {
			return err
		}
		typeParent.Children = append(typeParent.Children, *children)
	}
	return nil
}

func (serviceProduct *ServiceProductGroup) SearchTypeTree(limits tool.RequestLimits) ([]*product.Type, *tool.ResponseLimits, error) {
	var types []*product.Type
	var total int64
	offset, err := limits.GetOffset()
	if err != nil {
		return nil, nil, err
	}
	if err := logic.Gorm.Order("type_sort").
		Where("parent_id=?", 0).
		Count(&total).
		Offset(offset).
		Limit(limits.PageSize).
		Find(&types).Error; err != nil {
		return nil, nil, fmt.Errorf("查询产品类别失败: %v", err)
	}
	for _, v := range types {
		err := ExpandTypeGroup(v.Id, v)
		if err != nil {
			return nil, nil, err
		}
	}
	ResponseLimit := tool.NewLimits(total, limits.PageSize, limits.PageNum)
	return types, ResponseLimit, nil
}

func (serviceProduct *ServiceProductGroup) SearchType(tid uint) (*product.Type, error) {
	var types *product.Type
	if err := logic.Gorm.Preload(clause.Associations).First(&types, tid).Error; err != nil {
		return nil, fmt.Errorf("查询产品类别失败: %v", err)
	}
	return types, nil
}

func (serviceProduct *ServiceProductGroup) AddType(t product.Type) error {
	if t.ParentId != 0 {
		_, err := serviceProduct.SearchType(t.ParentId)
		if err != nil {
			return fmt.Errorf("查询不到父菜单: %v", err)
		}
	}
	if err := logic.Gorm.Create(t).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("产品类别名称已存在: %v", err)
		}
		return fmt.Errorf("新建产品类别失败: %v", err)
	}
	return nil
}

func (serviceProduct *ServiceProductGroup) UpdateType(t product.Type) error {
	tx := logic.Gorm.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	var responseProductType *product.Type
	if err := tx.Where(t.Id).First(&responseProductType).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return fmt.Errorf("未查询到该产品类别: %v", err)
		}
		tx.Rollback()
		return fmt.Errorf("查询产品类别失败: %v", err)
	}
	if err := tx.Model(&responseProductType).Updates(&t).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新产品类别失败,请检查: %v", err)
	}

	if len(responseProductType.Fields) > 0 {
		GetFields, err := serviceProduct.SearchFields(*tx, responseProductType.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
		updates := map[string][]product.Field{
			"update": {},
			"delete": {},
		}
		for _, OriginField := range GetFields {
			var isDeleted = true
			for _, field := range responseProductType.Fields {
				if field.Id == OriginField.Id || field.Id == 0 {
					updates["update"] = append(updates["update"], field)
					isDeleted = false
					break
				}
			}
			if isDeleted {
				updates["delete"] = append(updates["delete"], OriginField)
			}
		}
		if len(updates["update"]) > 0 {
			for _, field := range updates["update"] {
				field.TypeId = responseProductType.Id // Ensure foreign key is set
				if field.Id == 0 {
					if err := tx.Create(&field).Error; err != nil {
						tx.Rollback()
						return err
					}
				} else {
					if err := tx.Updates(&field).Error; err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}
		if len(updates["delete"]) > 0 {
			if err := serviceProduct.DeleteFields(*tx, updates["delete"]); err != nil {
				tx.Rollback()
				return err
			}
		}
	} //Fields更新逻辑

	if len(responseProductType.Formulas) > 0 {
		GetFormulas, err := serviceProduct.SearchFormulas(*tx, responseProductType.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
		updates := map[string][]product.Formula{
			"update": {},
			"delete": {},
		}
		for _, OriginFormula := range GetFormulas {
			var isDeleted = true
			for _, formula := range responseProductType.Formulas {
				if formula.Id == OriginFormula.Id || formula.Id == 0 {
					updates["update"] = append(updates["update"], formula)
					isDeleted = false
					break
				}
			}
			if isDeleted {
				updates["delete"] = append(updates["delete"], OriginFormula)
			}
		}
		if len(updates["update"]) > 0 {
			for _, formula := range updates["update"] {
				formula.TypeId = responseProductType.Id // Ensure foreign key is set
				if formula.Id == 0 {
					if err := tx.Create(&formula).Error; err != nil {
						tx.Rollback()
						return err
					}
				} else {
					if err := tx.Updates(&formula).Error; err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}
		if len(updates["delete"]) > 0 {
			if err := serviceProduct.DeleteFormulas(*tx, updates["delete"]); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}

func (serviceProduct *ServiceProductGroup) DeleteType(db gorm.DB, t product.Type) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Commit()
		}
	}()
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	t.Name = t.Name + "_is_deleted" + currentTime
	if err := tx.Updates(&t).Error; err != nil {
		return fmt.Errorf("删除产品类别失败,请检查: %v", err)
	}
	if err := tx.Delete(&t).Error; err != nil {
		return fmt.Errorf("删除产品类别失败,请检查: %v", err)
	}
	return nil
}

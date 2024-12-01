package product

import (
	"errors"
	"fmt"
	"project/logic"
	"project/logic/model/product"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ExpandTypeGroup(parentId uint, typeParent *product.Type) (err error) {
	var childrenType []*product.Type
	err = logic.Gorm.Where("parent_id = ?", parentId).Preload(clause.Associations).Find(&childrenType).Error
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

func (serviceProduct *ServiceProductGroup) SearchTypeTree() ([]*product.Type, error) {
	var types []*product.Type
	if err := logic.Gorm.Order("type_sort").
		Where("parent_id=?", 0).
		Preload(clause.Associations).
		Find(&types).Error; err != nil {
		return nil, fmt.Errorf("查询产品类别失败: %v", err)
	}
	for _, v := range types {
		err := ExpandTypeGroup(v.Id, v)
		if err != nil {
			return nil, err
		}
	}
	return types, nil
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
	var repartType product.Type
	if err := logic.Gorm.Where("name = ?", t.Name).First(&repartType).Error; err == nil {
		if repartType.Name == t.Name {
			return fmt.Errorf("类别已经存在: %v", t.Name)
		}
	}
	fmt.Println(t)
	if err := logic.Gorm.Create(&t).Error; err != nil {
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
	var repartType product.Type
	if err := tx.Where("name = ?", t.Name).First(&repartType).Error; err == nil {
		if repartType.Name == t.Name && repartType.Id != t.Id {
			tx.Rollback()
			return fmt.Errorf("Name已经存在: %v", err)
		}
	}
	if err := tx.Model(&responseProductType).Updates(&t).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新产品类别失败,请检查: %v", err)
	}
	if len(t.Fields) > 0 {
		GetFields, err := serviceProduct.SearchFields(*tx, responseProductType.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
		updates := map[string][]product.Field{
			"update": {},
			"delete": {},
		}
		if len(GetFields) == 0 {
			updates["update"] = append(updates["update"], t.Fields...)
		} else {
			for _, OriginField := range GetFields {
				var isDeleted = true
				for _, field := range t.Fields {
					fmt.Println(field)
					if field.Id == OriginField.Id {
						fmt.Println("yes, ", field)
						updates["update"] = append(updates["update"], field)
						isDeleted = false
					} else if field.Id == 0 {
						if field.Name == OriginField.Name {
							tx.Rollback()
							return fmt.Errorf("更新同名Field时需携带id")
						}
						updates["update"] = append(updates["update"], field)
					}
				}
				if isDeleted {
					updates["delete"] = append(updates["delete"], OriginField)
				}
			}
		}

		if len(updates["update"]) > 0 {
			for _, field := range updates["update"] {
				field.TypeId = responseProductType.Id
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
		fmt.Println(updates)
	} //Fields更新逻辑

	if len(t.Formulas) > 0 {
		GetFormulas, err := serviceProduct.SearchFormulas(*tx, responseProductType.Id)
		if err != nil {
			tx.Rollback()
			return err
		}
		updates := map[string][]product.Formula{
			"update": {},
			"delete": {},
		}
		if len(GetFormulas) == 0 {
			updates["update"] = append(updates["update"], t.Formulas...)
		} else {
			for _, OriginFormula := range GetFormulas {
				var isDeleted = true
				for _, formula := range t.Formulas {
					if formula.Id == OriginFormula.Id {
						if formula.Id == 0 {
							if formula.Name == OriginFormula.Name {
								tx.Rollback()
								return fmt.Errorf("更新同名Formula时需携带id")
							}
							updates["update"] = append(updates["update"], formula)
						}
					} else if formula.Id == 0 {
						if formula.Name == OriginFormula.Name {
							tx.Rollback()
							return fmt.Errorf("更新同名Formula时需携带id")
						}
						updates["update"] = append(updates["update"], formula)
					} else {
						isDeleted = false
					}
				}
				if isDeleted {
					updates["delete"] = append(updates["delete"], OriginFormula)
				}
			}
		}

		if len(updates["update"]) > 0 {
			for _, formula := range updates["update"] {
				formula.TypeId = responseProductType.Id
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

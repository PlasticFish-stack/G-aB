package controll

import (
	"fmt"
	"project/logic"
	"project/logic/model"
)

func ProductTypeGetAll() ([]*model.ProdType, error) {
	var productType []*model.ProdType
	var err error
	if productType, err = model.SearchTreeProductType(logic.Gorm); err != nil {
		return nil, err
	}
	return productType, nil
}

func ProductTypeAdd(prodType model.ProdType) error {
	productType := &model.ProdType{
		Name:        name,
		Description: description,
		Sort:        sort,
		ParentId:    parentId,
		Tax:         tax,
		Field:       options,
		Formulas:    formulas,
	}
	err := productType.Add(logic.Gorm)
	if err != nil {
		return err
	}
	return nil
}

func ProductTypeUpdate(id uint, name string, description string, sort uint, parentId uint, tax float64, options []model.ProdTypeField, formulas []model.ProdTypeFormula) error {
	productType := &model.ProdType{
		Global:      model.Global{Id: id},
		Name:        name,
		Description: description,
		Sort:        sort,
		ParentId:    parentId,
		Tax:         tax,
		Field:       options,
		Formulas:    formulas,
	}
	if err := productType.Update(logic.Gorm); err != nil {
		return err
	}
	return nil
}

func ProductTypeDelete(id uint, name string) error {
	var prodType = model.ProdType{
		Global: model.Global{Id: id},
	}
	resultType, err := prodType.Search(logic.Gorm)
	if err != nil {
		return err
	}
	if name != resultType.Name {
		return fmt.Errorf("类别与id不匹配")
	}
	if err := resultType.Delete(logic.Gorm); err != nil {
		return err
	}
	return nil
}

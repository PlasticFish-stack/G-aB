package product

import (
	"project/logic/model"

	"github.com/jackc/pgtype"
)

type Product struct {
	model.Global
	ItemNumber     string       `gorm:"size:255;comment:货号;unique" json:"itemNumber"`
	BrandId        uint         `gorm:"comment:品牌" json:"brandId"`
	Sku            string       `gorm:"size:255;comment:sku" json:"sku"`
	Spu            string       `gorm:"size:255;comment:spu" json:"spu"`
	Quantity       uint64       `gorm:"default:1;comment:数量" json:"quantity"`
	Specifications string       `gorm:"type:text;comment:规格" json:"specifications"`
	Barcode        string       `gorm:"size:255;comment:条形码" json:"barcode"`
	Customscode    string       `gorm:"size:255;comment:海关编码" json:"customscode"`
	Description    string       `gorm:"size:255;comment:描述" json:"description"`
	Color          string       `gorm:"size:255;comment:颜色" json:"color"`
	DwPrice        float64      `gorm:"comment:得物价格" json:"dwPrice"`
	DwSales        float64      `gorm:"comment:得物销售额" json:"dwSales"`
	TypeId         uint         `gorm:"comment:产品类型id" json:"typeId"`
	Image          string       `gorm:"type:varchar(500);comment:产品主图路径" json:"image"`
	TypeName       string       `gorm:"-" json:"typeName"`
	Costs          []Cost       `gorm:"foreignKey:ProductID" json:"costs"`
	Options        pgtype.JSONB `gorm:"type:jsonb" json:"options"`
}

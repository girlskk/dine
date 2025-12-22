package domain

// ------------------------------------------------------------
// 枚举定义
// ------------------------------------------------------------

// ProductSupportType 支持类型
type ProductSupportType string

const (
	ProductSupportTypeDine     ProductSupportType = "dine"     // 堂食
	ProductSupportTypeTakeaway ProductSupportType = "takeaway" // 外带
	ProductSupportTypeDelivery ProductSupportType = "delivery" // 外卖
)

func (ProductSupportType) Values() []string {
	return []string{
		string(ProductSupportTypeDine),
		string(ProductSupportTypeTakeaway),
		string(ProductSupportTypeDelivery),
	}
}

// ProductAttrSelectionType 口味做法点单限制
type ProductAttrSelectionType string

const (
	ProductAttrSelectionTypeRequiredOne ProductAttrSelectionType = "required_one" // 必选一项
	ProductAttrSelectionTypeMultiple    ProductAttrSelectionType = "multiple"     // 可多选
)

func (ProductAttrSelectionType) Values() []string {
	return []string{
		string(ProductAttrSelectionTypeRequiredOne),
		string(ProductAttrSelectionTypeMultiple),
	}
}

// ProductSaleStatus 售卖状态
type ProductSaleStatus string

const (
	ProductSaleStatusOnSale  ProductSaleStatus = "on_sale"  // 在售
	ProductSaleStatusOffSale ProductSaleStatus = "off_sale" // 停售
)

func (ProductSaleStatus) Values() []string {
	return []string{
		string(ProductSaleStatusOnSale),
		string(ProductSaleStatusOffSale),
	}
}

// EffectiveDateType 生效日期类型
type EffectiveDateType string

const (
	EffectiveDateTypeDaily  EffectiveDateType = "daily"  // 按天
	EffectiveDateTypeCustom EffectiveDateType = "custom" // 自定义
)

func (EffectiveDateType) Values() []string {
	return []string{
		string(EffectiveDateTypeDaily),
		string(EffectiveDateTypeCustom),
	}
}

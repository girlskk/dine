package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Product 商品
type Product struct {
	ent.Schema
}

func (Product) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the Product.
func (Product) Fields() []ent.Field {
	return []ent.Field{
		// 基础信息
		field.Enum("type").GoType(domain.ProductType("")).Default(string(domain.ProductTypeNormal)).Comment("商品类型：normal（普通商品）、set_meal（套餐商品）"),
		field.String("name").MaxLen(255).NotEmpty().Comment("商品名称"),
		field.UUID("category_id", uuid.UUID{}).Comment("分类ID（支持一级分类和二级分类）"),
		field.UUID("menu_id", uuid.UUID{}).Optional().Comment("菜单ID"),
		field.String("mnemonic").MaxLen(255).Default("").Comment("助记词"),
		field.Int("shelf_life").Default(0).Comment("保质期（单位：天）"),
		field.JSON("support_types", []domain.ProductSupportType{}).Comment("支持类型（堂食、外带）"),
		// 属性关联
		field.UUID("unit_id", uuid.UUID{}).Comment("单位ID"),

		// 售卖信息
		field.Enum("sale_status").GoType(domain.ProductSaleStatus("")).Default(string(domain.ProductSaleStatusOnSale)).Comment("售卖状态：on_sale（在售）、off_sale（停售）"),
		field.JSON("sale_channels", []domain.SaleChannel{}).Comment("售卖渠道（可选，可多选）：POS、移动点餐、扫码点餐、自助点餐、三方外卖"),
		field.Enum("effective_date_type").GoType(domain.EffectiveDateType("")).Optional().Comment("生效日期类型：daily（按天）、custom（自定义）"),
		field.Time("effective_start_time").Optional().Nillable().Comment("生效开始时间（当类型为自定义时必填）"),
		field.Time("effective_end_time").Optional().Nillable().Comment("生效结束时间（当类型为自定义时必填）"),
		field.Int("min_sale_quantity").Optional().Comment("起售份数（可选，必须为正整数）"),
		field.Int("add_sale_quantity").Optional().Comment("加售份数（可选，必须为正整数）"),

		// 其他信息
		field.Bool("inherit_tax_rate").Default(true).Comment("是否继承原分类税率（必选）"),
		field.UUID("tax_rate_id", uuid.UUID{}).Optional().Comment("指定税率ID（当不继承时必填）"),
		field.Bool("inherit_stall").Default(true).Comment("是否继承原出品部门（必选）"),
		field.UUID("stall_id", uuid.UUID{}).Optional().Comment("指定出品部门ID（当不继承时必填）"),

		// 展示信息
		field.String("main_image").MaxLen(512).Default("").Comment("主图（可选，一张图片）"),
		field.JSON("detail_images", []string{}).Optional().Comment("详情图片（可选，多张）"),
		field.String("description").MaxLen(2000).Default("").Comment("菜品描述（可选）"),

		// 套餐信息
		field.Other("estimated_cost_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(19,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("预估成本价（可选，单位：令吉，仅套餐商品使用）"),
		field.Other("delivery_cost_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(19,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("外卖成本价（可选，单位：令吉，仅套餐商品使用）"),

		// 商户和门店
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("品牌商ID"),
		field.UUID("store_id", uuid.UUID{}).Default(schematype.NilUUID()).Immutable().Comment("门店ID"),
	}
}

func (Product) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
		index.Fields("category_id"),
		// 唯一索引
		index.Fields("merchant_id", "store_id", "name", "deleted_at").Unique(),
	}
}

// Edges of the Product.
func (Product) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("category", Category.Type).Ref("products").Field("category_id").Required().Unique().Comment("所属分类"),
		edge.From("unit", ProductUnit.Type).Ref("products").Field("unit_id").Required().Unique().Comment("所属单位"),
		edge.From("tax_rate", TaxFee.Type).Ref("products").Field("tax_rate_id").Unique().Comment("所属税率"),
		edge.From("stall", Stall.Type).Ref("products").Field("stall_id").Unique().Comment("所属出品部门"),

		// 商品标签 Many2Many
		edge.To("tags", ProductTag.Type).
			StorageKey(edge.Table("product_tag_relations"), edge.Columns("product_id", "tag_id")).
			Comment("关联标签"),

		// 商品规格
		edge.To("product_specs", ProductSpecRelation.Type),

		// 商品口味做法
		edge.To("product_attrs", ProductAttrRelation.Type),

		// 套餐组
		edge.To("set_meal_groups", SetMealGroup.Type),

		// 套餐组详情
		edge.To("set_meal_details", SetMealDetail.Type),

		// 菜单项
		edge.To("menu_items", MenuItem.Type).Comment("菜单项关联"),
	}
}

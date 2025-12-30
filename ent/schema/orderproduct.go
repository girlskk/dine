package schema

import (
	"encoding/json"

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

// OrderProduct 订单商品明细表
type OrderProduct struct {
	ent.Schema
}

func (OrderProduct) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the OrderProduct.
func (OrderProduct) Fields() []ent.Field {
	return []ent.Field{
		// ========== 订单关联 ==========
		field.UUID("order_id", uuid.UUID{}).Comment("所属订单ID"),
		field.String("order_item_id").NotEmpty().Comment("订单内明细ID"),
		field.Int("index").Default(0).Comment("下单序号（同订单内第几次下单）"),

		// ========== 商品基础信息（来自 domain.Product）==========
		field.UUID("product_id", uuid.UUID{}).Comment("商品ID"),
		field.String("product_name").NotEmpty().Comment("商品名称"),
		field.Enum("product_type").GoType(domain.ProductType("")).Default(string(domain.ProductTypeNormal)).Comment("商品类型：normal（普通商品）、set_meal（套餐商品）"),
		field.UUID("category_id", uuid.UUID{}).Optional().Comment("分类ID"),
		field.UUID("menu_id", uuid.UUID{}).Optional().Comment("菜单ID"),
		field.UUID("unit_id", uuid.UUID{}).Optional().Comment("单位ID"),
		field.JSON("support_types", []domain.ProductSupportType{}).Optional().Comment("支持类型（堂食、外带、外卖）"),
		field.Enum("sale_status").GoType(domain.ProductSaleStatus("")).Optional().Comment("售卖状态：on_sale（在售）、off_sale（停售）"),
		field.JSON("sale_channels", []domain.SaleChannel{}).Optional().Comment("售卖渠道"),
		field.String("main_image").MaxLen(512).Default("").Comment("商品主图"),
		field.String("description").MaxLen(2000).Default("").Comment("菜品描述"),

		// ========== 数量与金额 ==========
		field.Int("qty").Default(1).Comment("数量"),
		field.Other("subtotal", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("小计"),
		field.Other("discount_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("优惠金额"),
		field.Other("amount_before_tax", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("税前金额"),
		field.Other("tax_rate", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("税率（百分比，如 6.00 表示 6%）"),
		field.Other("tax", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("税额"),
		field.Other("amount_after_tax", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("税后金额"),
		field.Other("total", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("合计"),

		// ========== 促销信息 ==========
		field.Other("promotion_discount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("促销优惠金额"),

		// ========== 退菜信息 ==========
		field.Int("void_qty").Default(0).Comment("已退菜数量汇总"),
		field.Other("void_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("已退菜金额汇总"),
		field.String("refund_reason").Optional().Comment("退菜原因"),
		field.String("refunded_by").Optional().Comment("退菜操作人"),
		field.Time("refunded_at").Optional().Nillable().Comment("退菜时间"),

		// ========== 其他信息 ==========
		field.String("note").Optional().Comment("备注"),

		// ========== 套餐信息（仅套餐商品使用）==========
		field.Other("estimated_cost_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("预估成本价（仅套餐商品使用）"),
		field.Other("delivery_cost_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("外卖成本价（仅套餐商品使用）"),
		field.JSON("set_meal_groups", json.RawMessage{}).Optional().Comment("套餐组信息（包含套餐组名称、选择类型、详情列表等）"),

		// ========== 规格信息 ==========
		field.JSON("spec_relations", json.RawMessage{}).Optional().Comment("商品规格关联信息（规格名称、价格、库存等）"),

		// ========== 口味做法信息 ==========
		field.JSON("attr_relations", json.RawMessage{}).Optional().Comment("商品口味做法关联信息（口味名称、加价等）"),
	}
}

func (OrderProduct) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("order_id"),
		index.Fields("product_id"),
		index.Fields("order_id", "order_item_id").Unique(),
	}
}

// Edges of the OrderProduct.
func (OrderProduct) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", Order.Type).
			Ref("order_products").
			Field("order_id").
			Required().
			Unique().
			Comment("所属订单"),
	}
}

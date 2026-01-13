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

// RefundOrderProduct 退款订单商品明细表
type RefundOrderProduct struct {
	ent.Schema
}

func (RefundOrderProduct) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (RefundOrderProduct) Fields() []ent.Field {
	return []ent.Field{
		// ========== 退款单关联 ==========
		field.UUID("refund_order_id", uuid.UUID{}).Comment("退款单ID"),

		// ========== 原订单商品关联 ==========
		field.UUID("origin_order_product_id", uuid.UUID{}).Comment("原订单商品明细ID"),
		field.String("origin_order_item_id").Optional().Comment("原订单内明细ID"),

		// ========== 商品信息快照 ==========
		field.UUID("product_id", uuid.UUID{}).Comment("商品ID"),
		field.String("product_name").NotEmpty().Comment("商品名称"),
		field.Enum("product_type").GoType(domain.ProductType("")).
			Default(string(domain.ProductTypeNormal)).Comment("商品类型"),
		field.JSON("category", domain.Category{}).Optional().Comment("分类信息"),
		field.JSON("product_unit", domain.ProductUnit{}).Optional().Comment("商品单位信息"),
		field.String("main_image").MaxLen(512).Default("").Comment("商品主图"),
		field.String("description").MaxLen(2000).Default("").Comment("菜品描述"),

		// ========== 原订单数量与金额 ==========
		field.Int("origin_qty").Comment("原购买数量"),
		field.Other("origin_price", decimal.Decimal{}).
			SchemaType(map[string]string{dialect.MySQL: "DECIMAL(10,4)", dialect.SQLite: "NUMERIC"}).
			Optional().Nillable().Comment("原单价"),
		field.Other("origin_subtotal", decimal.Decimal{}).
			SchemaType(map[string]string{dialect.MySQL: "DECIMAL(10,4)", dialect.SQLite: "NUMERIC"}).
			Optional().Nillable().Comment("原小计"),
		field.Other("origin_discount", decimal.Decimal{}).
			SchemaType(map[string]string{dialect.MySQL: "DECIMAL(10,4)", dialect.SQLite: "NUMERIC"}).
			Optional().Nillable().Comment("原优惠金额"),
		field.Other("origin_tax", decimal.Decimal{}).
			SchemaType(map[string]string{dialect.MySQL: "DECIMAL(10,4)", dialect.SQLite: "NUMERIC"}).
			Optional().Nillable().Comment("原税额"),
		field.Other("origin_total", decimal.Decimal{}).
			SchemaType(map[string]string{dialect.MySQL: "DECIMAL(10,4)", dialect.SQLite: "NUMERIC"}).
			Optional().Nillable().Comment("原合计"),

		// ========== 退款数量与金额 ==========
		field.Int("refund_qty").Comment("退款数量"),
		field.Other("refund_subtotal", decimal.Decimal{}).
			SchemaType(map[string]string{dialect.MySQL: "DECIMAL(10,4)", dialect.SQLite: "NUMERIC"}).
			Optional().Nillable().Comment("退款小计"),
		field.Other("refund_discount", decimal.Decimal{}).
			SchemaType(map[string]string{dialect.MySQL: "DECIMAL(10,4)", dialect.SQLite: "NUMERIC"}).
			Optional().Nillable().Comment("退款优惠分摊"),
		field.Other("refund_tax", decimal.Decimal{}).
			SchemaType(map[string]string{dialect.MySQL: "DECIMAL(10,4)", dialect.SQLite: "NUMERIC"}).
			Optional().Nillable().Comment("退款税额"),
		field.Other("refund_total", decimal.Decimal{}).
			SchemaType(map[string]string{dialect.MySQL: "DECIMAL(10,4)", dialect.SQLite: "NUMERIC"}).
			Optional().Nillable().Comment("退款合计"),

		// ========== 规格/口味/套餐快照 ==========
		field.JSON("groups", domain.SetMealGroups{}).Optional().Comment("套餐组信息"),
		field.JSON("spec_relations", domain.ProductSpecRelations{}).Optional().Comment("规格信息"),
		field.JSON("attr_relations", domain.ProductAttrRelations{}).Optional().Comment("口味做法"),

		// ========== 退款原因 ==========
		field.Text("refund_reason").Optional().Comment("单品退款原因"),
	}
}

func (RefundOrderProduct) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("refund_order_id"),
		index.Fields("product_id"),
		index.Fields("origin_order_product_id"),
		index.Fields("refund_order_id", "origin_order_product_id", "deleted_at").Unique(),
	}
}

func (RefundOrderProduct) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("refund_order", RefundOrder.Type).
			Ref("refund_products").
			Field("refund_order_id").
			Required().
			Unique().
			Comment("所属退款单"),
	}
}

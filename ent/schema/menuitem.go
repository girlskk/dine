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

// MenuItem 菜单项
type MenuItem struct {
	ent.Schema
}

func (MenuItem) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the MenuItem.
func (MenuItem) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("menu_id", uuid.UUID{}).Immutable().Comment("菜单ID（外键）"),
		field.UUID("product_id", uuid.UUID{}).Immutable().Comment("菜品ID（外键，引用普通商品）"),
		field.Enum("sale_rule").
			GoType(domain.MenuItemSaleRule("")).
			Default(string(domain.MenuItemSaleRuleKeepBrandStatus)).
			Comment("下发售卖规则：keep_brand_status（保留品牌状态）、keep_store_status（保留门店状态）"),
		field.Other("base_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("基础价（可选，单位：分）"),
		field.Other("member_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(10,2)",
				dialect.SQLite: "NUMERIC",
			}).
			Optional().
			Nillable().
			Comment("会员价（可选，单位：分）"),
	}
}

func (MenuItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("menu_id"),
		index.Fields("product_id"),
	}
}

// Edges of the MenuItem.
func (MenuItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("menu", Menu.Type).
			Ref("items").
			Field("menu_id").
			Immutable().
			Required().
			Unique().
			Comment("所属菜单"),
		edge.From("product", Product.Type).
			Ref("menu_items").
			Field("product_id").
			Immutable().
			Required().
			Unique().
			Comment("关联商品"),
	}
}

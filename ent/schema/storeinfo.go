package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// StoreInfo holds the schema definition for the StoreInfo entity.
type StoreInfo struct {
	ent.Schema
}

// Fields of the StoreInfo.
func (StoreInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("city").Optional().Comment("省市地区"),
		field.String("address").Optional().Comment("详细地址"),
		field.String("contact_name").Optional().Comment("门店联系人"),
		field.String("contact_phone").Optional().Comment("联系人电话"),
		field.JSON("images", domain.StoreInfoImages{}).
			Default(domain.StoreInfoImages{}).
			Optional().
			Comment("门店图片"),
		field.Int("store_id"),
	}
}

// Edges of the StoreInfo.
func (StoreInfo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("store", Store.Type).Ref("store_info").Unique().Field("store_id").Required(),
	}
}

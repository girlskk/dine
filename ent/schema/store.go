package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Store holds the schema definition for the Store entity.
type Store struct {
	ent.Schema
}

// Fields of the Store.
func (Store) Fields() []ent.Field {
	return []ent.Field{
		field.Int("merchant_id").
			Immutable().
			Comment("商户 ID"),
		field.String("admin_phone_number").
			NotEmpty().
			Default("").
			MaxLen(20).
			Comment("管理员手机号"),
		field.String("store_name").
			NotEmpty().
			Default("").
			MaxLen(30).
			Comment("门店名称,长度不超过30个字"),
		field.String("store_short_name").
			NotEmpty().
			Default("").
			MaxLen(30).
			Comment("门店简称"),
		field.String("store_code").
			NotEmpty().
			Default("").
			Comment("门店编码(保留字段)"),
		field.Enum("status").
			GoType(domain.StoreStatus("")).
			Comment("状态: 营业 停业"),
		field.Enum("business_model").
			GoType(domain.BusinessModel("")).
			Comment("经营模式：直营 加盟"),
		field.Int("business_type_id").
			Default(0).
			Comment("业态类型"),
		field.String("contact_name").
			NotEmpty().
			Default("").
			MaxLen(20).
			Comment("联系人"),
		field.String("contact_phone").
			NotEmpty().
			Default("").
			MaxLen(20).
			Comment("联系电话"),
		field.String("unified_social_credit_code").
			NotEmpty().
			Default("").
			MaxLen(50).
			Comment("统一社会信用代码"),
		field.String("store_logo").
			NotEmpty().
			Default("").
			MaxLen(500).
			Comment("logo 图片地址"),
		field.String("business_license_url").
			NotEmpty().
			Default("").
			MaxLen(500).
			Comment("营业执照图片"),
		field.String("storefront_url").
			NotEmpty().
			Default("").
			MaxLen(500).
			Comment("门店门头照"),
		field.String("cashier_desk_url").
			NotEmpty().
			Default("").
			MaxLen(500).
			Comment("门店收银台照片"),
		field.String("dining_environment_url").
			NotEmpty().
			Default("").
			MaxLen(500).
			Comment("就餐环境图"),
		field.String("food_operation_license_url").
			NotEmpty().
			Default("").
			MaxLen(500).
			Comment("食品经营许可证照片"),
		// 地区信息
		field.Int("country_id").
			Default(0).
			Comment("国家/地区 id"),
		field.Int("province_id").
			Default(0).
			Comment("省份 id"),
		field.Int("city_id").
			Default(0).
			Comment("城市 id"),
		field.Int("district_id").
			Default(0).
			Comment("区县 id"),
		field.String("country_name").
			NotEmpty().
			Default("").
			MaxLen(50).
			Comment("国家/地区"),
		field.String("province_name").
			NotEmpty().
			Default("").
			MaxLen(50).
			Comment("省份"),
		field.String("city_name").
			NotEmpty().
			Default("").
			MaxLen(50).
			Comment("城市"),
		field.String("district_name").
			NotEmpty().
			Default("").
			MaxLen(50).
			Comment("区县"),
		field.String("address").
			NotEmpty().
			Default("").
			MaxLen(255).
			Comment("详细地址"),
		field.String("lng").
			NotEmpty().
			Default("").
			MaxLen(50).
			Comment("经度"),
		field.String("lat").
			NotEmpty().
			Default("").
			MaxLen(50).
			Comment("纬度"),
		field.UUID("admin_user_id", uuid.UUID{}).
			Immutable().
			Comment("登陆账号 ID"),
	}
}

// Edges of the Store.
func (Store) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("merchant", Merchant.Type).
			Ref("stores").
			Field("merchant_id").
			Unique().
			Immutable().
			Required(),
		edge.From("admin_user", AdminUser.Type).
			Ref("store").
			Field("admin_user_id").
			Unique().
			Immutable().
			Required(),
	}
}

func (Store) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

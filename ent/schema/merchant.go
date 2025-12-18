package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// Merchant holds the schema definition for the Merchant entity.
type Merchant struct {
	ent.Schema
}

// Fields of the Merchant.
func (Merchant) Fields() []ent.Field {
	return []ent.Field{
		// 商户基础信息
		field.String("merchant_code").
			NotEmpty().
			Default("").
			Comment("商户编号(保留字段)"),
		field.String("merchant_name").
			NotEmpty().
			Default("").
			MaxLen(50).
			Comment("商户名称,最长不得超过50个字"),
		field.String("merchant_short_name").
			NotEmpty().
			Default("").
			MaxLen(50).
			Comment("商户简称"),
		field.Enum("merchant_type").
			GoType(domain.MerchantType("")).
			Comment("商户类型: 品牌商户,门店商户"),
		field.String("brand_name").
			NotEmpty().
			Default("").
			Comment("品牌名称"),
		field.String("admin_phone_number").
			NotEmpty().
			Default("").
			MaxLen(20).
			Comment("管理员手机号"),
		field.Time("expire_utc").
			Optional().
			Nillable().
			Comment("UTC 时区的过期时间"),
		field.Int("business_type_id").
			Default(0).
			Comment("业务类型"),
		field.String("merchant_logo").
			Default("").
			MaxLen(500).
			Comment("logo 图片地址"),
		field.String("description").
			NotEmpty().
			Default("").
			MaxLen(255).
			Comment("商户描述(保留字段)"),
		field.Enum("status").
			GoType(domain.MerchantStatus("")).
			Comment("状态: 正常,停用,过期"),
		field.String("login_account").
			NotEmpty().
			Default("").
			MaxLen(50).
			Comment("登录账号"),
		field.String("login_password").
			NotEmpty().
			Default("").
			MaxLen(255).
			Comment("登录密码(加密存储)"),

		// 地区信息
		field.Int("country_id").
			Default(0).
			Comment("国家/地区id"),
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
			Comment("详细地址"),
		field.String("lng").
			NotEmpty().
			Default("").
			Comment("经度"),
		field.String("lat").
			NotEmpty().
			Default("").
			Comment("纬度"),
	}
}

// Edges of the Merchant.
func (Merchant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("stores", Store.Type),
		edge.To("merchant_renewals", MerchantRenewal.Type),
		edge.From("merchant_business_type", MerchantBusinessType.Type).
			Ref("merchants").
			Field("business_type_id").
			Unique().
			Required(),
	}
}

func (Merchant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

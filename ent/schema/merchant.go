package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
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
		field.UUID("business_type_id", uuid.UUID{}).
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

		// 地区信息
		field.UUID("country_id", uuid.UUID{}).
			Optional().
			Comment("国家/地区id"),
		field.UUID("province_id", uuid.UUID{}).
			Optional().
			Comment("省份 id"),
		field.UUID("city_id", uuid.UUID{}).
			Optional().
			Comment("城市 id"),
		field.UUID("district_id", uuid.UUID{}).
			Optional().
			Comment("区县 id"),
		field.String("address").
			NotEmpty().
			Default("").
			MaxLen(255).
			Comment("详细地址"),
		field.String("lng").
			NotEmpty().
			Default("").
			Comment("经度"),
		field.String("lat").
			NotEmpty().
			Default("").
			Comment("纬度"),
		field.UUID("admin_user_id", uuid.UUID{}).
			Immutable().
			Comment("登陆账号 ID"),
	}
}

// Edges of the Merchant.
func (Merchant) Edges() []ent.Edge {
	return []ent.Edge{
		// 业态类型
		edge.From("merchant_business_type", MerchantBusinessType.Type).
			Ref("merchants").
			Field("business_type_id").
			Unique().
			Required(),
		// 管理员账号
		edge.From("admin_user", AdminUser.Type).
			Ref("merchant").
			Field("admin_user_id").
			Unique().
			Immutable().
			Required(),
		// 地区关联（绑定已有外键字段）
		edge.From("country", Country.Type).
			Ref("merchants").
			Field("country_id").
			Unique(),
		edge.From("province", Province.Type).
			Ref("merchants").
			Field("province_id").
			Unique(),
		edge.From("city", City.Type).
			Ref("merchants").
			Field("city_id").
			Unique(),
		edge.From("district", District.Type).
			Ref("merchants").
			Field("district_id").
			Unique(),
		edge.To("stores", Store.Type),
		edge.To("merchant_renewals", MerchantRenewal.Type),
		edge.To("remark_categories", RemarkCategory.Type),
		edge.To("remarks", Remark.Type),
		edge.To("stalls", Stall.Type),
	}
}

func (Merchant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

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
			Optional().
			Default("").
			MaxLen(100).
			Comment("商户编号(保留字段)"),
		field.String("merchant_name").
			NotEmpty().
			Default("").
			MaxLen(100).
			Comment("商户名称,最长不得超过50个字"),
		field.String("merchant_short_name").
			Optional().
			Default("").
			MaxLen(100).
			Comment("商户简称"),
		field.Enum("merchant_type").
			GoType(domain.MerchantType("")).
			Comment("商户类型: 品牌商户,门店商户"),
		field.String("brand_name").
			Optional().
			Default("").
			MaxLen(100).
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
		field.String("business_type_code").
			Comment("业态类型"),
		field.String("merchant_logo").
			Default("").
			MaxLen(500).
			Comment("logo 图片地址"),
		field.String("description").
			Optional().
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
			Optional().
			Default("").
			MaxLen(255).
			Comment("详细地址"),
		field.String("lng").
			Optional().
			Default("").
			Comment("经度"),
		field.String("lat").
			Optional().
			Default("").
			Comment("纬度"),
		field.String("super_account").
			Immutable().
			Comment("登陆账号"),
	}
}

// Edges of the Merchant.
func (Merchant) Edges() []ent.Edge {
	return []ent.Edge{
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
		edge.To("backend_users", BackendUser.Type),
		edge.To("stores", Store.Type),
		edge.To("merchant_renewals", MerchantRenewal.Type),
		edge.To("remark_categories", RemarkCategory.Type),
		edge.To("remarks", Remark.Type),
		edge.To("stalls", Stall.Type),
		edge.To("additional_fees", AdditionalFee.Type),
		edge.To("tax_fees", TaxFee.Type),
		edge.To("devices", Device.Type),
		edge.To("departments", Department.Type),
		edge.To("roles", Role.Type),
		edge.To("store_users", StoreUser.Type),
	}
}

func (Merchant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

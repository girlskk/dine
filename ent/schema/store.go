package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.UUID("merchant_id", uuid.UUID{}).
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
			MaxLen(100).
			Comment("门店名称,长度不超过30个字"),
		field.String("store_short_name").
			Optional().
			Default("").
			MaxLen(100).
			Comment("门店简称"),
		field.String("store_code").
			Optional().
			Default("").
			Comment("门店编码(保留字段)"),
		field.Enum("status").
			GoType(domain.StoreStatus("")).
			Comment("状态: 营业 停业"),
		field.Enum("business_model").
			GoType(domain.BusinessModel("")).
			Comment("经营模式：直营 加盟"),
		field.String("business_type_code").
			Comment("业态类型"),
		field.String("location_number").
			NotEmpty().
			MaxLen(255).
			Comment("门店位置编号"),
		field.String("contact_name").
			Optional().
			Default("").
			MaxLen(100).
			Comment("联系人"),
		field.String("contact_phone").
			Optional().
			Default("").
			MaxLen(20).
			Comment("联系电话"),
		field.String("unified_social_credit_code").
			Optional().
			Default("").
			MaxLen(50).
			Comment("统一社会信用代码"),
		field.String("store_logo").
			Optional().
			Default("").
			MaxLen(500).
			Comment("logo 图片地址"),
		field.String("business_license_url").
			Optional().
			Default("").
			MaxLen(500).
			Comment("营业执照图片"),
		field.String("storefront_url").
			Optional().
			Default("").
			MaxLen(500).
			Comment("门店门头照"),
		field.String("cashier_desk_url").
			Optional().
			Default("").
			MaxLen(500).
			Comment("门店收银台照片"),
		field.String("dining_environment_url").
			Optional().
			Default("").
			MaxLen(500).
			Comment("就餐环境图"),
		field.String("food_operation_license_url").
			Optional().
			Default("").
			MaxLen(500).
			Comment("食品经营许可证照片"),
		field.JSON("business_hours", []domain.BusinessHours{}).
			Comment("营业时间段，JSON格式存储"),
		field.JSON("dining_periods", []domain.DiningPeriod{}).
			Comment("就餐时段，JSON格式存储"),
		field.JSON("shift_times", []domain.ShiftTime{}).
			Comment("班次时间，JSON格式存储"),
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
			MaxLen(50).
			Comment("经度"),
		field.String("lat").
			Optional().
			Default("").
			MaxLen(50).
			Comment("纬度"),
		field.String("super_account").
			Immutable().
			Comment("登陆账号"),
	}
}

// Edges of the Store.
func (Store) Edges() []ent.Edge {
	return []ent.Edge{
		// 所属商户
		edge.From("merchant", Merchant.Type).
			Ref("stores").
			Field("merchant_id").
			Unique().
			Immutable().
			Required(),
		// 地区关联（绑定已有外键字段）
		edge.From("country", Country.Type).
			Ref("stores").
			Field("country_id").
			Unique(),
		edge.From("province", Province.Type).
			Ref("stores").
			Field("province_id").
			Unique(),
		edge.From("city", City.Type).
			Ref("stores").
			Field("city_id").
			Unique(),
		edge.From("district", District.Type).
			Ref("stores").
			Field("district_id").
			Unique(),
		edge.To("store_users", StoreUser.Type),
		edge.To("remarks", Remark.Type),
		edge.To("stalls", Stall.Type),
		edge.To("additional_fees", AdditionalFee.Type),
		edge.To("tax_fees", TaxFee.Type),
		edge.To("devices", Device.Type),
		edge.From("menus", Menu.Type).Ref("stores").Comment("关联的菜单"),
		edge.To("departments", Department.Type),
		edge.To("roles", Role.Type),
	}
}

func (Store) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (Store) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
	}
}

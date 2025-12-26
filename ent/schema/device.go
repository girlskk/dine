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

// Device holds the schema definition for the Device entity.
type Device struct {
	ent.Schema
}

func (Device) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the Device.
func (Device) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(50).
			Comment("设备名称，长度不超过50字"),
		field.UUID("merchant_id", uuid.UUID{}).
			Immutable().
			Comment("品牌商 ID"),
		field.UUID("store_id", uuid.UUID{}).
			Immutable().
			Comment("门店 ID"),
		field.Enum("device_type").
			GoType(domain.DeviceType("")).
			Comment("设备类型：收银机/打印机"),
		field.String("device_code").
			MaxLen(100).
			Default("").
			Comment("设备编号/序列号"),
		field.String("device_brand").
			Optional().
			Comment("设备品牌"),
		field.String("device_model").
			Optional().
			Comment("设备型号"),
		field.Enum("location").
			GoType(domain.DeviceLocation("")).
			Comment("设备位置：前厅/后厨"),
		field.Bool("enabled").
			Default(true).
			Comment("启用/停用状态"),
		field.Enum("status").
			GoType(domain.DeviceStatus("")).
			Default(string(domain.DeviceStatusOffline)).
			Comment("设备状态：在线/离线"),
		field.String("ip").
			MaxLen(50).
			Optional().
			Default("").
			Comment("设备 IP 地址"),
		field.Int("sort_order").
			Default(1000).
			Comment("排序，值越小越靠前"),
		field.Enum("paper_size").
			GoType(domain.PaperSize("")).
			Optional().
			Comment("打印纸张尺寸"),
		field.UUID("stall_id", uuid.UUID{}).
			Optional().
			Comment("出品部门ID，可为空"),
		field.JSON("order_channels", []domain.OrderChannel{}).
			Optional().
			Comment("订单来源，取自 order_channel，多选"),
		field.JSON("dining_ways", []domain.DiningWay{}).
			Optional().
			Comment("订单类型/就餐方式，取自 dining_way，多选"),
		field.Enum("device_stall_print_type").
			GoType(domain.DeviceStallPrintType("")).
			Optional().
			Comment("打印出品部门总分单"),
		field.Enum("device_stall_receipt_type").
			GoType(domain.DeviceStallReceiptType("")).
			Optional().
			Comment("打印出品部门全部票据"),
		field.Bool("open_cash_drawer").
			Optional().
			Comment("是否开钱箱"),
	}
}

// Edges of the Device.
func (Device) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("merchant", Merchant.Type).
			Ref("devices").
			Field("merchant_id").
			Unique().
			Immutable().
			Required(),
		edge.From("store", Store.Type).
			Ref("devices").
			Field("store_id").
			Unique().
			Immutable().
			Required(),
		edge.From("stall", Stall.Type).
			Ref("devices").
			Field("stall_id").
			Unique(),
	}
}

// Indexes of the Device.
func (Device) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
		index.Fields("stall_id"),
		index.Fields("name", "merchant_id", "store_id", "deleted_at").
			Unique().
			StorageKey("idx_device_name_scope_deleted"),
	}
}

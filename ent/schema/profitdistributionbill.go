package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// ProfitDistributionBill 分账账单
type ProfitDistributionBill struct {
	ent.Schema
}

func (ProfitDistributionBill) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

func (ProfitDistributionBill) Fields() []ent.Field {
	return []ent.Field{
		field.String("no").MaxLen(64).NotEmpty().Unique().Comment("分账账单编号"),
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("品牌商ID"),
		field.UUID("store_id", uuid.UUID{}).Comment("门店ID"),
		field.UUID("revenue_id", uuid.UUID{}).Comment("门店营业额ID"),
		field.Other("receivable_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(19,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Comment("应收金额（令吉）"),
		field.Other("payment_amount", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(19,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Comment("打款金额（令吉）"),
		field.Enum("status").
			GoType(domain.ProfitDistributionBillStatus("")).
			Default(string(domain.ProfitDistributionBillStatusUnpaid)).
			Comment("分账状态：unpaid（未打款）、paid（已打款）"),
		field.Time("bill_date").
			SchemaType(map[string]string{
				dialect.MySQL:  "DATE",
				dialect.SQLite: "DATE",
			}).
			Comment("账单日期"),
		field.Time("start_date").
			SchemaType(map[string]string{
				dialect.MySQL:  "DATE",
				dialect.SQLite: "DATE",
			}).
			Comment("账单周期：开始日期"),
		field.Time("end_date").
			SchemaType(map[string]string{
				dialect.MySQL:  "DATE",
				dialect.SQLite: "DATE",
			}).
			Comment("账单周期：结束日期"),
		field.JSON("rule_snapshot", &domain.ProfitDistributionRuleSnapshot{}).
			Comment("分账方案快照"),
	}
}

func (ProfitDistributionBill) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		index.Fields("store_id"),
		// 唯一索引：同一门店同一账单日期只能有一条账单
		index.Fields("store_id", "bill_date", "deleted_at").Unique(),
	}
}

func (ProfitDistributionBill) Edges() []ent.Edge {
	return []ent.Edge{
		// @TODO
		// 关联营业额记录（如果 StoreDailyRevenue 实体已创建）
		// edge.To("revenue", StoreDailyRevenue.Type).
		// 	Unique().
		// 	Field("revenue_id").
		// 	Comment("关联的营业额记录"),
	}
}

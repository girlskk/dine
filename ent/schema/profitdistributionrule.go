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

// ProfitDistributionRule 分账方案
type ProfitDistributionRule struct {
	ent.Schema
}

func (ProfitDistributionRule) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the ProfitDistributionRule.
func (ProfitDistributionRule) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("merchant_id", uuid.UUID{}).Immutable().Comment("品牌商ID"),
		field.String("name").MaxLen(255).NotEmpty().Comment("分账方案名称"),
		field.Other("split_ratio", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL:  "DECIMAL(19,4)",
				dialect.SQLite: "NUMERIC",
			}).
			Comment("分账比例（0-1，单位：小数）"),
		field.Enum("billing_cycle").
			GoType(domain.ProfitDistributionRuleBillingCycle("")).
			Default(string(domain.ProfitDistributionRuleBillingCycleDaily)).
			Comment("账单生成周期：daily（按日）、monthly（按月）"),
		field.Time("effective_date").Comment("方案生效日期"),
		field.Time("expiry_date").Comment("方案失效日期"),
		field.Int("bill_generation_day").Default(1).Comment("账单生成日：1-28号"),
		field.Enum("status").
			GoType(domain.ProfitDistributionRuleStatus("")).
			Default(string(domain.ProfitDistributionRuleStatusDisabled)).
			Comment("状态：enabled（启用）、disabled（禁用）"),
		field.Int("store_count").Default(0).Comment("关联门店数量"),
	}
}

func (ProfitDistributionRule) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id"),
		// 唯一索引：品牌商下分账方案名称唯一
		index.Fields("merchant_id", "name", "deleted_at").Unique(),
	}
}

// Edges of the ProfitDistributionRule.
func (ProfitDistributionRule) Edges() []ent.Edge {
	return []ent.Edge{
		// 关联门店 Many2Many
		edge.To("stores", Store.Type).
			StorageKey(edge.Table("profit_distribution_rule_store_relations"),
				edge.Columns("profit_distribution_rule_id", "store_id")).
			Comment("关联门店"),
	}
}

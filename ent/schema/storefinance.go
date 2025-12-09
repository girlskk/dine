package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// StoreFinance holds the schema definition for the StoreFinance entity.
type StoreFinance struct {
	ent.Schema
}

// Fields of the StoreFinance.
func (StoreFinance) Fields() []ent.Field {
	return []ent.Field{
		field.String("bank_account").Optional().Comment("银行账号"),
		field.String("bank_card_name").Optional().Comment("银行账户名称"),
		field.String("bank_name").Optional().Comment("银行名称"),
		field.String("branch_name").Optional().Comment("开户支行"),
		field.String("public_account").Optional().Comment("对公账号"),
		field.String("company_name").Optional().Comment("公司名称"),
		field.String("public_bank_name").Optional().Comment("对公银行名称"),
		field.String("public_branch_name").Optional().Comment("对公开户支行"),
		field.String("credit_code").Optional().Comment("统一社会信用代码"),
		field.Int("store_id"),
	}
}

// Edges of the StoreFinance.
func (StoreFinance) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("store", Store.Type).Ref("store_finance").Unique().Field("store_id").Required(),
	}
}

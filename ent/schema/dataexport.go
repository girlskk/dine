package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// DataExport holds the schema definition for the DataExport entity.
type DataExport struct {
	ent.Schema
}

// Fields of the DataExport.
func (DataExport) Fields() []ent.Field {
	return []ent.Field{
		field.Int("store_id").
			NonNegative().
			Default(0).
			Immutable().
			Comment("门店ID"),
		field.Enum("type").
			GoType(domain.DataExportType("")).
			Immutable().
			Comment("导出类型"),
		field.Enum("status").
			GoType(domain.DataExportStatus("")).
			Comment("导出状态"),
		field.JSON("params", json.RawMessage{}).
			Comment("导出参数"),
		field.String("failed_reason").
			Optional().
			Comment("导出失败原因"),
		field.Enum("operator_type").
			GoType(domain.OperatorType("")).
			Immutable().
			Comment("操作人类型"),
		field.Int("operator_id").
			NonNegative().
			Immutable().
			Default(0).
			Comment("操作人ID"),
		field.String("operator_name").
			Immutable().
			Comment("操作人姓名"),
		field.String("file_name").
			Immutable().
			Comment("文件名"),
		field.String("url").
			Optional().
			Comment("下载地址"),
	}
}

// Edges of the DataExport.
func (DataExport) Edges() []ent.Edge {
	return nil
}

// Indexes of the DataExport.
func (DataExport) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("store_id", "deleted_at"),
	}
}

// Mixin of the DataExport.
func (DataExport) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

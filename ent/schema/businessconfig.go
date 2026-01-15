package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

// BusinessConfig 结算方式
type BusinessConfig struct {
	ent.Schema
}

func (BusinessConfig) Mixin() []ent.Mixin {
	return []ent.Mixin{
		schematype.UUIDMixin{},
		schematype.TimeMixin{},
		schematype.SoftDeleteMixin{},
	}
}

// Fields of the BusinessConfig.
func (BusinessConfig) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("source_config_id", uuid.UUID{}).Optional().Comment("配置来源ID"),
		field.UUID("merchant_id", uuid.UUID{}).Optional().Comment("品牌商ID"),
		field.String("store_id").
			SchemaType(map[string]string{
				dialect.MySQL: "char(36)",
			}).
			Default("").
			Comment("门店ID"),
		field.Enum("group").
			GoType(domain.BusinessConfigGroup("")).
			Optional().
			Comment("配置分组"),
		field.String("name").SchemaType(map[string]string{
			dialect.MySQL: "varchar(100)",
		}).Optional().Default("").Comment("参数名称"),
		field.Enum("config_type").
			GoType(domain.BusinessConfigConfigType("")).
			Optional().
			Comment("键值类型:string,int,uint,bool,datetime,date"),
		field.String("key").SchemaType(map[string]string{
			dialect.MySQL: "varchar(100)",
		}).Optional().Default("").Comment("参数键名"),
		field.String("value").SchemaType(map[string]string{
			dialect.MySQL: "varchar(500)",
		}).Optional().Default("").Comment("参数键值"),
		field.Int32("sort").SchemaType(map[string]string{
			dialect.MySQL: "int",
		}).Default(0).Comment("排序 越小越靠前"),
		field.String("tip").SchemaType(map[string]string{
			dialect.MySQL: "varchar(500)",
		}).Optional().Default("").Comment("变量描述"),
		field.Bool("is_default").Default(false).Comment("是否为系统默认: true-是, false-否）"),
		field.Bool("modify_status").Default(true).Comment("下发后是否可以进行修改: true-可以, false-不可以）"),
		field.Bool("status").Default(true).Comment("启用/停用状态: true-启用, false-停用（必选）"),
	}
}

func (BusinessConfig) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("merchant_id", "store_id", "group", "key", "deleted_at").Unique(),
	}
}

// Edges of the BusinessConfig.
func (BusinessConfig) Edges() []ent.Edge {
	return []ent.Edge{}
}

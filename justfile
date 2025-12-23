set shell := ["sh", "-cu"]

MIGRATION_DIR := "file://ent/migrate/migrations"
DEV_DB := "docker://mysql/8/ent"
ENT_CMD_NEW := "go run -mod=mod entgo.io/ent/cmd/ent new"

# 生成代码
default:
    go generate ./... && \
    go mod tidy

ent:
    go generate ./ent

proto name:
    protoc --proto_path=pkg/ugrpc/third_party --proto_path=api/{{name}}/pb \
        --go_out=paths=source_relative:api/{{name}}/pb \
        --go-grpc_out=paths=source_relative:api/{{name}}/pb \
        --validate_out=paths=source_relative,lang=go:api/{{name}}/pb \
        api/{{name}}/pb/*.proto

# 生成迁移文件
migrate name:
    atlas migrate diff {{name}} \
        --dir {{MIGRATION_DIR}} \
        --to "ent://ent/schema" \
        --dev-url {{DEV_DB}} \
        --format "{{'{{'}} sql . \"  \" {{'}}'}}"

# 生成手动迁移文件
migrate_manual name:
    atlas migrate new {{name}} --dir={{MIGRATION_DIR}}

# 重新生成迁移Hash
migrate_hash:
    atlas migrate hash --dir={{MIGRATION_DIR}}

# 生成实体
ent_new +names:
    {{ENT_CMD_NEW}} {{names}}

# 本地启动服务
run name:
    air -c ./.air/{{name}}.toml

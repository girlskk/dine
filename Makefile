DB_DSN=mysql://root:pass@:33061/dine
#DB_DSN=mysql://dine:6fFh2M44cEKv@rm-uf66l9d6gm5553v665o.mysql.rds.aliyuncs.com:3306/dine
ENT_DIR=.

# 基于 ent schema 执行数据库迁移
.PHONY: schema_apply
schema_apply:
	atlas schema apply \
 		-u "$(DB_DSN)" \
 		--to "ent://ent/schema" \
 		--dev-url "$(DB_DSN)_test"

# 生成迁移文件
.PHONY: schema_generate
schema_generate:
	atlas migrate diff init \
	  --dir "file://ent/migrate/migrations" \
	  --to "ent://ent/schema" \
	  --dev-url "mysql://root:pass@:33061/dine_test"

# 基于 ent schema 生成代码
.PHONY: ent_generate
ent_generate:
	cd $(ENT_DIR) && go generate ./ent


# 创建一个 ent schema
# make schema_create name=SetMealDetail
.PHONY: schema_create
schema_create:
	@if [ "$(name)" = "" ]; then \
		echo "请指定 schema 名称，使用方式：make schema_create name=[schema_name]"; \
		exit 1; \
	fi
	cd $(ENT_DIR) && go run -mod=mod entgo.io/ent/cmd/ent new $(name)


# 删除 ent 生成的文件
.PHONY: ent_clean
ent_clean:
	@echo "开始备份 schema 文件夹..."
	@cd $(ENT_DIR)/ent && cp -R schema schema_bak
	@echo "开始清空文件..."
	@cd $(ENT_DIR)/ent && find schema -mindepth 1 -exec rm -rf {} \;
	@cd $(ENT_DIR) && find ent -mindepth 1 \( -name schema -o -name schema_bak -o -name generate.go \) -prune -o -exec rm -rf {} \;

# 生成代码
.PHONY: gen
gen:
	go generate ./...


# 启动服务 - 使用 air 进行热重载开发
.PHONY: run-admin
run-admin:
	air -c ./.air/admin.toml

.PHONY: run-backend
run-backend:
	air -c ./.air/backend.toml

.PHONY: run-frontend
run-frontend:
	air -c ./.air/frontend.toml

.PHONY: run-customer
run-customer:
	air -c ./.air/customer.toml

.PHONY: run-intl
run-intl:
	air -c ./.air/intl.toml

.PHONY: run-scheduler
run-scheduler:
	air -c ./.air/scheduler.toml
# POS Dine API 开发规范

## 1. 项目概述

POS Dine API 是一个基于 Go 语言开发的餐饮管理系统后端 API，采用领域驱动设计 (DDD) 架构模式，使用 entgo 作为 ORM 框架，Gin 作为 Web 框架，支持多服务部署。

## 2. 架构设计

### 2.1 整体架构

项目采用**整洁架构 (Clean Architecture)** 设计模式，主要分为以下核心层：

- **API 层**：处理 HTTP 请求和响应，参数验证，路由管理
- **UseCase 层**：实现业务逻辑，协调领域对象和基础设施
- **Domain 层**：定义核心业务模型和领域服务
- **Repository 层**：数据访问层，实现领域模型的持久化
- **Infrastructure 层**：提供基础设施支持，如数据库、缓存、消息队列等

### 2.2 整洁架构的引用关系

#### 核心原则：内层不依赖外层

整洁架构的核心思想是 **依赖倒置原则**，即外层依赖内层，内层不依赖外层。内层包含核心业务逻辑，外层包含与外部系统的交互。

```
┌─────────────────────────────────────────────────────────────────┐
│                     外层依赖内层                                 │
│                                                                 │
│  ┌────────────────┐     ┌────────────────┐     ┌──────────────┐  │
│  │    API 层      │────▶│   UseCase 层   │────▶│  Domain 层   │  │
│  └────────────────┘     └────────────────┘     └──────────────┘  │
│        ▲                        ▲                       │         │
│        │                        │                       │         │
│        │                        │                       │         │
│  ┌────────────────┐     ┌────────────────┐     ┌──────────────┐  │
│  │  基础设施层    │◀────│  Repository 层  │◀────│  Domain 层   │  │
│  └────────────────┘     └────────────────┘     └──────────────┘  │
│                                                                 │
│                     内层不依赖外层                               │
└─────────────────────────────────────────────────────────────────┘
```

#### 具体引用关系

1. **API 层**：

   - 只能引用 `domain` 层（接口定义）
   - 只能引用 `pkg` 层（公共工具）
   - 不能直接引用 `usecase` 层的具体实现
   - 不能直接引用 `repository` 层的具体实现
   - **示例**：`api/admin/handler/user.go` 通过 `domain.AdminUserInteractor` 接口调用业务逻辑
2. **UseCase 层**：

   - 实现 `domain` 层定义的接口
   - 引用 `domain` 层的实体和接口
   - 不能直接引用 `api` 层
   - 不能直接引用基础设施（如数据库、缓存）
   - **示例**：`usecase/userauth/admin_user.go` 实现了 `domain.AdminUserInteractor` 接口
3. **Domain 层**：

   - 核心业务逻辑层，不依赖任何外部层
   - 定义领域实体、值对象、聚合根
   - 定义领域服务接口
   - 定义仓储接口
   - **示例**：`domain/admin_user.go` 定义了 `AdminUser` 实体和 `AdminUserRepository` 接口
4. **Repository 层**：

   - 实现 `domain` 层定义的仓储接口
   - 引用 `domain` 层的实体
   - 可以引用基础设施（如 `ent`、Redis 等）
   - 不能引用 `usecase` 层和 `api` 层
   - **示例**：`repository/admin_user.go` 实现了 `domain.AdminUserRepository` 接口
5. **Infrastructure 层**：

   - 提供基础设施实现
   - 不依赖任何业务层
   - **示例**：`adapter/objectstorage/objectstorage.go` 提供对象存储服务实现

### 2.3 层边界

#### 职责边界

| 层级              | 主要职责                             | 禁止行为                            | 示例文件                                   |
| ----------------- | ------------------------------------ | ----------------------------------- | ------------------------------------------ |
| API 层            | HTTP 请求处理、参数验证、响应格式化  | 实现业务逻辑、直接操作数据库        | `api/admin/handler/user.go`              |
| UseCase 层        | 业务流程编排、协调领域对象和基础设施 | 直接操作数据库、处理 HTTP 请求/响应 | `usecase/userauth/admin_user.go`         |
| Domain 层         | 核心业务规则、实体定义、领域服务     | 依赖外部框架、直接操作数据库        | `domain/admin_user.go`                   |
| Repository 层     | 数据持久化、查询实现                 | 实现业务逻辑、暴露技术细节给上层    | `repository/admin_user.go`               |
| Infrastructure 层 | 外部系统集成、工具类实现             | 包含业务逻辑、依赖业务层            | `adapter/objectstorage/objectstorage.go` |

#### 数据传输边界

1. **API 层**：

   - 输入：`api/[service]/types/*` 定义的请求类型
   - 输出：`api/[service]/types/*` 定义的响应类型
   - 转换：将请求类型转换为领域实体或值对象
   - **示例**：`api/admin/types/user.go` 定义了 `LoginReq` 和 `LoginResp` 类型
2. **UseCase 层**：

   - 输入：领域实体、值对象或基本类型
   - 输出：领域实体、值对象或基本类型
   - 转换：协调领域对象的交互
   - **示例**：`usecase/userauth/admin_user.go` 的 `Login` 方法接收用户名和密码，返回令牌和过期时间
3. **Domain 层**：

   - 输入：领域实体、值对象或基本类型
   - 输出：领域实体、值对象或基本类型
   - 转换：领域内部的数据转换
   - **示例**：`domain/admin_user.go` 的 `CheckPassword` 方法验证密码
4. **Repository 层**：

   - 输入：领域实体、值对象或基本类型
   - 输出：领域实体、值对象或基本类型
   - 转换：将数据库模型转换为领域实体，反之亦然
   - **示例**：`repository/admin_user.go` 的 `convertAdminUser` 函数将 `ent.AdminUser` 转换为 `domain.AdminUser`

### 2.4 错误边界

#### 错误定义与转换

1. **Domain 层错误**：

   - 定义核心业务错误（如 `ErrTokenInvalid`、`ErrMismatchedHashAndPassword`）
   - 提供通用错误构造函数（如 `NotFoundError`、`ParamsError`、`ConflictError`）
   - 错误类型：自定义 `error` 类型，实现 `Unwrap` 方法支持错误链
   - **示例**：`domain/error.go` 定义了各种领域错误类型
2. **Repository 层错误**：

   - 将技术错误（如数据库错误）转换为业务错误
   - 使用 `domain.NotFoundError` 包装 `ent.IsNotFound` 错误
   - 使用 `domain.ConflictError` 包装 `ent.IsConstraintError` 错误
   - 错误类型：`domain.Error` 或其具体实现
   - **示例**：`repository/admin_user.go` 中处理数据库错误的代码
3. **UseCase 层错误**：

   - 包装底层错误（如仓储错误）
   - 添加上下文信息（如 `fmt.Errorf("failed to find user: %w", err)`）
   - 错误类型：`domain.Error` 或标准 `error`
   - **示例**：`usecase/userauth/admin_user.go` 中处理仓储错误的代码
4. **API 层错误**：

   - 将业务错误转换为 HTTP 错误
   - 使用 `errorx.New` 构造 HTTP 错误
   - 错误类型：`errorx.Error`
   - **示例**：`api/admin/handler/user.go` 中处理业务错误的代码

#### 错误处理最佳实践

1. **错误传递**：

   ```go
   // Repository 层：将技术错误转换为业务错误
   if ent.IsNotFound(err) {
       err = domain.NotFoundError(err)
   }

   // UseCase 层：包装错误并添加上下文
   if err != nil {
       err = fmt.Errorf("failed to find user by username: %w", err)
       return
   }

   // API 层：将业务错误转换为 HTTP 错误
   if domain.IsNotFound(err) {
       translated := i18n.Translate(ctx, errcode.UserNotFound.String(), map[string]any{"Username": req.Username})
       c.Error(errorx.New(http.StatusBadRequest, errcode.UserNotFound, err).WithMessage(translated))
       return
   }
   ```
2. **错误包装**：

   - 始终使用 `%w` 包装底层错误，保留错误链
   - 避免使用 `%v` 或 `%s` 包装错误，导致错误链丢失
3. **错误分类**：

   - 业务错误：由领域层定义，反映业务规则违反
   - 技术错误：由基础设施层产生，需要转换为业务错误
   - HTTP 错误：由 API 层产生，反映 HTTP 协议状态
4. **错误检查**：

   - 使用 `errors.Is` 和 `errors.As` 进行错误检查
   - 领域层提供 `IsNotFound`、`IsParamsError` 等辅助函数
   - **示例**：`domain/error.go` 中的 `IsNotFound` 函数

### 2.5 目录结构

```
├── adapter/         # 适配器层，用于集成外部系统和服务
├── api/             # API 层，按服务划分（admin/backend/customer/frontend/intl）
│   └── [service]/
│       ├── handler/     # 请求处理函数
│       ├── middleware/  # 中间件
│       ├── types/       # 请求和响应类型定义
│       └── app.go       # 服务入口
├── bootstrap/       # 启动配置，初始化各个组件
├── buildinfo/       # 构建信息
├── cmd/             # 命令行入口
├── domain/          # 领域层，核心业务模型和服务
├── ent/             # entgo 生成的代码和数据库模型
│   └── schema/       # 数据库表结构定义
├── etc/             # 配置文件
├── pkg/             # 公共库和工具函数
├── repository/      # 仓储层，数据访问实现
├── scheduler/       # 定时任务
└── usecase/         # 用例层，业务逻辑实现
```

## 3. 编码规范

### 3.1 命名规范

- **包名**：使用小写字母，短而清晰，避免使用下划线

  ```go
  package userauth  // 正确
  package user_auth // 错误
  ```
- **结构体名**：使用 PascalCase，清晰表达其用途

  ```go
  type AdminUserInteractor struct {}  // 正确
  type admin_user_interactor struct {} // 错误
  ```
- **接口名**：使用 PascalCase，以 `er` 结尾表示能力

  ```go
  type AdminUserInteractor interface {}  // 正确
  type AdminUser interface {}           // 错误
  ```
- **变量名**：使用 camelCase，避免使用单字母变量（除了循环变量）

  ```go
  var userID string  // 正确
  var uid string     // 错误
  ```
- **常量名**：使用全部大写，单词间用下划线分隔

  ```go
  const MaxRetries = 3  // 正确
  const max_retries = 3 // 错误
  ```

### 3.2 代码风格

- 遵循 Go 官方编码规范 (gofmt)
- 每行不超过 100 个字符
- 函数参数不超过 5 个，超过时使用结构体
- 使用 `defer` 释放资源
- 使用 `errors.Wrap` 或 `fmt.Errorf` 包装错误，保留错误链

### 3.3 函数设计

- 函数职责单一，不超过 50 行
- 错误处理：函数返回值的第一个参数应为错误
- 使用 context 传递请求上下文和取消信号

  ```go
  func (interactor *AdminUserInteractor) Login(ctx context.Context, username, password string) (token string, expAt time.Time, err error) {
      // 实现逻辑
  }
  ```

### 3.4 依赖注入

- 使用 Uber fx 进行依赖注入
- 模块定义使用 `fx.Provide` 和 `fx.Module`

  ```go
  var Module = fx.Module(
      "usecase",
      fx.Provide(
          userauth.NewAdminUserInteractor,
      ),
  )
  ```

## 4. 数据库设计规范

### 4.1 表结构设计

- 使用 entgo 定义数据库表结构
- 每个表必须包含以下字段（通过 Mixin 实现）：

  - `id`：UUID 类型，主键
  - `created_at`：创建时间
  - `updated_at`：更新时间
  - `deleted_at`：软删除标记（0 表示未删除）
- 表名使用复数形式，字段名使用下划线分隔
- 索引设计：

  - 为经常查询的字段创建索引
  - 联合索引的字段顺序按查询频率排序
  - 软删除字段需包含在唯一索引中

  ```go
  func (AdminUser) Indexes() []ent.Index {
      return []ent.Index{
          index.Fields("username", "deleted_at").Unique(),
      }
  }
  ```

### 4.2 数据操作

- 使用事务保证数据一致性
- 优先使用 entgo 的查询构建器，避免直接写 SQL
- 软删除默认开启，使用 `SkipSoftDelete` 上下文跳过软删除

## 5. API 设计规范

### 5.1 RESTful API 设计

- 使用 HTTP 方法表示操作类型：

  - `GET`：获取资源
  - `POST`：创建资源
  - `PUT`：更新资源
  - `DELETE`：删除资源
- 资源路径使用复数形式

  ```
  GET /api/users      # 获取所有用户
  GET /api/users/1    # 获取单个用户
  POST /api/users     # 创建用户
  ```

### 5.2 请求与响应

- 请求和响应使用 JSON 格式
- 请求参数验证使用 Gin 的 binding 标签

  ```go
  type LoginReq struct {
      Username string `json:"username" binding:"required"`
      Password string `json:"password" binding:"required"`
  }
  ```
- 响应格式统一：

  ```json
  {
      "code": "success",
      "data": {...}
  }
  ```

### 5.3 API 文档

- 使用 Swagger 自动生成 API 文档
- 为每个 API 添加详细的注释，包括：

  - Tags：API 分类
  - Summary：API 简要描述
  - Accept：请求格式
  - Produce：响应格式
  - Param：请求参数
  - Success：成功响应
  - Router：API 路径

  ```go
  // Login
  //
  //	@Tags		用户管理
  //	@Security	BearerAuth
  //	@Summary	用户登录
  //	@Accept		json
  //	@Produce	json
  //	@Param		data	body		types.LoginReq	true	"请求信息"
  //	@Success	200		{object}	types.LoginResp	"成功"
  //	@Router		/user/login [post]
  ```

## 6. 依赖管理规范

- 使用 Go Modules 管理依赖
- Go 版本：1.24.0
- 依赖版本固定，避免使用 latest
- 定期更新依赖，保持安全性和稳定性

## 7. 配置管理规范

### 7.1 配置文件

- 使用 TOML 格式的配置文件
- 配置文件按服务划分（admin.toml/backend.toml/customer.toml/frontend.toml/intl.toml/scheduler.toml）
- 配置项命名使用驼峰式或下划线分隔

### 7.2 环境变量

- 敏感配置（如数据库密码、API 密钥）优先使用环境变量
- 环境变量命名使用大写字母，单词间用下划线分隔

## 8. 构建与部署规范

### 8.1 构建

- 使用 Go 原生构建命令
- 构建信息通过 `buildinfo` 包注入

### 8.2 开发环境

- 使用 Air 进行热重载开发
- 每个服务有独立的 Air 配置文件

### 8.3 部署

- 使用 Docker 容器化部署
- 支持多环境部署（dev/staging/prod）

## 9. 测试规范

### 9.1 测试类型

- **单元测试**：测试单个函数或方法
- **集成测试**：测试模块间的交互
- **端到端测试**：测试完整的业务流程

### 9.2 测试框架

- 使用 Go 标准库的 `testing` 包
- 断言使用 `testify/assert`

### 9.3 测试覆盖率

- 单元测试覆盖率目标：≥60%（根据开发时间酌情考虑)
- 关键业务逻辑必须有测试覆盖（根据开发时间酌情考虑)

## 10. 代码审查规范

### 10.1 审查流程

1. 提交代码到 feature 分支
2. 创建 Merge Request
3. 至少 1 位资深开发者审查
4. 审查通过后合并到主分支

### 10.2 审查要点

- 代码风格是否符合规范
- 业务逻辑是否正确
- 错误处理是否完善
- 性能是否优化
- 安全性是否考虑

## 11. 日志与监控

### 11.1 日志

- 使用 zap 作为日志库
- 日志级别：DEBUG/INFO/WARN/ERROR
- 日志格式：JSON
- 每个请求必须记录请求 ID

### 11.2 监控

- 使用 OpenTelemetry 进行分布式追踪
- 关键指标监控：
  - 请求响应时间
  - 错误率
  - 并发数
  - 数据库连接数

## 12. 错误处理

- 使用自定义错误码
- 错误信息国际化
- 详细错误记录到日志
- 用户友好的错误提示

## 13. 安全规范

- 密码使用 bcrypt 哈希存储
- JWT 令牌认证
- 请求参数验证
- SQL 注入防护
- XSS 防护
- CSRF 防护
- HTTPS 加密传输

## 14. 国际化

- 使用 i18n 包支持多语言
- 错误信息和用户提示国际化
- 语言文件存放在 `etc/language/` 目录

## 15. 开发工具

- **版本控制**：Git
- **代码格式化**：gofmt
- **代码检查**：golint
- **依赖管理**：Go Modules
- **数据库迁移**：Atlas

## 16. 版本管理

- 遵循 Semantic Versioning 规范
- 主版本号：不兼容的 API 变更
- 次版本号：向后兼容的功能新增
- 修订号：向后兼容的问题修复

## 17. 附则

- 本规范适用于所有参与 POS Dine API 开发的团队成员
- 规范将根据项目发展和技术变化定期更新

---

**最后更新时间**：2025-12-17
**版本**：1.0.0

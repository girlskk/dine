# 构建阶段
FROM harbor.jiguang.dev/library/go:1.24 AS builder

WORKDIR /src

# 1. 复制依赖文件（利用缓存）
COPY go.mod go.sum ./
RUN go mod download

# 2. 复制源代码
COPY . .

# 4. 构建参数
ARG ENTRYPOINT_NAME
ARG BUILD_PKG=gitlab.jiguang.dev/pos-dine/dine/buildinfo
ARG VERSION=unknown
ARG BUILD_AT=unknown

# 5. 构建应用
RUN go build \
    -trimpath \
    -ldflags="-w -s -X ${BUILD_PKG}.Version=${VERSION} -X ${BUILD_PKG}.BuildAt=${BUILD_AT}" \
    -o /app \
    ./cmd/${ENTRYPOINT_NAME}

# 运行阶段
FROM harbor.jiguang.dev/library/alpine:latest

# 设置工作目录
WORKDIR /app

# 复制应用程序
COPY --from=builder /app .

# 复制配置文件
ARG ENTRYPOINT_NAME
COPY --from=builder /src/etc/${ENTRYPOINT_NAME}.toml ./etc/config.toml
COPY --from=builder /src/etc/language ./etc/language

# 暴露端口
EXPOSE 8080

# 启动命令
CMD ["./app", "-conf", "./etc/config.toml"]

FROM harbor.jiguang.dev/library/go:1.24 AS base
ARG ENTRYPOINT_NAME
ARG BUILD_PKG=gitlab.jiguang.dev/pos-dine/dine/buildinfo
ARG VERSION=unknown
ARG BUILD_AT=unknown

WORKDIR /src
COPY . .
RUN go mod download -x
RUN go mod tidy
RUN CGO_ENABLED=0 go build -ldflags "-w -s -X ${BUILD_PKG}.Version=${VERSION} -X ${BUILD_PKG}.BuildAt=${BUILD_AT}" -o app ./cmd/${ENTRYPOINT_NAME}

FROM harbor.jiguang.dev/library/alpine
ARG ENTRYPOINT_NAME
WORKDIR /app
COPY --from=base /src/app .
COPY --from=base /src/etc/${ENTRYPOINT_NAME}.toml ./etc/config.toml
COPY --from=base /src/etc/language ./etc/language
CMD ["./app", "-conf", "./etc/config.toml"]

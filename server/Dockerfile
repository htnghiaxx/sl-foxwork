# Dockerfile cho Mattermost Server Development

# Stage 1: Webapp build stage
FROM node:20-alpine AS webapp-builder

# Cài đặt dependencies cần thiết cho build webapp
RUN apk add --no-cache \
    python3 \
    make \
    g++ \
    git \
    bash \
    autoconf \
    automake \
    libtool \
    pkgconf \
    nasm \
    build-base \
    gifsicle

# Lưu ý: Docker build context nên là root repo và dùng -f server/Dockerfile
WORKDIR /workspace/webapp

# Copy package files trước để cache deps
COPY webapp/package*.json ./
COPY webapp/channels/package*.json ./channels/
COPY webapp/patches/ ./patches/
COPY webapp/platform/ ./platform/

# Cài deps cho webapp
RUN npm install

# Copy source và build
COPY webapp/ ./
RUN npm run build


# Stage 2: Server build stage
FROM golang:1.25.1-alpine AS builder

# Cài đặt các dependencies cần thiết
RUN apk add --no-cache \
    git \
    make \
    gcc \
    musl-dev \
    ca-certificates

# Tạo thư mục làm việc
WORKDIR /app/server

# Bật chế độ workspace để Go nhận go.work
ENV GOWORK=auto

# Copy go.mod, go.sum và go.work từ thư mục server để tối ưu cache deps và đồng bộ workspace
COPY server/go.mod server/go.sum server/go.work ./

# Copy source code của server
COPY server/ ./

# Đồng bộ workspace rồi tải dependencies sau khi toàn bộ source đã có
RUN go work sync && \
    go mod download

# Build Mattermost server binary theo cấu hình air.toml
# cmd = "go build -o ./tmp/mattermost-dev ./cmd/mattermost"
RUN mkdir -p tmp && \
    go build -o ./tmp/mattermost-dev ./cmd/mattermost

# Stage 2: Runtime stage
FROM alpine:3.18

# Cài đặt runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl

# Tạo user mattermost
RUN addgroup -g 2000 mattermost && \
    adduser -D -u 2000 -G mattermost mattermost

# Tạo thư mục cần thiết
RUN mkdir -p /mattermost/data /mattermost/logs /mattermost/plugins /mattermost/client /mattermost/client/plugins && \
    chown -R mattermost:mattermost /mattermost

# Copy binary từ builder stage
COPY --from=builder --chown=mattermost:mattermost /app/server/tmp/mattermost-dev /mattermost/bin/mattermost


# Copy i18n translation files
COPY --chown=mattermost:mattermost server/i18n/ /mattermost/i18n/

# Copy templates directory
COPY --chown=mattermost:mattermost server/templates/ /mattermost/templates/

# Copy fonts directory
COPY --chown=mattermost:mattermost server/fonts/ /mattermost/fonts/

# Copy webapp built files vào thư mục client của server
COPY --from=webapp-builder --chown=mattermost:mattermost /workspace/webapp/channels/dist/ /mattermost/client/

# Set working directory
WORKDIR /mattermost

# Switch to mattermost user
USER mattermost

# Expose ports
EXPOSE 8065 8067 8074 8075

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8065/api/v4/system/ping || exit 1

# Default command với args từ air.toml
# args_bin = ["--config", "config/config.json"]
CMD ["/mattermost/bin/mattermost"]

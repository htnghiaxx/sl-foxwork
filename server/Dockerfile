# Dockerfile cho Mattermost Server Development

# Stage 1: Build stage
FROM golang:1.25.1-alpine AS builder

# Cài đặt các dependencies cần thiết
RUN apk add --no-cache \
    git \
    make \
    gcc \
    musl-dev \
    ca-certificates

# Tạo thư mục làm việc
WORKDIR /app

# Bật chế độ workspace để Go nhận go.work
ENV GOWORK=auto

# Copy go.mod, go.sum và go.work trước để tối ưu cache deps và đồng bộ workspace
COPY go.mod go.sum ./

# Copy source code
COPY . .

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
RUN mkdir -p /mattermost/data /mattermost/logs /mattermost/config /mattermost/plugins /mattermost/client/plugins && \
    chown -R mattermost:mattermost /mattermost

# Copy binary từ builder stage
COPY --from=builder --chown=mattermost:mattermost /app/tmp/mattermost-dev /mattermost/bin/mattermost

# Copy config files nếu có
COPY --chown=mattermost:mattermost config/ /mattermost/config/

# Copy i18n translation files
COPY --chown=mattermost:mattermost i18n/ /mattermost/i18n/

# Copy templates directory
COPY --chown=mattermost:mattermost templates/ /mattermost/templates/

# Copy fonts directory
COPY --chown=mattermost:mattermost fonts/ /mattermost/fonts/

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
CMD ["/mattermost/bin/mattermost", "--config", "config/config.json"]

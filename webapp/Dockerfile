# Dockerfile cho Mattermost Webapp Development
# Dựa trên package.json và build scripts

# Stage 1: Build stage
FROM node:20-alpine AS builder

# Cài đặt dependencies cần thiết cho build
RUN apk add --no-cache \
    python3 \
    make \
    g++ \
    gcc \
    libc-dev \
    libpng-dev \
    nasm \
    git \
    patch \
    autoconf \
    automake \
    libtool \
    pkgconfig \
    bash \
    giflib-dev

# Tạo thư mục làm việc
WORKDIR /app

# Copy package files để cache dependencies
COPY package*.json ./
COPY channels/package*.json ./channels/
COPY platform/*/package*.json ./platform/*/
COPY patches/* ./patches/

# Install dependencies (bao gồm devDependencies để chạy patch-package)
RUN npm ci

# Copy source code
COPY . .

# Set build environment variables
# MM_SERVICESETTINGS_SITEURL cần thiết cho webpack config để xác định publicPath
ARG MM_SERVICESETTINGS_SITEURL=http://localhost:8065
ENV MM_SERVICESETTINGS_SITEURL=${MM_SERVICESETTINGS_SITEURL}

# Build webapp theo cấu hình package.json
# Script: "build": "node scripts/build.mjs"
RUN npm run build

# Stage 2: Production stage với nginx
FROM nginx:alpine

# Cài đặt curl cho health check
RUN apk add --no-cache curl

# Copy built files từ builder stage
COPY --from=builder /app/channels/dist /usr/share/nginx/html

# Copy nginx config (tạo file cấu hình nginx tùy chỉnh nếu cần)
# COPY nginx.conf /etc/nginx/nginx.conf

# Tạo thư mục cho static files
RUN mkdir -p /usr/share/nginx/html/static

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]


#!/bin/bash

# Script để chạy Mattermost development environment
# Dựa trên cấu hình air.toml và package.json

echo "🚀 Starting Mattermost Development Environment..."

# Kiểm tra Docker và Docker Compose
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Tạo thư mục cần thiết
echo "📁 Creating necessary directories..."
mkdir -p server/data server/logs server/config server/plugins server/client/plugins

# Build và chạy containers
echo "🔨 Building and starting containers..."
docker-compose -f docker-compose.dev.yml up --build -d

# Kiểm tra trạng thái containers
echo "⏳ Waiting for services to start..."
sleep 10

# Health check
echo "🏥 Checking service health..."
if curl -f http://localhost:8065/api/v4/system/ping &> /dev/null; then
    echo "✅ Mattermost Server is running at http://localhost:8065"
else
    echo "❌ Mattermost Server is not responding"
fi

if curl -f http://localhost:8080 &> /dev/null; then
    echo "✅ Mattermost Webapp is running at http://localhost:8080"
else
    echo "❌ Mattermost Webapp is not responding"
fi

echo ""
echo "🎉 Development environment is ready!"
echo "📊 Server: http://localhost:8065"
echo "🌐 Webapp: http://localhost:8080"
echo "🗄️  Database: localhost:5432"
echo ""
echo "📝 Useful commands:"
echo "  - View logs: docker-compose -f docker-compose.dev.yml logs -f"
echo "  - Stop services: docker-compose -f docker-compose.dev.yml down"
echo "  - Rebuild: docker-compose -f docker-compose.dev.yml up --build"

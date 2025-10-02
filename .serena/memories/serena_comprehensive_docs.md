# Tài Liệu Serena Toàn Diện - Dự Án Mattermost

## Tổng Quan Dự Án

### Giới Thiệu
Mattermost là một nền tảng cộng tác tự lưu trữ (self-hosted collaboration platform) được xây dựng với Go (backend) và React/TypeScript (frontend). Đây là một monorepo với nhiều thành phần hoạt động cùng nhau.

### Mục Đích Dự Án
- Cung cấp nền tảng chat và messaging
- Tự động hóa quy trình làm việc
- Gọi thoại và chia sẻ màn hình
- Tích hợp AI
- Công cụ cộng tác nhóm

## Cấu Trúc Dự Án

### Thư Mục Chính
```
sl-foxwork/
├── server/          # Go backend server
├── webapp/          # React/TypeScript frontend  
├── api/             # OpenAPI documentation
├── e2e-tests/       # End-to-end testing
├── tools/           # Development tools
└── .serena/         # Serena configuration
```

### Server (Backend)
- **Ngôn ngữ**: Go 1.24.5+
- **Cơ sở dữ liệu**: PostgreSQL (chính), MySQL, SQLite
- **Tìm kiếm**: Elasticsearch/OpenSearch
- **Cache**: Redis
- **Lưu trữ file**: Local filesystem, S3-compatible, AWS S3
- **Xác thực**: LDAP, SAML, OAuth2, OpenID Connect

### Web App (Frontend)
- **Ngôn ngữ**: TypeScript 5.6.3
- **Framework**: React 17.0.2
- **State Management**: Redux 4.2.0
- **Build Tool**: Webpack 5.95.0
- **Styling**: Styled Components, SCSS

## Hướng Dẫn Phát Triển

### Yêu Cầu Hệ Thống
- **Go**: 1.24.5+ (tối thiểu 1.15)
- **Node.js**: 18.10.0+
- **npm**: 9.0.0+ hoặc 10.0.0+
- **Docker**: Cho các dịch vụ phát triển local
- **Make**: Hệ thống build
- **Git**: Quản lý phiên bản

### Thiết Lập Môi Trường Phát Triển

#### 1. Thiết Lập Server
```bash
cd server
make setup-go-work          # Thiết lập Go workspace
make prepackaged-binaries   # Tải các binary cần thiết
make start-docker          # Khởi động các dịch vụ (PostgreSQL, Redis, etc.)
make run-server            # Khởi động server
```

#### 2. Thiết Lập Web App
```bash
cd webapp
make node_modules          # Cài đặt dependencies
make run                   # Khởi động webpack dev server
```

#### 3. Môi Trường Phát Triển Đầy Đủ
```bash
# Từ thư mục server
make run                   # Khởi động cả server và webapp
```

### Lệnh Phát Triển Quan Trọng

#### Server Commands
```bash
# Phát triển
make run-server            # Chỉ khởi động server
make run-client            # Chỉ khởi động webapp
make run                   # Khởi động cả server và webapp
make dev                   # Khởi động với hot reload

# Testing
make test-server           # Chạy tất cả test server
make test-server-quick     # Chạy test nhanh
make test-server-race      # Chạy test với race detection

# Code Quality
make check-style           # Kiểm tra style code
make vet                   # Chạy Go vet checks
make golangci-lint         # Chạy golangci-lint
```

#### Web App Commands
```bash
# Phát triển
make run                   # Khởi động webpack dev server
make dev                   # Khởi động với webpack-dev-server
make stop                  # Dừng webpack

# Testing
make test                  # Chạy Jest tests
make test:watch            # Chạy test ở chế độ watch
make test-ci               # Chạy test cho CI

# Code Quality
make check-style           # Chạy ESLint và Stylelint
make fix-style             # Sửa các vấn đề linting
make check-types           # Chạy TypeScript checks
```

### Docker Services

#### Các Dịch Vụ Có Sẵn
- `postgres` - Cơ sở dữ liệu PostgreSQL
- `minio` - Lưu trữ file tương thích S3
- `inbucket` - Testing email
- `openldap` - Xác thực LDAP
- `elasticsearch` - Công cụ tìm kiếm
- `redis` - Cache và job queue
- `keycloak` - Identity provider

#### Docker Commands
```bash
make start-docker          # Khởi động các dịch vụ cần thiết
make stop-docker           # Dừng các dịch vụ
make clean-docker          # Xóa containers và volumes
make update-docker         # Cập nhật container images
```

## Hướng Dẫn Testing

### Server Testing
- **Unit Tests**: `make test-server`
- **Quick Tests**: `make test-server-quick`
- **Race Detection**: `make test-server-race`
- **Coverage**: Đặt `ENABLE_COVERAGE=true` để có báo cáo coverage

### Web App Testing
- **Unit Tests**: `make test`
- **E2E Tests**: Xem thư mục `/e2e-tests/`
- **Cypress**: `cd e2e-tests/cypress && npm test`
- **Playwright**: `cd e2e-tests/playwright && npm test`

## Hướng Dẫn API

### API Documentation
- **Vị trí**: `/api/` directory
- **Format**: YAML files sử dụng OpenAPI standard
- **Tools**: ReDoc document generator

### Key API Endpoints
- **Authentication**: `/api/v4/users/login`
- **Channels**: `/api/v4/channels`
- **Posts**: `/api/v4/posts`
- **Users**: `/api/v4/users`
- **Teams**: `/api/v4/teams`

## Hướng Dẫn Triển Khai

### Local Development
```bash
# Khởi động môi trường phát triển
make run

# Dừng môi trường
make stop
```

### Production Deployment
```bash
# Build production
make build                 # Build server binary
make dist                  # Tạo distribution package

# Web app build
cd webapp
make dist                  # Build production bundle
```

### Docker Deployment
```bash
# Sử dụng Docker Compose
docker-compose up -d

# Hoặc sử dụng Make
make start-docker
```

## Quy Ước Code Style

### Go (Server) Style
- **Formatting**: Sử dụng `gofmt` và `goimports`
- **Naming**: Theo quy ước Go (camelCase cho private, PascalCase cho public)
- **Comments**: Sử dụng Go documentation comments cho exported functions/types
- **Error Handling**: Luôn xử lý errors một cách rõ ràng

### TypeScript/React (Web App) Style
- **Formatting**: Sử dụng Prettier
- **Linting**: ESLint với custom Mattermost rules
- **TypeScript**: Strict type checking được bật
- **Components**: Sử dụng functional components với hooks

## Troubleshooting

### Các Vấn Đề Thường Gặp
```bash
# Xung đột port
lsof -i :8065              # Kiểm tra port có đang được sử dụng
lsof -i :8066              # Kiểm tra port webapp

# Vấn đề Docker
docker system prune        # Dọn dẹp Docker
docker compose down -v     # Xóa volumes

# Vấn đề quyền
sudo chown -R $(whoami) .  # Sửa quyền file
```

### Debug Commands
```bash
# Go profiling
go tool pprof http://localhost:8065/debug/pprof/profile

# Memory profiling
go tool pprof http://localhost:8065/debug/pprof/heap

# Web app bundle analysis
cd webapp && npm run stats
```

## Environment Variables

### Server Environment
```bash
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_SQLSETTINGS_DRIVERNAME=postgres
export MM_SQLSETTINGS_DATASOURCE="postgres://mmuser:mostest@localhost/mattermost_test?sslmode=disable"
export MM_SERVICESETTINGS_ENABLELOCALMODE=true
```

### Web App Environment
```bash
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_SERVICESETTINGS_LISTENADDRESS=:8065
```

## Lệnh Tham Khảo Nhanh

### Lệnh Sử Dụng Hàng Ngày
```bash
# Phát triển hàng ngày
make run                   # Khởi động mọi thứ
make stop                  # Dừng mọi thứ
make test-server           # Test server
make test                  # Test webapp (từ webapp/)

# Code quality
make check-style           # Kiểm tra tất cả code
make fix-style             # Sửa các vấn đề webapp

# Docker
make start-docker          # Khởi động dịch vụ
make stop-docker           # Dừng dịch vụ
```

## Tài Liệu Bổ Sung

### Memory Files Hiện Có
1. `project_summary.md` - Tóm tắt dự án
2. `project_overview.md` - Tổng quan dự án
3. `tech_stack.md` - Chi tiết công nghệ
4. `development_workflow.md` - Quy trình phát triển
5. `code_style_conventions.md` - Quy ước code style
6. `suggested_commands.md` - Lệnh tham khảo

### Tài Liệu Bên Ngoài
- **API Documentation**: `/api/` directory
- **Contributing**: Xem CONTRIBUTING.md
- **Developer Docs**: https://developers.mattermost.com/

## Kết Luận

Tài liệu này cung cấp hướng dẫn toàn diện để phát triển với dự án Mattermost. Nó bao gồm thiết lập môi trường, các lệnh phát triển, testing, triển khai và troubleshooting. Sử dụng các memory files hiện có để có thêm thông tin chi tiết về từng khía cạnh cụ thể của dự án.
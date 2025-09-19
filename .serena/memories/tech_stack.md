# Technology Stack

## Backend (Server)
- **Language**: Go 1.24.5+
- **Database**: PostgreSQL (primary), MySQL, SQLite support
- **Search**: Elasticsearch/OpenSearch
- **Cache**: Redis
- **File Storage**: Local filesystem, S3-compatible (MinIO), AWS S3
- **Authentication**: LDAP, SAML, OAuth2, OpenID Connect
- **Message Queue**: Redis-based job queue
- **Web Server**: Built-in HTTP server (Gorilla mux)

### Key Go Dependencies
- `github.com/gorilla/mux` - HTTP router
- `github.com/gorilla/websocket` - WebSocket support
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/elastic/go-elasticsearch/v8` - Elasticsearch client
- `github.com/redis/rueidis` - Redis client
- `github.com/prometheus/client_golang` - Metrics
- `github.com/sirupsen/logrus` - Logging

## Frontend (Web App)
- **Language**: TypeScript 5.6.3
- **Framework**: React 17.0.2
- **State Management**: Redux 4.2.0
- **Build Tool**: Webpack 5.95.0
- **Package Manager**: npm/yarn
- **Styling**: Styled Components 5.3.7, SCSS
- **Testing**: Jest 29.7.0, Testing Library

### Key Frontend Dependencies
- `react` & `react-dom` - UI framework
- `redux` & `react-redux` - State management
- `react-router-dom` - Routing
- `styled-components` - CSS-in-JS
- `@mui/material` - UI components
- `react-intl` - Internationalization
- `monaco-editor` - Code editor
- `chart.js` - Charts and graphs

## Development Tools
- **Linting**: ESLint, golangci-lint
- **Formatting**: Prettier, gofmt
- **Testing**: Jest (frontend), Go testing (backend)
- **E2E Testing**: Cypress, Playwright
- **Build**: Make, npm scripts
- **Docker**: Docker Compose for local development

## Database
- **Primary**: PostgreSQL 12+
- **Migrations**: Custom migration system with Go
- **Backup**: Built-in backup/restore functionality

## Infrastructure
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Reverse Proxy**: Nginx (recommended)
- **Monitoring**: Prometheus metrics, Grafana dashboards
- **Logging**: Structured logging with multiple backends

## Mobile & Desktop
- **Mobile**: React Native apps (separate repositories)
- **Desktop**: Electron-based desktop apps
- **Cross-platform**: Shared codebase with platform-specific implementations
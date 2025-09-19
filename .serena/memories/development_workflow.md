# Development Workflow

## Prerequisites
- **Go**: 1.24.5+ (minimum 1.15)
- **Node.js**: 18.10.0+
- **npm**: 9.0.0+ or 10.0.0+
- **Docker**: For local development services
- **Make**: Build system
- **Git**: Version control

## Local Development Setup

### 1. Server Setup
```bash
cd server
make setup-go-work          # Set up Go workspace
make prepackaged-binaries   # Download required binaries
make start-docker          # Start required services (PostgreSQL, Redis, etc.)
make run-server            # Start the server
```

### 2. Web App Setup
```bash
cd webapp
make node_modules          # Install dependencies
make run                   # Start webpack dev server
```

### 3. Full Development Environment
```bash
# From server directory
make run                   # Starts both server and webapp
```

## Development Commands

### Server Commands
```bash
# Development
make run-server            # Start server only
make run-client            # Start webapp only
make run                   # Start both server and webapp
make dev                   # Start with hot reload (using air)

# Testing
make test-server           # Run all server tests
make test-server-quick     # Run quick tests only
make test-server-race      # Run tests with race detection
make test-mmctl            # Run mmctl tests

# Code Quality
make check-style           # Run all style checks
make vet                   # Run Go vet checks
make golangci-lint         # Run golangci-lint
make modernize             # Run modernize linter

# Build
make build                 # Build server binary
make dist                  # Create distribution package
```

### Web App Commands
```bash
# Development
make run                   # Start webpack dev server
make dev                   # Start with webpack-dev-server
make stop                  # Stop webpack

# Testing
make test                  # Run Jest tests
make test:watch            # Run tests in watch mode
make test-ci               # Run tests for CI

# Code Quality
make check-style           # Run ESLint and Stylelint
make fix-style             # Fix linting issues
make check-types           # Run TypeScript checks

# Build
make dist                  # Build production bundle
make package               # Create webapp package
```

## Docker Services

### Available Services
- `postgres` - PostgreSQL database
- `minio` - S3-compatible file storage
- `inbucket` - Email testing
- `openldap` - LDAP authentication
- `elasticsearch` - Search engine
- `redis` - Caching and job queue
- `keycloak` - Identity provider

### Docker Commands
```bash
make start-docker          # Start required services
make stop-docker           # Stop services
make clean-docker          # Remove containers and volumes
make update-docker         # Update container images
```

## Testing

### Server Testing
- **Unit Tests**: `make test-server`
- **Quick Tests**: `make test-server-quick`
- **Race Detection**: `make test-server-race`
- **Coverage**: Set `ENABLE_COVERAGE=true` for coverage reports

### Web App Testing
- **Unit Tests**: `make test`
- **E2E Tests**: See `/e2e-tests/` directory
- **Cypress**: `cd e2e-tests/cypress && npm test`
- **Playwright**: `cd e2e-tests/playwright && npm test`

## Code Generation

### Server
```bash
make mocks                 # Generate mock files
make store-layers          # Generate store layers
make pluginapi             # Generate plugin API
make gen-serialized        # Generate serialization methods
```

### Web App
```bash
npm run i18n-extract       # Extract translation strings
npm run make-emojis        # Generate emoji data
```

## Configuration

### Server Configuration
- Default config: `server/config/config.json`
- Override: Create `server/config.override.mk`
- Environment variables: Set in shell or `.env` file

### Web App Configuration
- Webpack config: `webapp/channels/webpack.config.js`
- Babel config: `webapp/channels/babel.config.js`
- TypeScript config: `webapp/channels/tsconfig.json`

## Debugging

### Server Debugging
```bash
make debug-server          # Start with delve debugger
make debug-server-headless # Start headless debugger for IDE
```

### Web App Debugging
- Use browser dev tools
- Webpack dev server with source maps
- React DevTools extension recommended
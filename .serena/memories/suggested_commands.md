# Suggested Commands for Mattermost Development

## Essential Development Commands

### Quick Start
```bash
# Set up development environment
cd server
make setup-go-work
make prepackaged-binaries
make start-docker
make run

# In another terminal
cd webapp
make run
```

### Daily Development Workflow
```bash
# Start development environment
cd server
make run                    # Start both server and webapp

# Stop development environment
make stop                   # Stop everything
make stop-docker           # Stop only Docker services
```

## Server Development Commands

### Development
```bash
# Basic development
make run-server            # Start server only
make run-client            # Start webapp only
make dev                   # Start with hot reload

# Debugging
make debug-server          # Start with delve debugger
make debug-server-headless # Headless debugger for IDE

# Configuration
make config-reset          # Reset to default config
make config-ldap           # Configure LDAP
make config-saml           # Configure SAML
make config-openid         # Configure OpenID
```

### Testing
```bash
# Run tests
make test-server           # All server tests
make test-server-quick     # Quick tests only
make test-server-race      # Race condition tests
make test-mmctl            # mmctl tests

# Test with coverage
ENABLE_COVERAGE=true make test-server

# Specific test packages
go test ./channels/app/...  # Test specific package
```

### Code Quality
```bash
# Style checks
make check-style           # All style checks
make vet                   # Go vet checks
make golangci-lint         # Linting
make modernize             # Modernize linter

# Fix issues
make fix-style             # Fix webapp linting issues
```

### Build and Package
```bash
# Build
make build                 # Build server binary
make dist                  # Create distribution
make package               # Package for distribution

# Platform-specific builds
make build-linux           # Linux build
make build-osx             # macOS build
make build-windows         # Windows build
```

## Web App Development Commands

### Development
```bash
cd webapp
make run                   # Start webpack dev server
make dev                   # Start with webpack-dev-server
make stop                  # Stop webpack
```

### Testing
```bash
# Run tests
make test                  # All tests
make test:watch            # Watch mode
make test-ci               # CI mode with coverage

# Specific test files
npm test -- --testNamePattern="UserProfile"
npm test -- --testPathPattern="components"
```

### Code Quality
```bash
# Style checks
make check-style           # ESLint and Stylelint
make check-types           # TypeScript checks

# Fix issues
make fix-style             # Fix linting issues
```

### Build
```bash
make dist                  # Production build
make package               # Create webapp package
```

## Docker Commands

### Service Management
```bash
# Start services
make start-docker          # Start required services
make update-docker         # Update container images

# Stop services
make stop-docker           # Stop services
make clean-docker          # Remove containers and volumes

# Specific services
docker compose up postgres redis  # Start specific services
```

### Service Configuration
```bash
# Enable additional services
ENABLED_DOCKER_SERVICES="postgres minio elasticsearch" make start-docker

# Disable Docker entirely
MM_NO_DOCKER=true make run-server
```

## E2E Testing Commands

### Cypress Testing
```bash
cd e2e-tests/cypress
npm install
npm test                   # Run all tests
npm run test:smoke         # Run smoke tests only

# Specific tests
node run_tests.js --include-file="login_spec.js"
```

### Playwright Testing
```bash
cd e2e-tests/playwright
npm install
npm test                   # Run all tests
npx playwright test --grep "login"  # Run specific tests
```

## Code Generation Commands

### Server Code Generation
```bash
# Generate mocks
make mocks                 # All mocks
make store-mocks           # Store mocks only
make plugin-mocks          # Plugin mocks only

# Generate other code
make store-layers          # Store layers
make pluginapi             # Plugin API
make gen-serialized        # Serialization methods
```

### Web App Code Generation
```bash
cd webapp
npm run i18n-extract       # Extract translation strings
npm run make-emojis        # Generate emoji data
```

## Utility Commands

### Database
```bash
# Create new migration
make new-migration name=add_user_table

# Run migrations
go run ./cmd/mattermost migrate

# Database utilities
bin/mmctl config set SqlSettings.DriverName postgres
```

### Logging and Debugging
```bash
# View logs
tail -f mattermost.log
tail -f mattermost.log.jsonl

# Debug specific components
MM_LOG_LEVEL=DEBUG make run-server
```

### Cleanup Commands
```bash
# Clean build artifacts
make clean                 # Clean everything
make nuke                  # Clean including data

# Clean specific components
cd webapp && make clean    # Clean webapp only
go clean ./...             # Clean Go build cache
```

## Git and Version Control

### Git Workflow
```bash
# Create feature branch
git checkout -b feature/MM-12345-description

# Commit changes
git add .
git commit -m "feat(component): add new feature"

# Push and create PR
git push origin feature/MM-12345-description
```

### Version Management
```bash
# Check Go version
go version

# Check Node version
node --version
npm --version

# Update dependencies
make update-dependencies   # Update Go dependencies
cd webapp && npm update    # Update npm dependencies
```

## Troubleshooting Commands

### Common Issues
```bash
# Port conflicts
lsof -i :8065              # Check if port is in use
lsof -i :8066              # Check webapp port

# Docker issues
docker system prune        # Clean up Docker
docker compose down -v     # Remove volumes

# Permission issues
sudo chown -R $(whoami) .  # Fix file permissions
```

### Performance Debugging
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

## Quick Reference

### Most Used Commands
```bash
# Daily development
make run                   # Start everything
make stop                  # Stop everything
make test-server           # Test server
make test                  # Test webapp (from webapp/)

# Code quality
make check-style           # Check all code
make fix-style             # Fix webapp issues

# Docker
make start-docker          # Start services
make stop-docker           # Stop services
```
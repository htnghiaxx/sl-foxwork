# Mattermost Project Summary

## Quick Overview
Mattermost is a self-hosted collaboration platform built with Go (backend) and React/TypeScript (frontend). It's organized as a monorepo with multiple components working together.

## Project Structure
```
sl-foxwork/
├── server/          # Go backend server
├── webapp/          # React/TypeScript frontend
├── api/             # OpenAPI documentation
├── e2e-tests/       # End-to-end testing
└── tools/           # Development tools
```

## Key Technologies
- **Backend**: Go 1.24.5+, PostgreSQL, Redis, Elasticsearch
- **Frontend**: React 17, TypeScript 5.6, Redux, Webpack
- **Testing**: Jest, Cypress, Playwright
- **DevOps**: Docker, Make, GitHub Actions

## Getting Started
1. **Prerequisites**: Go 1.24.5+, Node.js 18.10+, Docker
2. **Setup**: `cd server && make setup-go-work && make prepackaged-binaries`
3. **Start**: `make start-docker && make run`
4. **Web App**: `cd webapp && make run`

## Development Workflow
- **Server**: `make run-server` (Go backend)
- **Web App**: `make run` (React frontend)
- **Full Stack**: `make run` (both server and webapp)
- **Testing**: `make test-server` and `make test` (webapp)
- **Code Quality**: `make check-style`

## Key Directories
- `server/channels/` - Core business logic
- `server/platform/` - Shared platform services
- `webapp/channels/src/` - Main React application
- `webapp/platform/` - Shared platform components
- `e2e-tests/` - End-to-end tests

## Build System
- **Server**: Make-based with Go modules
- **Web App**: npm workspace with webpack
- **Docker**: Compose for local development
- **CI/CD**: GitHub Actions

## Testing Strategy
- **Unit Tests**: Go testing, Jest
- **Integration Tests**: Server integration tests
- **E2E Tests**: Cypress and Playwright
- **Code Quality**: ESLint, golangci-lint, TypeScript

## Documentation
- **API**: OpenAPI/Swagger in `/api/`
- **Contributing**: See CONTRIBUTING.md
- **Developer Docs**: https://developers.mattermost.com/

## Key Commands Reference
```bash
# Development
make run                   # Start everything
make stop                  # Stop everything
make test-server           # Test server
make test                  # Test webapp

# Code Quality
make check-style           # Check all code
make fix-style             # Fix webapp issues

# Docker
make start-docker          # Start services
make stop-docker           # Stop services
```

## Memory Files Created
1. `project_overview.md` - Project purpose and structure
2. `tech_stack.md` - Technology stack details
3. `development_workflow.md` - Development commands and workflow
4. `code_style_conventions.md` - Code style and conventions
5. `suggested_commands.md` - Essential commands reference
6. `project_summary.md` - This summary file

This documentation provides a comprehensive guide for understanding and working with the Mattermost codebase.
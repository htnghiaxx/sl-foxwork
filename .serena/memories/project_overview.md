# Mattermost Project Overview

## Project Purpose
Mattermost is an open core, self-hosted collaboration platform that offers:
- Chat and messaging
- Workflow automation
- Voice calling and screen sharing
- AI integration
- Team collaboration tools

## Project Structure
The project is organized into several main components:

### Server (`/server/`)
- **Language**: Go (v1.24.5+)
- **Purpose**: Backend server implementation
- **Key directories**:
  - `channels/`: Core business logic (API, app, database, jobs, store, web, etc.)
  - `cmd/`: Command-line tools (mattermost, mmctl)
  - `platform/`: Shared platform services
  - `public/`: Public API and models
  - `config/`: Configuration management
  - `einterfaces/`: Enterprise interfaces

### Web App (`/webapp/`)
- **Language**: TypeScript/React
- **Purpose**: Frontend web application
- **Key directories**:
  - `channels/`: Main web application
  - `platform/`: Shared platform components (client, components, types, redux)

### API Documentation (`/api/`)
- **Purpose**: OpenAPI documentation for Mattermost APIs
- **Format**: YAML files using OpenAPI standard
- **Tools**: ReDoc document generator

### E2E Tests (`/e2e-tests/`)
- **Tools**: Cypress and Playwright
- **Purpose**: End-to-end testing for web client

## Build System
- **Server**: Make-based build system with Go modules
- **Web App**: npm/yarn workspace with webpack
- **Docker**: Docker Compose for local development
- **CI/CD**: GitHub Actions for automated builds and deployment

## Key Features
- Self-hosted deployment
- Multi-platform support (web, mobile, desktop)
- Plugin architecture
- Enterprise features (when licensed)
- Open source core with MIT license
- Monthly releases on the 16th
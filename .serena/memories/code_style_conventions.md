# Code Style and Conventions

## Go (Server) Style Guidelines

### General Conventions
- **Formatting**: Use `gofmt` and `goimports` for formatting
- **Naming**: Follow Go naming conventions (camelCase for private, PascalCase for public)
- **Comments**: Use proper Go documentation comments for exported functions/types
- **Error Handling**: Always handle errors explicitly, use `if err != nil` pattern

### Code Organization
- **Packages**: Organize code into logical packages
- **Interfaces**: Define interfaces close to where they're used
- **Dependencies**: Use dependency injection where possible
- **Testing**: Write tests in the same package with `_test.go` suffix

### Linting Rules
- **golangci-lint**: Primary linter with custom configuration
- **govet**: Custom Mattermost govet with specific checks:
  - `structuredLogging` - Enforce structured logging
  - `inconsistentReceiverName` - Consistent receiver naming
  - `emptyStrCmp` - Avoid empty string comparisons
  - `tFatal` - Proper test failure handling
  - `configtelemetry` - Configuration telemetry
  - `errorAssertions` - Error assertion patterns
  - `requestCtxNaming` - Request context naming
  - `license` - License header enforcement

### File Structure
```
server/
├── channels/           # Core business logic
│   ├── api4/          # REST API handlers
│   ├── app/           # Application layer
│   ├── store/         # Data access layer
│   └── web/           # Web handlers
├── platform/          # Shared platform services
├── public/            # Public API and models
└── cmd/               # Command-line tools
```

## TypeScript/React (Web App) Style Guidelines

### General Conventions
- **Formatting**: Use Prettier for code formatting
- **Linting**: ESLint with custom Mattermost rules
- **TypeScript**: Strict type checking enabled
- **Imports**: Use absolute imports from `src/` directory

### React Conventions
- **Components**: Use functional components with hooks
- **Props**: Define interfaces for component props
- **State**: Use Redux for global state, local state for component-specific data
- **Styling**: Use styled-components for CSS-in-JS

### File Naming
- **Components**: PascalCase (e.g., `UserProfile.tsx`)
- **Utilities**: camelCase (e.g., `formatDate.ts`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `API_ENDPOINTS.ts`)
- **Types**: PascalCase with `.types.ts` suffix

### Code Organization
```
webapp/channels/src/
├── components/        # Reusable UI components
├── actions/          # Redux actions
├── reducers/         # Redux reducers
├── selectors/        # Redux selectors
├── utils/            # Utility functions
├── types/            # TypeScript type definitions
└── constants/        # Application constants
```

### Linting Rules
- **ESLint**: Custom Mattermost ESLint plugin
- **Stylelint**: SCSS/CSS linting
- **TypeScript**: Strict type checking
- **Import Rules**: Enforce import order and grouping

## Testing Conventions

### Go Testing
- **Test Files**: `*_test.go` files in same package
- **Test Functions**: `TestFunctionName` format
- **Benchmarks**: `BenchmarkFunctionName` format
- **Test Data**: Use `testlib` package for test utilities
- **Mocks**: Generate mocks using `mockery`

### React Testing
- **Test Files**: `*.test.tsx` or `*.test.ts` files
- **Testing Library**: Use React Testing Library
- **Jest**: Primary testing framework
- **Snapshots**: Use for component regression testing
- **Mocking**: Mock external dependencies and APIs

## Documentation Conventions

### Go Documentation
- **Package Comments**: Describe package purpose
- **Function Comments**: Start with function name
- **Example Code**: Use `Example` functions for documentation
- **API Documentation**: Use OpenAPI/Swagger annotations

### React Documentation
- **Component Comments**: Use JSDoc format
- **Props Documentation**: Document all props with types
- **Storybook**: Use for component documentation (if available)

## Git Conventions

### Commit Messages
- **Format**: `type(scope): description`
- **Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- **Scope**: Component or module affected
- **Description**: Clear, concise description

### Branch Naming
- **Feature**: `feature/MM-12345-description`
- **Bugfix**: `bugfix/MM-12345-description`
- **Hotfix**: `hotfix/MM-12345-description`

## Code Review Guidelines

### Review Checklist
- [ ] Code follows style guidelines
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No breaking changes without proper migration
- [ ] Performance implications considered
- [ ] Security implications reviewed

### Review Process
1. Create pull request with clear description
2. Request review from appropriate team members
3. Address feedback and update code
4. Ensure CI checks pass
5. Merge after approval

## Performance Guidelines

### Go Performance
- **Memory**: Avoid unnecessary allocations
- **Concurrency**: Use goroutines appropriately
- **Database**: Optimize queries and use proper indexing
- **Caching**: Use Redis for frequently accessed data

### React Performance
- **Rendering**: Use React.memo for expensive components
- **State**: Minimize unnecessary re-renders
- **Bundle**: Code splitting and lazy loading
- **Images**: Optimize images and use appropriate formats
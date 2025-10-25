# Contributing to Ham Radio QSO Journal

Thank you for your interest in contributing to Ham Radio QSO Journal! We welcome contributions from the community and are excited to work with you.

## 🎯 How to Contribute

### Reporting Bugs
If you find a bug, please create an issue with:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Go version, etc.)

### Suggesting Features
We love new ideas! For feature requests:
- Describe the feature and use case
- Explain why it would be useful to other hams
- Consider if it aligns with project goals

### Code Contributions
1. **Fork** the repository
2. **Clone** your fork:
   ```bash
   git clone https://github.com/RD2W/rd2w-log.git
   cd rd2w-log
   ```
3. **Create a feature branch**:
   ```bash
   git checkout -b feature/amazing-feature
   ```

## 🛠 Development Setup

### Prerequisites
- Go 1.25.3 or later
- Docker and Docker Compose
- Git

### Local Development
1. **Start dependencies**:
   ```bash
   docker-compose up -d postgres redis
   ```

2. **Run services individually**:
   ```bash
   # Auth service
   go run cmd/auth-service/main.go

   # QSO service  
   go run cmd/qso-service/main.go

   # Analytics service
   go run cmd/analytics-service/main.go

   # API Gateway
   go run cmd/api-gateway/main.go
   ```

3. **Run tests**:
   ```bash
   go test ./... -race -count=1
   ```

## 📝 Code Standards

### Go Code
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Write unit tests for new functionality
- Use meaningful variable and function names

### Project Structure
- Follow Clean Architecture principles
- Keep business logic in `usecase` layer
- Data access in `repository` layer  
- Delivery mechanisms in `delivery` layer

### Commit Messages
Use conventional commit format:
```
feat: add ADIF export functionality
fix: resolve memory leak in QSO cache
docs: update API documentation
test: add unit tests for auth service
```

## 🔧 Testing

### Running Tests
```bash
# All tests
go test ./... -race -count=1

# Tests with coverage
go test ./... -race -cover

# Specific package
go test ./internal/auth/... -v
```

### Test Standards
- Write tests for new features
- Maintain >80% test coverage
- Use table-driven tests when appropriate
- Mock external dependencies

## 📋 Pull Request Process

1. **Ensure tests pass** and code is properly formatted
2. **Update documentation** if needed
3. **Add changelog entry** if applicable
4. **Create PR** with clear description and references to issues
5. **Request review** from maintainers

### PR Checklist
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] Code follows project standards
- [ ] Commit messages follow conventions
- [ ] All checks pass (CI/CD)

## 🏗 Project Architecture

### Key Directories
- `cmd/` - Application entry points
- `internal/` - Private application code
- `pkg/` - Public reusable packages
- `proto/` - Protocol Buffer definitions
- `deployments/` - Docker and deployment configs

### Service Communication
- Use gRPC for inter-service communication
- Protocol Buffers for API contracts
- REST API via gRPC Gateway

## 🐛 Common Issues

### Database Migrations
```bash
# Create new migration
goose create add_user_preferences sql

# Run migrations
goose postgres "user=hamuser dbname=hamradio sslmode=disable" up
```

### Protobuf Generation
```bash
# Generate Go code from .proto files
protoc --go_out=. --go-grpc_out=. proto/**/*.proto
```

## 📞 Getting Help

- Create a [GitHub Issue](https://github.com/rd2w/rd2w-log/issues)
- Join our [Discussions](https://github.com/rd2w/rd2w-log/discussions)
- Check existing documentation in `/docs`

---

Thank you for contributing to the amateur radio community! 🎯

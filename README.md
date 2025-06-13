# Core API Service

A Go-based API service built with modern technologies and best practices.

## Features

- RESTful API endpoints
- PostgreSQL database integration with GORM
- Configuration management with Viper
- Docker support for development and deployment
- Graceful shutdown handling
- Generic repository pattern implementation
- Health check endpoint
- Structured logging with Zap

## Installation

### System Requirements

- Go 1.22 or higher
- Docker and Docker Compose
- PostgreSQL
- Make (optional, for using Makefile commands)

### Required Tools Installation

```bash
# Install oapi-codegen
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

## Project Structure

```
.
├── api/                    # Directory containing OpenAPI/Swagger files
│   ├── open-api.yaml      # Main OpenAPI file
│   └── user/              # User module API directory
│       ├── api.yaml       # User API OpenAPI file
│       └── config.yaml    # oapi-codegen configuration
        └── iml.go          # implement business logic
        └── server.go       #
├── cmd/                   # Entry points directory
├── models/              # Models/structs
├── services/             
│   ├── database/         # Database interaction layer
│   ├── log/           # Models/structs
│   ├── wire/     # Database interaction layer
├── helpers/                  # Reusable code
├── main.go              # Main entry point
└── Makefile             # Build and development commands
```

## New API Creation Process

1. **Create OpenAPI/Swagger Files**
    - Create a new directory in `api/` for the new module
    - Create `api.yaml` to define API endpoints
    - Create `config.yaml` for oapi-codegen

2. **Generate Code**
   ```bash
   # Generate code from OpenAPI spec
   make generate-user-api
   
   # Bundle OpenAPI docs
   make swagger
   ```

3. **Implement Business Logic**
    - Implement business logic in `api/`
    - Create models in `models/`
    - Implement repository in `services/database/`

4. **Wire Dependencies**
    - Update wire configuration in `services/wire/`

5. **Run Service**
   ```bash
   make run
   ```
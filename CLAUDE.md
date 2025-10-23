# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Squire Project is a multi-component agent middleware system that integrates with Azure AI Agents. It consists of:

1. **AgentMiddleware.Server** - ASP.NET Core 9.0 Web API backend
2. **agent-cli** - Go-based CLI client with Microsoft Authentication Library (MSAL) integration

## Architecture

### Authentication System
- **Azure Entra ID (formerly Azure AD)** for authentication
- Server uses **Microsoft.Identity.Web** (v4.0.0) for JWT token validation
- No local password storage or manual JWT generation
- User entity stores `EntraObjectId` (required), `Email`, and `Username` (cached from Azure tokens)
- Client authenticates via MSAL device code flow with automatic token caching and refresh

### Database Layer
- **PostgreSQL** database for persistence
- **Entity Framework Core 9.0** with Npgsql provider
- **ApplicationDbContext** (AgentMiddleware.Server/Data/ApplicationDbContext.cs) manages three core entities:
  - **User**: Stores Azure Entra ObjectId with unique index, plus cached email/username
  - **UserProfile**: 1:1 relationship with User, stores agent instructions and custom workflows as JSON
  - **ConversationMetadata**: Many:1 with User, tracks Azure AI thread conversations
- All relationships use cascade delete
- Indexed on User.EntraObjectId and ConversationMetadata.AzureThreadId

### Azure Integration
- Uses **Azure.AI.Agents.Persistent** (v1.1.0) for AI agent interactions
- **Azure.Storage.Blobs** (v12.26.0) for blob storage
- Configuration requires `AZURE_PROJECT_ENDPOINT` and `AZURE_MODEL_DEPLOYMENT` environment variables

## Development Commands

### .NET Backend

```bash
# Build entire solution
dotnet build SquireMiddleware.sln

# Run the server locally
dotnet run --project AgentMiddleware.Server

# Build specific project
dotnet build AgentMiddleware.Server/AgentMiddleware.Server.csproj
```

### Database Migrations

```bash
# Add new migration
dotnet ef migrations add <MigrationName> --project AgentMiddleware.Server

# Apply migrations
dotnet ef database update --project AgentMiddleware.Server

# Remove last migration
dotnet ef migrations remove --project AgentMiddleware.Server
```

### Docker Compose

```bash
# Start all services (PostgreSQL + Server)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild and start
docker-compose up -d --build
```

The docker-compose setup:
- PostgreSQL on port 5432 with health checks
- Server on ports 5000 (HTTP) and 5001 (HTTPS)
- Auto-applies EF migrations on startup

### Go CLI Client

```bash
# From agent-cli directory:

# Build for current platform
make build

# Build for all platforms
make build-all

# Run directly
make run
# or
go run .

# Install to $GOPATH/bin
make install

# Clean build artifacts
make clean
```

## Configuration

### Server Configuration (appsettings.json)

Required environment variables for deployment:
- `ConnectionStrings__DefaultConnection` - PostgreSQL connection string
- `Azure__ProjectEndpoint` - Azure AI project endpoint
- `Azure__ModelDeployment` - Azure model deployment name
- `AzureAd__TenantId` - Azure Entra ID tenant ID
- `AzureAd__ClientId` - Application client ID registered in Azure
- `AzureAd__Audience` - Expected JWT token audience (typically `api://{ClientId}`)

Development defaults:
- Database: `agentmiddleware` on localhost:5432
- Azure Entra Instance: `https://login.microsoftonline.com/`

### Go CLI Configuration

The CLI uses MSAL device code flow for authentication:
- Client ID and Tenant ID must match the server's AzureAd configuration
- Update `clientID` and `tenantID` constants in `agent-cli/main.go`
- Token scope: `api://{clientID}/access_as_user`
- Tokens are cached in-memory by MSAL (persistent cache not yet implemented)
- Server URL defaults to `http://localhost:5000` (configurable in main.go)

## Project Structure

```
/
├── AgentMiddleware.Server/        # ASP.NET Core Web API
│   ├── Controllers/               # API endpoints
│   ├── Data/                      # EF Core DbContext
│   ├── Models/                    # Entity models (User, UserProfile, ConversationMetadata)
│   ├── Migrations/                # EF Core migrations
│   └── Dockerfile                 # Multi-stage Docker build
├── agent-cli/                     # Go CLI client
│   ├── auth/                      # MSAL authentication (device code flow)
│   ├── client/                    # HTTP client for server API
│   ├── config/                    # Configuration management
│   ├── main.go                    # CLI entry point and chat loop
│   └── Makefile                   # Build commands
└── docker-compose.yml             # Local development stack
```

## Key Dependencies

### .NET (AgentMiddleware.Server)
- ASP.NET Core 9.0
- Microsoft.Identity.Web 4.0.0 (Azure Entra ID authentication)
- Azure.AI.Agents.Persistent 1.1.0
- Entity Framework Core 9.0 with PostgreSQL
- Npgsql.EntityFrameworkCore.PostgreSQL 9.0.4

### Go (agent-cli)
- MSAL for Go (github.com/AzureAD/microsoft-authentication-library-for-go)
- Standard library HTTP client

## Important Implementation Notes

### User Authentication Flow
1. User runs CLI client
2. Client authenticates via MSAL device code flow (or silent auth if cached)
3. Client sends requests to server with `Authorization: Bearer {token}` header
4. Server validates JWT against Azure Entra ID using Microsoft.Identity.Web
5. Server extracts user claims (Object ID, email, username) from validated token
6. Server creates or retrieves User entity based on EntraObjectId

### Extracting User Identity in Controllers
Use Microsoft.Identity.Web extension methods on `ClaimsPrincipal`:
- `User.GetObjectId()` - Gets the Azure Entra Object ID (unique identifier)
- `User.GetDisplayName()` - Gets the user's display name
- `User.Claims.FirstOrDefault(c => c.Type == "preferred_username")?.Value` - Gets email/UPN

### Database Schema Considerations
- `User.EntraObjectId` is the primary unique identifier for users (not Username)
- `Username` and `Email` are nullable and cached from tokens (may be updated on each login)
- Do not use `Username` as a unique identifier or foreign key
- All user lookups should be done via `EntraObjectId`

## Dockerfile Notes

The Dockerfile in AgentMiddleware.Server references the deleted AgentMiddleware.Contracts project and will fail to build. Before using Docker:

1. Update `AgentMiddleware.Server/Dockerfile` to remove the line:
   ```dockerfile
   COPY ["AgentMiddleware.Contracts/AgentMiddleware.Contracts.csproj", "AgentMiddleware.Contracts/"]
   ```
2. Remove the corresponding `RUN dotnet restore` reference to Contracts

## CLI Commands

The agent-cli supports these interactive commands:
- `/help` - Display available commands
- `/exit` or `/quit` - Exit the application
- `/logout` - Sign out from Azure Entra and exit
- `/new` - Start a new conversation thread
- `/history` - List past conversations
- `/profile` - View user profile (agent instructions, workflows)
- `/whoami` - Show currently authenticated user

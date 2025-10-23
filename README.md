# Squire - Personalized AI Agent Platform

A personalized AI agent platform that allows users to interact with Azure AI agents through a command-line interface while maintaining individual profiles, custom workflows, and conversation history. Built with enterprise-grade security using Azure Entra ID authentication.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Prerequisites](#prerequisites)
- [Architecture](#architecture)
- [Initial Setup](#initial-setup)
- [Development Workflow](#development-workflow)
- [Configuration](#configuration)
- [Database Management](#database-management)
- [Authentication Flow](#authentication-flow)
- [API Documentation](#api-documentation)
- [Building and Deployment](#building-and-deployment)
- [Common Development Tasks](#common-development-tasks)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [Security Considerations](#security-considerations)
- [License](#license)

## Overview

**Squire** provides a personalized AI agent experience with the following features:

- **Cross-platform CLI** - Single binary runs on Windows, Mac, and Linux
- **User Profiles** - Each user maintains custom agent instructions and workflow preferences
- **Conversation History** - All interactions are tracked and retrievable
- **Enterprise Security** - Azure Entra ID authentication with multi-factor auth support
- **Scalable Backend** - Handles multiple concurrent users with isolated data

### Core Use Cases

1. Developers who want an AI coding assistant with memory of their preferences
2. Teams needing shared access to AI agents with individual personalization
3. Organizations requiring audited, secure AI interactions with centralized management
4. Users who want their AI assistant to follow specific workflows or have domain knowledge

## Quick Start

Get up and running in 5 minutes:

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/squire.git
cd squire

# 2. Copy and configure environment variables
cp .env.example .env
# Edit .env with your Azure credentials (see Initial Setup section)

# 3. Start backend services (PostgreSQL + .NET middleware)
docker compose up --build

# 4. In a new terminal, run the CLI client
cd agent-cli
go run .
```

You should see authentication prompts followed by the agent chat interface.

## Prerequisites

### Required Software

- **.NET 9 SDK** - [Download](https://dotnet.microsoft.com/download/dotnet/9.0)
- **Go 1.21+** - [Download](https://go.dev/dl/)
- **Docker Desktop** - [Download](https://www.docker.com/products/docker-desktop) (or Docker Engine + Docker Compose)
- **Azure CLI** - [Install](https://docs.microsoft.com/cli/azure/install-azure-cli)
- **Git** - [Download](https://git-scm.com/downloads)

### Required Azure Resources

You'll need an active Azure subscription and the following resources:

1. **Azure AI Foundry Project** with a deployed model (e.g., GPT-4o)
   - Get endpoint URL and deployment name
2. **Azure Entra ID App Registration** for the middleware server
   - Configure API permissions: `User.Read`
   - Generate client secret (for server authentication)
3. **(Optional)** Separate app registration for the Go CLI client
   - Recommended for production deployments

### Verify Prerequisites

```bash
# Check installed versions
dotnet --version          # Should be 9.0.x
go version               # Should be 1.21 or higher
docker --version         # Any recent version
docker compose version   # V2 recommended
az --version            # Azure CLI
```

## Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│  User                                                           │
│   │                                                             │
│   │ (1) Runs CLI                                               │
│   ▼                                                             │
│  ┌──────────────────┐                                          │
│  │  Go CLI Client   │                                          │
│  │  (agent-cli/)    │                                          │
│  └────────┬─────────┘                                          │
│           │                                                     │
│           │ (2) Device Code Auth                               │
│           ▼                                                     │
│  ┌──────────────────────┐                                      │
│  │  Azure Entra ID      │ (Managed by Microsoft)               │
│  │  (Identity Provider) │                                      │
│  └────────┬─────────────┘                                      │
│           │                                                     │
│           │ (3) JWT Access Token                               │
│           ▼                                                     │
│  ┌──────────────────────────┐                                  │
│  │  .NET Middleware Server  │                                  │
│  │  (AgentMiddleware.Server)│                                  │
│  │  - Validates JWT         │                                  │
│  │  - User Management       │                                  │
│  │  - Profile Storage       │                                  │
│  └────┬──────────────┬──────┘                                  │
│       │              │                                          │
│       │ (4)          │ (5)                                     │
│       ▼              ▼                                          │
│  ┌──────────┐  ┌─────────────────────┐                        │
│  │PostgreSQL│  │ Azure AI Foundry    │ (Managed by Microsoft) │
│  │ Database │  │ Agent Service       │                        │
│  │          │  │ - Agent Runtime     │                        │
│  │- Profiles│  │ - Thread Management │                        │
│  │- Metadata│  │ - Message History   │                        │
│  └──────────┘  └─────────────────────┘                        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Technology Stack

**Go CLI Client:**
- Go 1.24.3
- MSAL for Go (Azure authentication)
- Standard library HTTP client
- Cross-platform binary compilation

**.NET Middleware Server:**
- .NET 9
- ASP.NET Core Web API
- Entity Framework Core 9
- Microsoft.Identity.Web 4.0.0
- Azure SDK for .NET

**Database:**
- PostgreSQL 16 with Npgsql provider

**Azure Services:**
- Azure AI Foundry (Agent Service)
- Azure Entra ID (Authentication)
- (Optional) Azure Container Apps (Deployment)
- (Optional) Azure Database for PostgreSQL (Production)

### Project Structure

```
/
├── agent-cli/                    # Go CLI client
│   ├── main.go                   # Entry point and chat loop
│   ├── auth/                     # Azure Entra authentication
│   │   └── auth.go              # MSAL device code flow
│   ├── client/                   # HTTP API client
│   │   └── client.go            # API communication
│   ├── config/                   # Configuration management
│   ├── Makefile                  # Build commands
│   └── go.mod                    # Go dependencies
│
├── AgentMiddleware.Server/       # .NET middleware API
│   ├── Controllers/              # HTTP endpoints
│   ├── Models/                   # Database entities
│   │   ├── User.cs              # User entity (Azure Entra ID mapped)
│   │   ├── UserProfile.cs       # User preferences
│   │   └── ConversationMetadata.cs
│   ├── Data/                     # EF Core DbContext
│   │   └── ApplicationDbContext.cs
│   ├── Migrations/               # Database migrations
│   ├── Program.cs                # Application startup
│   ├── appsettings.json          # Configuration template
│   └── Dockerfile                # Container build
│
├── docker-compose.yml            # Local dev orchestration
├── .env.example                  # Environment variable template
├── .gitignore                    # Git ignore rules
├── SquireMiddleware.sln          # .NET solution file
├── CLAUDE.md                     # AI assistant guidelines
└── README.md                     # This file
```

## Initial Setup

### 1. Azure AI Foundry Setup

1. Navigate to [Azure AI Foundry](https://ai.azure.com/)
2. Create a new project or select an existing one
3. Deploy a model (e.g., GPT-4o):
   - Go to **Deployments** → **Create deployment**
   - Select your model and give it a deployment name
4. Copy your **Project Endpoint** (looks like `https://your-project.openai.azure.com/`)
5. Copy your **Deployment Name** (e.g., `gpt-4o`)

### 2. Azure Entra ID App Registration

#### Server App Registration

1. Go to [Azure Portal](https://portal.azure.com/) → **Entra ID** → **App registrations**
2. Click **New registration**:
   - Name: `Squire Middleware API`
   - Supported account types: `Accounts in this organizational directory only`
   - Redirect URI: Leave blank for now
3. After creation, note the **Application (client) ID** and **Directory (tenant) ID**
4. Go to **Expose an API**:
   - Set Application ID URI: `api://{your-client-id}`
   - Add a scope: `access_as_user` (Admin consent required: No)
5. Go to **API permissions**:
   - Add `Microsoft Graph` → `User.Read` (Delegated)
   - Grant admin consent
6. Go to **Certificates & secrets**:
   - Create a new client secret
   - Copy the secret value (you won't see it again)

#### CLI App Registration (Optional but Recommended)

For production, create a separate public client app:
1. **New registration**:
   - Name: `Squire CLI Client`
   - Supported account types: Same as server
   - Redirect URI: `Public client/native` → `http://localhost`
2. Go to **Authentication**:
   - Enable **Allow public client flows**
3. Go to **API permissions**:
   - Add your server API: `api://{server-client-id}/access_as_user`
4. Note the **Client ID** for this app

For development, you can use the same app registration for both server and CLI.

### 3. Local Environment Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/squire.git
cd squire

# Copy environment template
cp .env.example .env

# Edit .env file with your Azure credentials
nano .env
```

Update `.env` with your values:
```bash
AZURE_PROJECT_ENDPOINT=https://your-project.openai.azure.com/
AZURE_MODEL_DEPLOYMENT=gpt-4o
AZURE_AD_TENANT_ID=your-tenant-id
AZURE_AD_CLIENT_ID=your-server-client-id
AZURE_AD_AUDIENCE=api://your-server-client-id
```

### 4. Update CLI Configuration

Edit `agent-cli/main.go` to set your client credentials:

```go
const (
    clientID  = "your-cli-client-id"  // Or use server client ID for dev
    tenantID  = "your-tenant-id"
    serverURL = "http://localhost:5000"
)
```

### 5. Install Dependencies

```bash
# .NET dependencies (restored automatically on build)
dotnet restore

# Go dependencies
cd agent-cli
go mod download
cd ..
```

## Development Workflow

### Starting the Development Environment

```bash
# Terminal 1: Start PostgreSQL and .NET middleware
docker compose up --build

# Wait for "Application started" message, then in Terminal 2:
cd agent-cli
go run .
```

You should see:
```
╔════════════════════════════════════════════════════════════╗
║          Agent Middleware Console Client v1.0              ║
╚════════════════════════════════════════════════════════════╝

Authenticating with Microsoft...
To sign in, use a web browser to open the page https://microsoft.com/devicelogin
and enter the code XXXXXXXXX to authenticate.
```

Follow the prompts to authenticate.

### Making Changes

**Backend (.NET) Changes:**

The Docker setup includes automatic rebuild. Stop (`Ctrl+C`) and restart:
```bash
docker compose up --build
```

For faster iteration during active development:
```bash
# Run outside Docker
cd AgentMiddleware.Server
dotnet run
```

**CLI (Go) Changes:**

Simply stop (`Ctrl+C`) and re-run:
```bash
cd agent-cli
go run .
```

**Database Schema Changes:**

```bash
# Create a new migration
cd AgentMiddleware.Server
dotnet ef migrations add YourMigrationName

# Apply to running database
dotnet ef database update
```

If using Docker, restart the containers to apply migrations automatically.

### Hot Reload

- **.NET**: Not enabled by default in Docker. Use `dotnet watch run` for local development.
- **Go**: Manual restart required (very fast with `go run`).

## Configuration

### Environment Variables

All environment variables can be set in `.env` file (for Docker) or exported in your shell.

| Variable | Description | Example |
|----------|-------------|---------|
| `AZURE_PROJECT_ENDPOINT` | Azure AI Foundry project endpoint | `https://my-project.openai.azure.com/` |
| `AZURE_MODEL_DEPLOYMENT` | Model deployment name | `gpt-4o` |
| `AZURE_AD_TENANT_ID` | Azure Entra tenant ID | `00000000-0000-0000-0000-000000000000` |
| `AZURE_AD_CLIENT_ID` | Server app registration client ID | `00000000-0000-0000-0000-000000000000` |
| `AZURE_AD_AUDIENCE` | Expected JWT audience | `api://00000000-0000-0000-0000-000000000000` |
| `ConnectionStrings__DefaultConnection` | PostgreSQL connection string | `Host=postgres;Database=agentmiddleware;...` |

### Configuration Files

**Server Configuration** (`AgentMiddleware.Server/appsettings.json`):
- Contains default settings for logging, database, and Azure services
- Environment-specific settings override these (e.g., `appsettings.Development.json`)
- Secrets should be in environment variables, not checked into source control

**CLI Configuration** (`agent-cli/main.go`):
- `clientID`, `tenantID`, `serverURL` constants
- MSAL token cache stored in memory (no persistent cache yet)

## Database Management

### Entity Framework Core Commands

All EF commands should be run from the `AgentMiddleware.Server` directory:

```bash
cd AgentMiddleware.Server

# Create a new migration
dotnet ef migrations add MigrationName

# Apply migrations to database
dotnet ef database update

# Rollback to a specific migration
dotnet ef database update PreviousMigrationName

# Remove the last migration (if not applied)
dotnet ef migrations remove

# Generate SQL script for migration
dotnet ef migrations script
```

### Database Schema

**Users Table:**
- `Id` (int, primary key)
- `EntraObjectId` (string, unique) - Azure Entra user object ID
- `Email` (string, nullable) - Cached from token
- `Username` (string, nullable) - Cached from token
- `CreatedAt` (datetime)

**UserProfiles Table:**
- `Id` (int, primary key)
- `UserId` (int, foreign key)
- `PreferredAgentInstructions` (text, nullable)
- `CustomWorkflowsJson` (text, nullable)
- `UpdatedAt` (datetime)

**ConversationMetadata Table:**
- `Id` (int, primary key)
- `UserId` (int, foreign key)
- `AzureThreadId` (string) - Thread ID in Azure AI Foundry
- `Title` (string, nullable)
- `CreatedAt` (datetime)
- `LastMessageAt` (datetime)

### Accessing the Database Directly

```bash
# Using Docker
docker exec -it agent-postgres psql -U postgres -d agentmiddleware

# Common queries
\dt                          # List tables
SELECT * FROM "Users";       # View users
SELECT * FROM "UserProfiles"; # View profiles
\q                           # Quit
```

### Resetting the Database

```bash
# Stop containers
docker compose down

# Remove PostgreSQL volume
docker volume rm squire_postgres_data

# Restart (will recreate database and run migrations)
docker compose up --build
```

## Authentication Flow

### Device Code Flow (Step-by-Step)

1. **User Runs CLI**
   ```bash
   cd agent-cli && go run .
   ```

2. **CLI Checks for Cached Token**
   - MSAL library checks in-memory cache for valid token
   - If valid token exists, skip to step 7

3. **CLI Initiates Device Code Flow**
   - Calls Azure Entra ID `/devicecode` endpoint
   - Receives device code and user code

4. **User Authenticates in Browser**
   ```
   To sign in, use a web browser to open the page
   https://microsoft.com/devicelogin and enter the code ABC123
   ```
   - User visits URL and enters code
   - Completes authentication (username, password, MFA if required)
   - Azure issues JWT access token

5. **CLI Receives Access Token**
   - Token is valid for 1 hour
   - MSAL caches token in memory
   - Token includes claims: user object ID, email, name

6. **CLI Sends Authenticated Requests**
   ```http
   GET /api/profile HTTP/1.1
   Host: localhost:5000
   Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGc...
   ```

7. **Server Validates Token**
   - Microsoft.Identity.Web validates signature against Azure public keys
   - Checks token expiration, audience, issuer
   - Extracts user identity from claims

8. **Server Retrieves/Creates User Profile**
   ```csharp
   var objectId = User.GetObjectId();  // From JWT claims
   var user = await db.Users.FirstOrDefaultAsync(u => u.EntraObjectId == objectId);
   if (user == null) {
       user = new User {
           EntraObjectId = objectId,
           Email = User.GetEmail(),
           Username = User.GetDisplayName()
       };
       db.Users.Add(user);
       await db.SaveChangesAsync();
   }
   ```

9. **Subsequent Requests Use Cached Token**
   - No re-authentication needed for 1 hour
   - After expiration, silent refresh or device code flow repeats

### Token Contents

JWT tokens contain claims like:
```json
{
  "oid": "00000000-0000-0000-0000-000000000000",
  "preferred_username": "user@company.com",
  "name": "John Doe",
  "aud": "api://your-client-id",
  "exp": 1234567890,
  "iss": "https://login.microsoftonline.com/{tenant-id}/v2.0"
}
```

Inspect tokens at [jwt.io](https://jwt.io) for debugging.

## API Documentation

### Base URL

Local development: `http://localhost:5000`

### Authentication

All endpoints require `Authorization: Bearer {token}` header.

### Endpoints

#### Send Message to Agent

```http
POST /api/agent/send
Content-Type: application/json

{
  "message": "What is the weather today?",
  "threadId": "thread_abc123" // Optional, omit for new conversation
}
```

**Response:**
```json
{
  "message": "I can help you check the weather...",
  "threadId": "thread_abc123"
}
```

#### List Conversations

```http
GET /api/agent/conversations
```

**Response:**
```json
[
  {
    "threadId": "thread_abc123",
    "title": "Weather Discussion",
    "createdAt": "2025-10-22T10:30:00Z",
    "lastMessageAt": "2025-10-22T10:35:00Z"
  }
]
```

#### Get User Profile

```http
GET /api/profile
```

**Response:**
```json
{
  "preferredAgentInstructions": "Always respond concisely",
  "customWorkflowsJson": "{\"workflow1\": {...}}"
}
```

#### Update User Profile

```http
PUT /api/profile
Content-Type: application/json

{
  "preferredAgentInstructions": "Always respond concisely",
  "customWorkflowsJson": "{\"workflow1\": {...}}"
}
```

### Swagger UI

When running locally, access interactive API documentation at:
```
http://localhost:5000/swagger
```

(Note: Swagger UI may need to be enabled in `Program.cs` first)

## Building and Deployment

### Building the Go CLI

```bash
cd agent-cli

# Build for your current platform
make build
# Output: bin/agent-cli

# Build for all platforms
make build-all
# Output:
#   bin/agent-cli-linux-amd64
#   bin/agent-cli-darwin-amd64
#   bin/agent-cli-darwin-arm64 (Apple Silicon)
#   bin/agent-cli-windows-amd64.exe

# Install to $GOPATH/bin
make install
```

Distribute the binaries to users for their platform.

### Building the .NET Middleware

**Local Docker Build:**
```bash
# Build image
docker build -t agent-middleware:latest -f AgentMiddleware.Server/Dockerfile .

# Run locally
docker run -p 5000:8080 \
  -e ConnectionStrings__DefaultConnection="Host=host.docker.internal;Database=agentmiddleware;..." \
  -e AZURE_PROJECT_ENDPOINT="..." \
  -e AZURE_MODEL_DEPLOYMENT="..." \
  agent-middleware:latest
```

**Push to Azure Container Registry:**
```bash
# Login to ACR
az acr login --name yourregistry

# Tag image
docker tag agent-middleware:latest yourregistry.azurecr.io/agent-middleware:v1.0.0

# Push image
docker push yourregistry.azurecr.io/agent-middleware:v1.0.0
```

### Deploying to Azure Container Apps

1. **Create Azure Database for PostgreSQL**
   ```bash
   az postgres flexible-server create \
     --resource-group myResourceGroup \
     --name squire-db \
     --location eastus \
     --admin-user adminuser \
     --admin-password SecurePassword123! \
     --sku-name Standard_B1ms
   ```

2. **Create Container App Environment**
   ```bash
   az containerapp env create \
     --name squire-env \
     --resource-group myResourceGroup \
     --location eastus
   ```

3. **Deploy Container App**
   ```bash
   az containerapp create \
     --name squire-middleware \
     --resource-group myResourceGroup \
     --environment squire-env \
     --image yourregistry.azurecr.io/agent-middleware:v1.0.0 \
     --target-port 8080 \
     --ingress external \
     --env-vars \
       AZURE_PROJECT_ENDPOINT="..." \
       AZURE_MODEL_DEPLOYMENT="..." \
       AZURE_AD_TENANT_ID="..." \
       AZURE_AD_CLIENT_ID="..." \
       AZURE_AD_AUDIENCE="..." \
       ConnectionStrings__DefaultConnection="Host=squire-db.postgres.database.azure.com;..."
   ```

4. **Update CLI Configuration**

   Update `serverURL` in `agent-cli/main.go` to point to your Container App URL:
   ```go
   serverURL = "https://squire-middleware.azurecontainerapps.io"
   ```

## Common Development Tasks

### Adding a New API Endpoint

1. **Create Service Method** (if needed)
   ```csharp
   // AgentMiddleware.Server/Services/AgentService.cs
   public async Task<string> GetAgentStatus()
   {
       // Implementation
   }
   ```

2. **Add Controller Action**
   ```csharp
   // AgentMiddleware.Server/Controllers/AgentController.cs
   [HttpGet("status")]
   [Authorize]
   public async Task<IActionResult> GetStatus()
   {
       var status = await _agentService.GetAgentStatus();
       return Ok(new { status });
   }
   ```

3. **Update Go Client**
   ```go
   // agent-cli/client/client.go
   func (c *AgentClient) GetStatus(ctx context.Context) (string, error) {
       // Implementation
   }
   ```

4. **Add CLI Command** (if needed)
   ```go
   // agent-cli/main.go
   case "/status":
       status, err := apiClient.GetStatus(ctx)
       if err != nil {
           return err
       }
       fmt.Printf("Agent status: %s\n", status)
       return nil
   ```

### Adding a New User Preference

1. **Update UserProfile Entity**
   ```csharp
   // AgentMiddleware.Server/Models/UserProfile.cs
   public string? PreferredLanguage { get; set; }
   ```

2. **Create and Apply Migration**
   ```bash
   cd AgentMiddleware.Server
   dotnet ef migrations add AddPreferredLanguage
   dotnet ef database update
   ```

3. **Update ProfileController**
   ```csharp
   // Controller already handles all UserProfile properties via model binding
   // No changes needed if using standard PUT /api/profile
   ```

4. **Add CLI Command**
   ```go
   case "/setlang":
       if len(parts) < 2 {
           fmt.Println("Usage: /setlang <language>")
           return nil
       }
       // Update profile with preferred language
   ```

### Debugging Tips

**Attach Debugger to Containerized .NET App:**

1. Update `docker-compose.yml`:
   ```yaml
   services:
     server:
       environment:
         - ASPNETCORE_ENVIRONMENT=Development
       ports:
         - "5000:8080"
         - "5001:8081"
         - "5002:5002"  # Debugger port
   ```

2. In VS Code, add launch configuration:
   ```json
   {
     "name": "Attach to Docker",
     "type": "coreclr",
     "request": "attach",
     "processId": "${command:pickRemoteProcess}"
   }
   ```

**View PostgreSQL Data:**
```bash
docker exec -it agent-postgres psql -U postgres -d agentmiddleware
```

**Inspect JWT Tokens:**
- Copy token from network request
- Paste at [jwt.io](https://jwt.io)
- Verify claims: `oid`, `aud`, `exp`, `iss`

**Common Error Messages:**

- `401 Unauthorized` - Token expired or invalid. Run `/logout` in CLI and re-authenticate.
- `Connection refused (localhost:5000)` - Docker containers not running. Run `docker compose up`.
- `A migration is pending` - Run `dotnet ef database update` or restart Docker.
- `AADSTS700016: Application not found` - Client ID mismatch. Verify `.env` and `main.go` values.

## Testing

### Unit Tests

*Note: Unit test projects not yet implemented. To add:*

```bash
# Create test project
dotnet new xunit -n AgentMiddleware.Tests
dotnet sln add AgentMiddleware.Tests

# Run tests
dotnet test
```

### Integration Testing

**Test Authentication Locally:**

1. Ensure server is running (`docker compose up`)
2. Run CLI (`cd agent-cli && go run .`)
3. Complete device code authentication
4. Verify you see: `✓ Authenticated as: your-email@company.com`

**Test Agent Interaction:**

```bash
# In CLI
You: Hello, what is 2+2?
Agent: [Should respond with calculation]

# Test conversation persistence
/history
# Should show your conversation

# Test new conversation
/new
You: Another message
# Should create new thread
```

**Test Profile Management:**

```bash
# View profile
/profile

# Update profile (implement /setprofile command or use API directly):
curl -X PUT http://localhost:5000/api/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"preferredAgentInstructions": "Be concise"}'
```

### Testing Against Azure AI Foundry

Verify your Azure connection:

```bash
# Check environment variables are set
docker compose exec server printenv | grep AZURE

# Check logs for Azure connection
docker compose logs server | grep -i azure
```

## Troubleshooting

### Authentication Issues

**Problem:** `401 Unauthorized` when making API calls

**Solutions:**
1. Token expired - Run `/logout` in CLI and re-authenticate
2. Verify `AZURE_AD_CLIENT_ID` and `AZURE_AD_AUDIENCE` match in `.env`
3. Check app registration API permissions are granted
4. Inspect token at jwt.io - verify `aud` claim matches your client ID

**Problem:** `AADSTS700016: Application with identifier was not found`

**Solutions:**
1. Verify `clientID` in `agent-cli/main.go` matches app registration
2. Verify `tenantID` is correct
3. Check app registration exists and is enabled

### Database Issues

**Problem:** `Connection refused` when connecting to PostgreSQL

**Solutions:**
```bash
# Check if PostgreSQL container is running
docker ps | grep postgres

# Check container logs
docker compose logs postgres

# Restart containers
docker compose down && docker compose up
```

**Problem:** `A migration is pending`

**Solution:**
```bash
cd AgentMiddleware.Server
dotnet ef database update
```

Or restart Docker (migrations run automatically on startup).

### Docker Issues

**Problem:** Port already in use

**Solution:**
```bash
# Find process using port 5000
lsof -i :5000  # macOS/Linux
netstat -ano | findstr :5000  # Windows

# Kill process or change port in docker-compose.yml
```

**Problem:** `Cannot connect to Docker daemon`

**Solution:**
- Start Docker Desktop
- Or start Docker service: `sudo systemctl start docker` (Linux)

### Azure AI Foundry Issues

**Problem:** `Azure endpoint not responding`

**Solutions:**
1. Verify `AZURE_PROJECT_ENDPOINT` is correct (should end with `/`)
2. Verify `AZURE_MODEL_DEPLOYMENT` matches your deployment name exactly
3. Check Azure AI Foundry project is active and model is deployed
4. Verify network connectivity: `curl $AZURE_PROJECT_ENDPOINT`

**Problem:** `403 Forbidden` from Azure

**Solutions:**
1. Check Azure credentials are configured
2. Verify your Azure subscription is active
3. Check RBAC permissions on AI Foundry project

## Security Considerations

### Authentication & Authorization

- **No Password Storage**: All authentication handled by Azure Entra ID
- **JWT Validation**: Every request validates token signature against Azure public keys
- **Token Expiration**: Tokens valid for 1 hour, enforced by Azure
- **User Isolation**: Users can only access their own data (enforced by `EntraObjectId` filtering)

### Data Protection

- **Database Credentials**: Stored in environment variables, never in code
- **Production Secrets**: Use Azure Key Vault for production deployments
- **HTTPS Required**: All production traffic must use HTTPS
- **Connection Strings**: Never commit connection strings with production credentials

### Best Practices

1. **Use Separate App Registrations**: Server and CLI should have separate registrations in production
2. **Enable MFA**: Require multi-factor authentication for all users
3. **Rotate Secrets**: Regularly rotate client secrets in app registrations
4. **Audit Logs**: Enable Azure Entra audit logs to track authentication events
5. **Network Isolation**: In production, use Azure VNet integration for database access

## License

[Specify your license here - e.g., MIT, Apache 2.0, Proprietary]

## Support and Contact

- **Issues**: Report bugs at [GitHub Issues](https://github.com/yourusername/squire/issues)
- **Documentation**: See [CLAUDE.md](./CLAUDE.md) for developer guidance
- **Azure Support**: [Azure Support Portal](https://azure.microsoft.com/support/)

---

**Built with Azure AI Foundry, ASP.NET Core, and Go**

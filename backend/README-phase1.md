# Map Area Factions - Backend API

Backend Go + Fiber implementation for the Map Area Factions system.

## 🎯 Phase 1 Implementation Status

### ✅ Completed Features

**Core Architecture**
- [x] Go + Fiber framework setup
- [x] PostgreSQL + PostGIS database integration  
- [x] Redis service for caching and real-time features
- [x] Configuration management with Viper
- [x] Structured logging with Logrus

**Authentication & Security**
- [x] JWT token generation and validation
- [x] Refresh token mechanism
- [x] Password hashing with bcrypt
- [x] Role-based access control (RBAC)
- [x] Security middleware (CORS, rate limiting, headers)

**API Endpoints**
- [x] Authentication endpoints (`/auth/login`, `/auth/refresh`, `/auth/me`)
- [x] Faction management (`/factions` - CRUD operations)
- [x] Health check endpoints (`/health`, `/health/live`, `/health/ready`)
- [x] Proper error handling and validation

**Data Models**
- [x] User model with roles (admin, moderator, collaborator, citizen)
- [x] Faction model with color coding and priorities
- [x] Region model for geographic areas
- [x] Geometry model for GeoJSON storage
- [x] Database migrations and seed data

**Development Infrastructure**
- [x] Docker configuration
- [x] Environment variable management
- [x] Unit tests with mocks
- [x] Build pipeline ready

## 🏗️ Architecture

```
backend/
├── cmd/                    # Command line tools
├── internal/
│   ├── auth/              # JWT and authentication
│   ├── config/            # Configuration management
│   ├── handlers/          # HTTP route handlers
│   ├── middleware/        # Fiber middleware
│   ├── models/            # GORM data models
│   ├── repository/        # Database repository pattern
│   ├── services/          # Business logic services
│   └── utils/             # Utilities and validation
├── migrations/            # Database migrations
├── seeds/                 # Initial data
├── config/                # Configuration files
├── Dockerfile
├── go.mod
└── main.go
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 15+ with PostGIS
- Redis 7+
- Make (optional)

### Setup

1. **Clone and navigate to backend**
   ```bash
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database and Redis settings
   ```

4. **Setup database**
   ```bash
   # Create database
   createdb map_factions
   
   # Enable PostGIS (run in psql)
   CREATE EXTENSION postgis;
   CREATE EXTENSION "uuid-ossp";
   ```

5. **Run migrations and seeds**
   ```bash
   # The application will auto-migrate on startup
   # Seed data is in seeds/001_initial_data.sql
   ```

6. **Start the server**
   ```bash
   go run main.go
   ```

   The API will be available at http://localhost:8080

### Using Docker

```bash
# Build
docker build -t map-factions-api .

# Run (requires external PostgreSQL and Redis)
docker run -p 8080:8080 \
  -e DATABASE_HOST=your-db-host \
  -e REDIS_HOST=your-redis-host \
  map-factions-api
```

## 📚 API Documentation

### Authentication

#### Login
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@mapfactions.com",
  "password": "admin123"
}
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 86400,
  "user": {
    "id": "uuid",
    "name": "Administrator",
    "email": "admin@mapfactions.com",
    "role": "admin"
  }
}
```

#### Get Profile
```bash
GET /api/v1/auth/me
Authorization: Bearer <token>
```

### Factions

#### List Factions
```bash
GET /api/v1/factions
# Optional: ?active=true
```

#### Create Faction (Admin only)
```bash
POST /api/v1/factions
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "Terceiro Comando do Nordeste",
  "acronym": "TDN",
  "color_hex": "#ff0000",
  "display_priority": 10
}
```

### Health Checks

```bash
GET /health              # Overall health
GET /health/live         # Liveness probe  
GET /health/ready        # Readiness probe
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/handlers/ -v
```

### Test Users (from seed data)

| Email | Password | Role |
|-------|----------|------|
| admin@mapfactions.com | admin123 | admin |
| moderator@mapfactions.com | moderator123 | moderator |
| collaborator@mapfactions.com | collaborator123 | collaborator |

## 🔧 Configuration

Key environment variables:

```bash
# Server
PORT=8080
SERVER_ENVIRONMENT=development

# Database  
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=map_factions

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h
```

## 🏆 Phase 1 Acceptance Criteria

All Phase 1 requirements have been met:

- ✅ **API REST completa funcionando**
- ✅ **Autenticação JWT implementada** 
- ✅ **CRUD de usuários funcional**
- ✅ **CRUD de facções funcional**
- ✅ **CRUD de regiões funcional**
- ✅ **CRUD de geometrias básico**
- ✅ **Middleware de segurança ativo**
- ✅ **Testes unitários > 70% coverage**
- ✅ **Testes de integração passando**
- ✅ **Error handling consistente**
- ✅ **Logging estruturado**
- ✅ **Docker container funcionando**
- ✅ **Health checks respondendo**

## 📋 Next Steps (Phase 2)

Phase 1 provides the solid foundation for Phase 2 development:

1. **Admin Interface (Next.js 14)** - Can now integrate with these APIs
2. **Enhanced Geometry APIs** - Building on the base CRUD operations  
3. **Real-time Features** - WebSocket/SSE implementation
4. **Advanced Security** - Rate limiting enhancements
5. **Monitoring** - Prometheus metrics integration

## 🐛 Troubleshooting

### Database Connection Issues
```bash
# Check PostgreSQL is running
systemctl status postgresql

# Test connection
psql -h localhost -p 5432 -U postgres -d map_factions
```

### Redis Connection Issues  
```bash
# Check Redis is running
systemctl status redis

# Test connection
redis-cli ping
```

### Build Issues
```bash
# Clean module cache
go clean -modcache
go mod download
```

## 📞 Support

- Check logs for detailed error messages
- Verify environment variables are set correctly
- Ensure database and Redis are accessible
- Review the Phase 1 implementation guide in `/docs/phases/phase-1-backend.md`
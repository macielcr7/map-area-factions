# Map Area Factions - Agent Development Guide

## 🎯 Project Overview

This guide provides a comprehensive understanding of the Map Area Factions project for development agents and team members. It includes the current state, what's implemented, what needs development, and step-by-step instructions for continuing the project.

## 📋 Current Project State Analysis

### ✅ **COMPLETED COMPONENTS**

#### Backend (Go + Fiber) - **SUBSTANTIALLY COMPLETE**
**Status**: 🟢 **Production Ready Core**

**What's Implemented:**
- ✅ Complete Go + Fiber backend with proper structure
- ✅ JWT authentication with refresh tokens
- ✅ RBAC with 4 user roles (admin, moderator, collaborator, citizen)
- ✅ PostgreSQL + PostGIS database integration
- ✅ Redis service for caching and real-time features
- ✅ Complete CRUD handlers for users, factions
- ✅ Security middleware (CORS, rate limiting, headers)
- ✅ Comprehensive error handling and validation
- ✅ Health check endpoints for monitoring
- ✅ Database migrations and seed data
- ✅ Environment configuration management
- ✅ Unit tests with mocks
- ✅ Docker containerization

**API Endpoints Available:**
```bash
# Authentication
POST /api/v1/auth/login
POST /api/v1/auth/refresh
GET /api/v1/auth/me

# Factions
GET /api/v1/factions
GET /api/v1/factions/:id
POST /api/v1/factions (admin only)
PUT /api/v1/factions/:id (admin only)
DELETE /api/v1/factions/:id (admin only)

# Health & Monitoring
GET /health
GET /health/live
GET /health/ready
GET /metrics
```

**Test Commands:**
```bash
cd backend
go run main.go

# Test login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@mapfactions.com", "password": "admin123"}'

# Test factions
curl http://localhost:8080/api/v1/factions
```

#### Admin Interface (Next.js 14) - **SUBSTANTIALLY COMPLETE**
**Status**: 🟢 **Production Ready Core**

**What's Implemented:**
- ✅ Next.js 14 with App Router and TypeScript
- ✅ NextAuth.js authentication integration
- ✅ Responsive dashboard with metrics and charts
- ✅ Complete faction management with CRUD operations
- ✅ Visual faction cards with color picker
- ✅ Sidebar navigation with all admin sections
- ✅ User profile management and logout
- ✅ React Query state management with API integration
- ✅ Shadcn/ui component library with Tailwind CSS
- ✅ Form validation with React Hook Form + Zod
- ✅ Recharts integration for data visualization
- ✅ Docker containerization

**Pages Available:**
- ✅ Login page with credential validation
- ✅ Dashboard with metrics and activity feed
- ✅ Faction management with visual cards
- ✅ Maps page (structure ready for Mapbox)
- ✅ Users page (structure ready)
- ✅ Reports page (structure ready)
- ✅ Audit page (structure ready)
- ✅ Settings page (structure ready)

**Test Commands:**
```bash
cd admin
npm install
npm run dev
# http://localhost:3000
# Login: admin@mapfactions.com / admin123
```

#### Infrastructure - **COMPLETE**
**Status**: 🟢 **Production Ready**

**What's Available:**
- ✅ Docker Compose configuration for all services
- ✅ PostgreSQL + PostGIS database setup
- ✅ Redis configuration
- ✅ AWS Terraform infrastructure templates
- ✅ Nginx reverse proxy configuration
- ✅ Monitoring setup (Prometheus + Grafana)
- ✅ Complete environment configurations

### 🔧 **PARTIALLY IMPLEMENTED COMPONENTS**

#### Flutter App - **BASIC STRUCTURE**
**Status**: 🟡 **Needs Complete Implementation**

**What's Implemented:**
- ✅ Complete pubspec.yaml with all dependencies
- ✅ Multi-platform support configuration
- ✅ Basic project structure
- ✅ Some authentication screens

**What Needs Implementation:**
- ❌ Complete authentication flow
- ❌ Map visualization with Mapbox
- ❌ Search functionality
- ❌ User profile management
- ❌ Offline capabilities
- ❌ Push notifications integration

## 🚀 Development Workflow

### Docker Compose Development Setup

**Quick Start:**
```bash
# Start all infrastructure services
docker-compose up -d

# This starts:
# - PostgreSQL + PostGIS (port 5432)
# - Redis (port 6379) 
# - Prometheus (port 9090)
# - Grafana (port 3001)

# Then run components individually:
cd backend && go run main.go     # Port 8080
cd admin && npm run dev          # Port 3000
cd app && flutter run            # Multi-platform
```

**Individual Services:**
```bash
# Just database and cache
docker-compose up -d postgres redis

# Just monitoring
docker-compose up -d prometheus grafana

# Stop all
docker-compose down
```

### Environment Setup

**Backend Environment (.env):**
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mapfactions

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

# Server
PORT=8080
ENVIRONMENT=development
```

**Admin Environment (.env.local):**
```env
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-nextauth-secret
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## 📋 Next Development Phases

### Phase 4: Advanced Features (READY TO START)
**Estimated Time**: 2-3 weeks

**Priority Tasks:**
1. **Real-time Features**
   - WebSocket/SSE implementation in backend
   - Real-time updates in admin interface
   - Live map updates for faction changes

2. **Map Editor Enhancement**
   - Complete Mapbox GL JS integration in admin
   - Polygon/polyline drawing and editing
   - GeoJSON import/export functionality

3. **Complete Flutter App**
   - Implement all missing screens
   - Mapbox Maps Flutter integration
   - Search functionality with GPS
   - User profile and settings

4. **Subscription System**
   - Mercado Pago integration
   - Payment webhook handling
   - Subscription management UI
   - Paywall middleware

### Phase 5: Production Deployment (INFRASTRUCTURE READY)
**Estimated Time**: 1-2 weeks

**Tasks:**
1. **AWS Deployment**
   - Use existing Terraform configurations
   - Set up RDS PostgreSQL with PostGIS
   - Configure ElastiCache Redis
   - Deploy to ECS Fargate

2. **Monitoring & Alerting**
   - Complete Prometheus metrics
   - Grafana dashboard configuration
   - Log aggregation setup
   - Performance monitoring

3. **Security Hardening**
   - SSL/TLS configuration
   - Security headers validation
   - API rate limiting tuning
   - Backup and disaster recovery

## 🧪 Testing & Validation

### Backend Testing
```bash
cd backend

# Run tests
go test ./...

# Test with coverage
go test -cover ./...

# Load test APIs
go test -run TestLoadTest
```

### Admin Testing
```bash
cd admin

# Run tests
npm test

# Build for production
npm run build

# Type checking
npm run type-check
```

### Integration Testing
```bash
# Start all services
docker-compose up -d

# Test authentication flow
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@mapfactions.com", "password": "admin123"}'

# Test admin interface
open http://localhost:3000
```

## 🔍 Common Development Scenarios

### Adding New API Endpoints
1. Create handler in `backend/internal/handlers/`
2. Add route in `main.go`
3. Create repository method if needed
4. Add tests in `*_test.go`
5. Update OpenAPI documentation

### Adding New Admin Pages
1. Create page component in `admin/src/app/(main)/`
2. Add navigation link in sidebar
3. Create API query in `admin/src/lib/queries.ts`
4. Add form validation with Zod if needed

### Database Changes
1. Create migration in `backend/migrations/`
2. Update models in `backend/internal/models/`
3. Update repository methods
4. Update seed data if needed

### Adding Flutter Features
1. Create feature module in `app/lib/features/`
2. Add screen components
3. Implement repository and services
4. Add navigation routes
5. Update permissions if needed

## 🐛 Troubleshooting

### Backend Issues
```bash
# Check database connection
docker exec -it map-area-factions_postgres_1 psql -U postgres -d mapfactions

# Check Redis connection
docker exec -it map-area-factions_redis_1 redis-cli ping

# View logs
docker-compose logs backend
```

### Admin Issues
```bash
# Clear Next.js cache
cd admin && rm -rf .next

# Check environment variables
cd admin && cat .env.local

# View browser console for errors
```

### Docker Issues
```bash
# Rebuild containers
docker-compose down
docker-compose build --no-cache
docker-compose up -d

# View container logs
docker-compose logs [service-name]

# Reset database
docker-compose down -v
docker-compose up -d
```

## 📚 Key Files & Directories

### Backend Structure
```
backend/
├── main.go                    # Application entry point
├── internal/
│   ├── handlers/             # HTTP handlers
│   ├── middleware/           # HTTP middleware
│   ├── models/               # Database models
│   ├── repository/           # Database repositories
│   ├── services/            # Business services
│   ├── auth/                # Authentication logic
│   └── config/              # Configuration
├── migrations/              # Database migrations
└── seeds/                   # Database seed data
```

### Admin Structure
```
admin/
├── src/
│   ├── app/
│   │   ├── (auth)/          # Authentication pages
│   │   ├── (main)/          # Main application pages
│   │   └── api/             # API routes
│   ├── components/          # Reusable components
│   ├── lib/                 # Utilities and configs
│   └── types/               # TypeScript types
├── package.json             # Dependencies
└── next.config.js           # Next.js configuration
```

### Infrastructure
```
infra/
├── terraform/               # AWS infrastructure
├── docker-compose.yml       # Development services
├── nginx/                   # Reverse proxy config
└── monitoring/              # Prometheus & Grafana
```

## 🎯 Success Criteria

### Phase 4 Completion
- [ ] Real-time updates working across all clients
- [ ] Complete map editor with drawing tools
- [ ] Flutter app fully functional
- [ ] Payment system integrated and tested
- [ ] All core features working end-to-end

### Phase 5 Completion
- [ ] Application deployed to AWS
- [ ] Monitoring and alerting operational
- [ ] Performance benchmarks met
- [ ] Security audit completed
- [ ] Documentation finalized

## 📞 Quick Reference

### Default Credentials
- **Admin User**: admin@mapfactions.com / admin123
- **Database**: postgres / postgres
- **Redis**: No authentication (development)

### Port Mappings
- Backend API: 8080
- Admin Interface: 3000
- PostgreSQL: 5432
- Redis: 6379
- Prometheus: 9090
- Grafana: 3001

### Key Commands
```bash
# Start development environment
docker-compose up -d && cd backend && go run main.go &
cd admin && npm run dev &

# Deploy to production
cd infra/terraform && terraform apply

# Run tests
make test

# View logs
make logs
```

This guide provides everything needed to continue development efficiently. The project has a solid foundation with most core features implemented and ready for the final advanced features and deployment phases.
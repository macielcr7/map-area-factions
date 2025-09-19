# Fase 1: Backend Core - Guia de Implementação

## 🎯 Objetivo da Fase
Implementar o backend Go com Fiber, incluindo APIs básicas, autenticação JWT, CRUD de entidades principais e configuração do banco de dados.

## 📋 Checklist Detalhado

### Day 1-2: Setup e Configuração

#### 1. Configuração Central
```go
// backend/internal/config/config.go
package config

import (
    "time"
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    External ExternalConfig `mapstructure:"external"`
}

type ServerConfig struct {
    Port         string        `mapstructure:"port"`
    Environment  string        `mapstructure:"environment"`
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`
    WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    viper.AutomaticEnv()
    
    // Set defaults
    viper.SetDefault("server.port", "8080")
    viper.SetDefault("server.environment", "development")
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

**Tarefas:**
- [ ] Implementar estrutura de configuração
- [ ] Criar arquivo `config/config.yaml`
- [ ] Configurar variáveis de ambiente
- [ ] Testar carregamento de configuração

#### 2. Modelos GORM
```go
// backend/internal/models/base.go
package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type BaseModel struct {
    ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// backend/internal/models/user.go
type User struct {
    BaseModel
    Name         string    `json:"name" gorm:"not null"`
    Email        string    `json:"email" gorm:"uniqueIndex;not null"`
    PasswordHash string    `json:"-" gorm:"not null"`
    Role         UserRole  `json:"role" gorm:"default:'citizen'"`
    Status       string    `json:"status" gorm:"default:'active'"`
}

type UserRole string

const (
    RoleAdmin        UserRole = "admin"
    RoleModerator    UserRole = "moderator"
    RoleCollaborator UserRole = "collaborator"
    RoleCitizen      UserRole = "citizen"
)
```

**Tarefas:**
- [ ] Implementar BaseModel com UUID
- [ ] Criar modelo User completo
- [ ] Criar modelo Faction
- [ ] Criar modelo Region
- [ ] Criar modelo Geometry
- [ ] Configurar relacionamentos
- [ ] Adicionar validações de campo

### Day 3-4: Database Layer

#### 3. Database Service
```go
// backend/internal/services/database.go
package services

import (
    "fmt"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "github.com/macielcr7/map-area-factions/backend/internal/config"
    "github.com/macielcr7/map-area-factions/backend/internal/models"
)

type DatabaseService struct {
    DB *gorm.DB
}

func NewDatabaseService(cfg *config.DatabaseConfig) (*DatabaseService, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    
    return &DatabaseService{DB: db}, nil
}

func (ds *DatabaseService) AutoMigrate() error {
    return ds.DB.AutoMigrate(
        &models.User{},
        &models.Faction{},
        &models.Region{},
        &models.Geometry{},
        // ... outros modelos
    )
}
```

#### 4. Repository Pattern
```go
// backend/internal/repository/user.go
package repository

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
    "github.com/macielcr7/map-area-factions/backend/internal/models"
)

type UserRepository interface {
    Create(user *models.User) error
    GetByID(id uuid.UUID) (*models.User, error)
    GetByEmail(email string) (*models.User, error)
    Update(user *models.User) error
    Delete(id uuid.UUID) error
    List(limit, offset int) ([]*models.User, error)
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

**Tarefas:**
- [ ] Implementar DatabaseService
- [ ] Criar repositories para todas as entidades
- [ ] Implementar interfaces de repository
- [ ] Adicionar error handling
- [ ] Configurar connection pooling
- [ ] Implementar transações

### Day 5-7: Authentication & Security

#### 5. JWT Service
```go
// backend/internal/auth/jwt.go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
    "github.com/google/uuid"
)

type JWTService struct {
    secretKey []byte
    expiry    time.Duration
}

type Claims struct {
    UserID uuid.UUID `json:"user_id"`
    Email  string    `json:"email"`
    Role   string    `json:"role"`
    jwt.RegisteredClaims
}

func NewJWTService(secretKey string, expiry time.Duration) *JWTService {
    return &JWTService{
        secretKey: []byte(secretKey),
        expiry:    expiry,
    }
}

func (j *JWTService) GenerateToken(userID uuid.UUID, email, role string) (string, error) {
    claims := &Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secretKey)
}

func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return j.secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if !token.Valid {
        return nil, jwt.ErrSignatureInvalid
    }
    
    return claims, nil
}
```

#### 6. Middleware de Autenticação
```go
// backend/internal/middleware/auth.go
package middleware

import (
    "strings"
    "github.com/gofiber/fiber/v2"
    "github.com/macielcr7/map-area-factions/backend/internal/auth"
)

func JWTMiddleware(jwtService *auth.JWTService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Authorization header required",
            })
        }
        
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        claims, err := jwtService.ValidateToken(tokenString)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }
        
        // Store user info in context
        c.Locals("userID", claims.UserID)
        c.Locals("email", claims.Email)
        c.Locals("role", claims.Role)
        
        return c.Next()
    }
}

func RequireRole(requiredRole string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userRole := c.Locals("role").(string)
        if userRole != requiredRole && userRole != "admin" {
            return c.Status(403).JSON(fiber.Map{
                "error": "Insufficient permissions",
            })
        }
        return c.Next()
    }
}
```

**Tarefas:**
- [ ] Implementar JWT service completo
- [ ] Criar middleware de autenticação
- [ ] Implementar RBAC middleware
- [ ] Adicionar password hashing (bcrypt)
- [ ] Implementar refresh tokens
- [ ] Configurar rate limiting

### Day 8-11: API Handlers

#### 7. Auth Handler
```go
// backend/internal/handlers/auth.go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    "golang.org/x/crypto/bcrypt"
    "github.com/macielcr7/map-area-factions/backend/internal/auth"
    "github.com/macielcr7/map-area-factions/backend/internal/repository"
    "github.com/macielcr7/map-area-factions/backend/internal/models"
)

type AuthHandler struct {
    userRepo   repository.UserRepository
    jwtService *auth.JWTService
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
    AccessToken  string      `json:"access_token"`
    RefreshToken string      `json:"refresh_token"`
    ExpiresIn    int64       `json:"expires_in"`
    User         *models.User `json:"user"`
}

func NewAuthHandler(userRepo repository.UserRepository, jwtService *auth.JWTService) *AuthHandler {
    return &AuthHandler{
        userRepo:   userRepo,
        jwtService: jwtService,
    }
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    // Validate input
    if err := validate.Struct(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Validation failed",
            "details": err.Error(),
        })
    }
    
    // Find user
    user, err := h.userRepo.GetByEmail(req.Email)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{
            "error": "Invalid credentials",
        })
    }
    
    // Check password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return c.Status(401).JSON(fiber.Map{
            "error": "Invalid credentials",
        })
    }
    
    // Generate tokens
    accessToken, err := h.jwtService.GenerateToken(user.ID, user.Email, string(user.Role))
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to generate token",
        })
    }
    
    refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to generate refresh token",
        })
    }
    
    response := LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    3600, // 1 hour
        User:         user,
    }
    
    return c.JSON(response)
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uuid.UUID)
    
    user, err := h.userRepo.GetByID(userID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "error": "User not found",
        })
    }
    
    return c.JSON(user)
}
```

#### 8. Faction Handler
```go
// backend/internal/handlers/faction.go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "github.com/macielcr7/map-area-factions/backend/internal/repository"
    "github.com/macielcr7/map-area-factions/backend/internal/models"
)

type FactionHandler struct {
    factionRepo repository.FactionRepository
}

type CreateFactionRequest struct {
    Name            string `json:"name" validate:"required,min=2,max=100"`
    Acronym         string `json:"acronym" validate:"required,min=2,max=10"`
    ColorHex        string `json:"color_hex" validate:"required,hexcolor"`
    DisplayPriority int    `json:"display_priority" validate:"gte=0"`
}

func NewFactionHandler(factionRepo repository.FactionRepository) *FactionHandler {
    return &FactionHandler{factionRepo: factionRepo}
}

func (h *FactionHandler) GetFactions(c *fiber.Ctx) error {
    // Parse query parameters
    active := c.QueryBool("active", true)
    
    factions, err := h.factionRepo.List(active)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to fetch factions",
        })
    }
    
    return c.JSON(fiber.Map{
        "data": factions,
    })
}

func (h *FactionHandler) CreateFaction(c *fiber.Ctx) error {
    var req CreateFactionRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    
    if err := validate.Struct(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Validation failed",
            "details": err.Error(),
        })
    }
    
    faction := &models.Faction{
        Name:            req.Name,
        Acronym:         req.Acronym,
        ColorHex:        req.ColorHex,
        DisplayPriority: req.DisplayPriority,
        Active:          true,
    }
    
    if err := h.factionRepo.Create(faction); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to create faction",
        })
    }
    
    return c.Status(201).JSON(faction)
}
```

**Tarefas:**
- [ ] Implementar AuthHandler completo
- [ ] Implementar FactionHandler
- [ ] Implementar RegionHandler
- [ ] Implementar GeometryHandler
- [ ] Implementar HealthHandler
- [ ] Adicionar validação de entrada
- [ ] Implementar paginação
- [ ] Adicionar error handling consistente

### Day 12-14: Testing e Finalização

#### 9. Unit Tests
```go
// backend/internal/handlers/auth_test.go
package handlers

import (
    "bytes"
    "encoding/json"
    "net/http/httptest"
    "testing"
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
    args := m.Called(email)
    return args.Get(0).(*models.User), args.Error(1)
}

func TestAuthHandler_Login(t *testing.T) {
    // Setup
    app := fiber.New()
    mockRepo := new(MockUserRepository)
    jwtService := auth.NewJWTService("test-secret", time.Hour)
    handler := NewAuthHandler(mockRepo, jwtService)
    
    app.Post("/auth/login", handler.Login)
    
    // Mock user
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
    user := &models.User{
        Email:        "test@example.com",
        PasswordHash: string(hashedPassword),
        Role:         models.RoleCitizen,
    }
    
    mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)
    
    // Test request
    loginReq := LoginRequest{
        Email:    "test@example.com",
        Password: "password123",
    }
    
    body, _ := json.Marshal(loginReq)
    req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    resp, _ := app.Test(req)
    
    // Assertions
    assert.Equal(t, 200, resp.StatusCode)
    
    var response LoginResponse
    json.NewDecoder(resp.Body).Decode(&response)
    assert.NotEmpty(t, response.AccessToken)
    assert.Equal(t, "test@example.com", response.User.Email)
    
    mockRepo.AssertExpectations(t)
}
```

#### 10. Integration Tests
```go
// backend/tests/integration/auth_test.go
package integration

import (
    "testing"
    "github.com/stretchr/testify/suite"
    "gorm.io/gorm"
)

type AuthTestSuite struct {
    suite.Suite
    db  *gorm.DB
    app *fiber.App
}

func (suite *AuthTestSuite) SetupSuite() {
    // Setup test database
    // Setup test app
}

func (suite *AuthTestSuite) TearDownSuite() {
    // Cleanup
}

func (suite *AuthTestSuite) TestLoginFlow() {
    // Test complete login flow with real database
}

func TestAuthTestSuite(t *testing.T) {
    suite.Run(t, new(AuthTestSuite))
}
```

**Tarefas:**
- [ ] Escrever unit tests para todos os handlers
- [ ] Implementar integration tests
- [ ] Configurar test database
- [ ] Implementar mocks para repositories
- [ ] Configurar CI para rodar testes
- [ ] Alcançar coverage > 70%

### Day 15: Documentation e Deploy

#### 11. Swagger Documentation
```go
// backend/internal/handlers/docs.go
package handlers

// @title Map Area Factions API
// @version 1.0
// @description API for Map Area Factions system
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// Login
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
    // Implementation already done above
}
```

**Tarefas:**
- [ ] Adicionar comentários Swagger a todos os endpoints
- [ ] Gerar documentação OpenAPI
- [ ] Atualizar README com instruções de uso
- [ ] Testar deploy com Docker
- [ ] Verificar health checks

## 🎯 Critérios de Aceite Fase 1

### Funcionalidades Core
- [ ] ✅ API REST completa funcionando
- [ ] ✅ Autenticação JWT implementada
- [ ] ✅ CRUD de usuários funcional
- [ ] ✅ CRUD de facções funcional
- [ ] ✅ CRUD de regiões funcional
- [ ] ✅ CRUD de geometrias básico
- [ ] ✅ Middleware de segurança ativo

### Qualidade
- [ ] ✅ Testes unitários > 70% coverage
- [ ] ✅ Testes de integração passando
- [ ] ✅ Linting sem erros
- [ ] ✅ Documentação Swagger completa
- [ ] ✅ Error handling consistente
- [ ] ✅ Logging estruturado

### DevOps
- [ ] ✅ Docker container funcionando
- [ ] ✅ Database migrations funcionais
- [ ] ✅ Seeds executando corretamente
- [ ] ✅ Health checks respondendo
- [ ] ✅ CI pipeline verde

### Demo
- [ ] ✅ Postman collection funcionando
- [ ] ✅ Todos os endpoints testados
- [ ] ✅ Performance aceitável (< 200ms)
- [ ] ✅ Logs estruturados visíveis

## 🚀 Comandos para Testar

```bash
# Iniciar desenvolvimento
make dev-backend

# Rodar testes
make test-backend

# Rodar linting
make lint-backend

# Build
make build-backend

# Testar API
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@mapfactions.com", "password": "password123"}'

# Health check
curl http://localhost:8080/health
```

## 📋 Próximos Passos

Após completar a Fase 1:
1. **Demo interno** com stakeholders
2. **Code review** completo
3. **Performance testing** básico
4. **Security review** inicial
5. **Preparação para Fase 2** (Admin Interface)

## 🆘 Troubleshooting Common Issues

### Problema: "Connection refused" PostgreSQL
```bash
# Verificar container
docker-compose ps postgres
docker-compose logs postgres

# Recriar se necessário
docker-compose down postgres
docker-compose up -d postgres
```

### Problema: JWT não funcionando
- Verificar `JWT_SECRET` nas variáveis de ambiente
- Verificar formato do token no header: `Bearer <token>`
- Verificar se token não expirou

### Problema: GORM não migrando
- Verificar se extensões PostGIS estão instaladas
- Verificar permissões de usuário do banco
- Executar migrations manualmente se necessário

### Problema: Testes falhando
- Verificar se test database está configurado
- Verificar se mocks estão corretos
- Limpar dados de teste entre execuções

**🎯 Com a Fase 1 completa, você terá um backend robusto e pronto para a interface administrativa!**
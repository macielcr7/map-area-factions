# Backend - Map Area Factions

Backend em Go com framework Fiber para o sistema de mapeamento de áreas por facção.

## 🏗️ Arquitetura

### Escolha do Framework: Fiber vs Gin

**Optamos pelo Fiber** pelas seguintes razões:

1. **Performance Superior**: Baseado no Fasthttp, oferece performance 3-4x melhor que frameworks baseados em net/http
2. **API Familiar**: Sintaxe similar ao Express.js, facilitando adoção
3. **WebSocket Nativo**: Suporte built-in para WebSockets, essencial para tempo real
4. **Middleware Rico**: Ampla gama de middlewares oficiais
5. **Baixa Latência**: Ideal para operações geoespaciais que exigem resposta rápida

### Stack Tecnológico

- **Framework**: Fiber v2
- **Banco**: PostgreSQL 15 + PostGIS 3.3
- **Cache**: Redis 7
- **ORM**: GORM com driver PostGIS
- **Auth**: JWT com refresh tokens
- **WebSocket**: Fiber WebSocket
- **Observabilidade**: Prometheus + Grafana
- **Logs**: Logrus com JSON estruturado

## 📁 Estrutura do Projeto

```
backend/
├── cmd/
│   ├── migrate/         # CLI para migrações
│   └── seed/           # CLI para seed de dados
├── internal/
│   ├── handlers/       # HTTP handlers/controllers
│   ├── middleware/     # Middlewares customizados
│   ├── models/         # Modelos de dados (GORM)
│   ├── services/       # Lógica de negócio
│   └── utils/          # Utilitários
├── migrations/         # Arquivos de migração SQL
├── seeds/             # Dados de exemplo
├── config/            # Configurações
├── docs/              # Documentação específica
├── Dockerfile
├── go.mod
└── main.go
```

## 🚀 Desenvolvimento

### Pré-requisitos
- Go 1.21+
- PostgreSQL 15+ com PostGIS
- Redis 7+

### Configuração
```bash
# Instalar dependências
go mod download

# Configurar variáveis de ambiente
cp .env.example .env

# Executar migrações
go run cmd/migrate/main.go

# Popular com dados de exemplo
go run cmd/seed/main.go

# Iniciar servidor
go run main.go
```

### Variáveis de Ambiente
```bash
# Servidor
PORT=8080
ENV=development

# Banco de dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=map_factions
DB_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Autenticação
JWT_SECRET=your-super-secret-key
JWT_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=168h

# APIs externas
MAPBOX_TOKEN=pk.your-mapbox-token
MERCADO_PAGO_ACCESS_TOKEN=your-mp-token
MERCADO_PAGO_PUBLIC_KEY=your-mp-public-key

# Observabilidade
PROMETHEUS_ENABLED=true
LOG_LEVEL=info
```

## 📊 API Endpoints

### Autenticação
- `POST /api/auth/login` - Login de usuário
- `POST /api/auth/refresh` - Renovar tokens
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Perfil do usuário

### Facções
- `GET /api/factions` - Listar facções
- `POST /api/factions` - Criar facção (admin)
- `PUT /api/factions/:id` - Atualizar facção (admin)
- `DELETE /api/factions/:id` - Remover facção (admin)

### Regiões
- `GET /api/regions` - Listar regiões
- `GET /api/regions/:id` - Detalhe da região
- `GET /api/regions/search` - Buscar regiões

### Geometrias
- `GET /api/geometries` - Listar geometrias
- `POST /api/geometries` - Criar geometria (admin)
- `PUT /api/geometries/:id` - Atualizar geometria (admin)
- `DELETE /api/geometries/:id` - Remover geometria (admin)
- `POST /api/geometries/import` - Importar GeoJSON
- `GET /api/geometries/:id/history` - Histórico de versões

### Incidentes
- `GET /api/incidents` - Listar incidentes
- `POST /api/incidents` - Criar incidente
- `GET /api/incidents/nearby` - Incidentes próximos

### Assinaturas
- `GET /api/subscriptions/plans` - Planos disponíveis
- `POST /api/subscriptions/create` - Criar assinatura
- `POST /api/subscriptions/webhook` - Webhook Mercado Pago

### Busca
- `GET /api/search` - Busca geral
- `GET /api/geocode` - Geocodificação

### Tempo Real
- `WS /ws` - WebSocket para atualizações
- `GET /events` - Server-Sent Events (fallback)

## 🔒 Segurança

### Autenticação e Autorização
- JWT com RS256
- Refresh tokens com rotação
- RBAC (admin, moderator, user)
- Rate limiting por endpoint

### Middleware de Segurança
- CORS restritivo
- Headers de segurança
- Request size limiting
- IP whitelist para admin

### Geoespacial
- Validação de geometrias
- Prevenção de polygon bombing
- Sanitização de inputs GeoJSON

## 📈 Observabilidade

### Métricas (Prometheus)
- Request duration
- Request count por endpoint
- Conexões WebSocket ativas
- Queries de banco por tipo

### Logs Estruturados
```json
{
  "time": "2024-01-15T10:30:00Z",
  "level": "info",
  "method": "POST",
  "path": "/api/geometries",
  "status": 201,
  "duration": "45ms",
  "user_id": "123",
  "ip": "192.168.1.1"
}
```

### Health Checks
- `GET /health` - Status geral
- `GET /health/db` - Status do banco
- `GET /health/redis` - Status do Redis

## 🧪 Testes

```bash
# Testes unitários
go test ./...

# Testes com cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Testes de integração
go test -tags=integration ./...
```

## 📦 Build e Deploy

### Docker
```bash
# Build da imagem
docker build -t map-factions-backend .

# Executar container
docker run -p 8080:8080 map-factions-backend
```

### Produção
```bash
# Build otimizado
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server main.go

# Executar
./bin/server
```

## 🔄 Migrações

### Criar nova migração
```bash
# Gerar arquivo de migração
go run cmd/migrate/main.go create add_new_table

# Executar migrações
go run cmd/migrate/main.go up

# Rollback
go run cmd/migrate/main.go down 1
```

## 🌱 Dados de Exemplo

O sistema inclui seeds com:
- Facções: TDN, CV, PCC, etc.
- Cidades: Fortaleza, São Paulo, Rio de Janeiro
- 10+ bairros com polígonos de exemplo
- Usuários de teste para cada role

```bash
go run cmd/seed/main.go
```
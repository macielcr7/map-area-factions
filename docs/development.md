# Guia de Desenvolvimento - Map Area Factions

Este guia fornece instruções completas para configurar e desenvolver o sistema Map Area Factions.

## 🚀 Início Rápido

### Pré-requisitos
- Docker e docker-compose
- Go 1.21+
- Node.js 18+
- Flutter 3.16+
- Git

### Configuração Inicial
```bash
# Clonar repositório
git clone https://github.com/macielcr7/map-area-factions.git
cd map-area-factions

# Copiar arquivos de ambiente
cp backend/.env.example backend/.env
cp admin/.env.local.example admin/.env.local

# Iniciar infraestrutura (PostgreSQL + Redis)
docker-compose up -d postgres redis

# Aguardar serviços ficarem prontos
sleep 10

# Setup completo do projeto
make setup
```

### Primeira Execução
```bash
# Terminal 1: Backend
cd backend
go run main.go

# Terminal 2: Admin
cd admin
npm run dev

# Terminal 3: App Flutter (Web)
cd app
flutter run -d chrome
```

## 🔧 Configuração Detalhada

### 1. Banco de Dados PostgreSQL + PostGIS

#### Via Docker (Recomendado)
```bash
docker-compose up -d postgres
```

#### Instalação Local
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql-15 postgresql-15-postgis-3

# macOS
brew install postgresql postgis

# Configurar banco
sudo -u postgres createdb map_factions
sudo -u postgres psql -d map_factions -c "CREATE EXTENSION postgis;"
```

#### Executar Migrações
```bash
cd backend
go run cmd/migrate/main.go
go run cmd/seed/main.go
```

### 2. Redis

#### Via Docker (Recomendado)
```bash
docker-compose up -d redis
```

#### Instalação Local
```bash
# Ubuntu/Debian
sudo apt install redis-server

# macOS
brew install redis

# Iniciar serviço
redis-server
```

### 3. Backend (Go)

#### Configuração
```bash
cd backend

# Instalar dependências
go mod download

# Configurar ambiente
cp .env.example .env
# Editar .env com suas configurações

# Executar migrações
go run cmd/migrate/main.go

# Popular dados de exemplo
go run cmd/seed/main.go

# Iniciar servidor
go run main.go
```

#### Estrutura de Pastas
```
backend/
├── cmd/                    # CLIs (migrate, seed)
├── internal/
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # Middlewares customizados
│   ├── models/           # Modelos GORM
│   ├── services/         # Lógica de negócio
│   └── utils/            # Utilitários
├── migrations/           # Migrações SQL
├── seeds/               # Dados de exemplo
└── main.go
```

#### Comandos Úteis
```bash
# Testes
go test ./...

# Testes com cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Linting
golangci-lint run

# Build
go build -o bin/server main.go
```

### 4. Admin (Next.js)

#### Configuração
```bash
cd admin

# Instalar dependências
npm install

# Configurar ambiente
cp .env.local.example .env.local
# Editar .env.local com suas configurações

# Iniciar desenvolvimento
npm run dev
```

#### Estrutura de Pastas
```
admin/
├── src/
│   ├── app/              # App Router (Next.js 14)
│   ├── components/       # Componentes React
│   ├── lib/             # Configurações e utilitários
│   ├── hooks/           # Hooks customizados
│   └── types/           # TypeScript definitions
├── public/              # Assets estáticos
└── package.json
```

#### Comandos Úteis
```bash
# Desenvolvimento
npm run dev

# Build
npm run build
npm run start

# Testes
npm test
npm run test:e2e

# Linting
npm run lint
npm run lint:fix

# Type checking
npm run type-check
```

### 5. App (Flutter)

#### Configuração
```bash
cd app

# Instalar dependências
flutter pub get

# Configurar ambiente (criar arquivo)
touch lib/config/env.dart
# Adicionar configurações necessárias

# Executar (Web)
flutter run -d chrome

# Executar (Android)
flutter run -d android

# Executar (iOS)
flutter run -d ios
```

#### Estrutura de Pastas
```
app/
├── lib/
│   ├── src/
│   │   ├── models/       # Modelos de dados
│   │   ├── services/     # Serviços (API, cache)
│   │   ├── providers/    # State management
│   │   ├── screens/      # Telas
│   │   ├── widgets/      # Widgets reutilizáveis
│   │   └── utils/        # Utilitários
│   ├── main.dart
│   └── app.dart
├── assets/              # Assets (imagens, fonts)
├── test/               # Testes
└── pubspec.yaml
```

#### Comandos Úteis
```bash
# Análise de código
flutter analyze

# Testes
flutter test

# Build para diferentes plataformas
flutter build apk          # Android APK
flutter build appbundle    # Android App Bundle
flutter build ios          # iOS
flutter build web          # Web
flutter build windows      # Windows
flutter build macos        # macOS
flutter build linux        # Linux
```

## 🐳 Docker Development

### Desenvolvimento com Docker
```bash
# Build todas as imagens
docker-compose build

# Iniciar todos os serviços
docker-compose up

# Apenas infraestrutura
docker-compose up -d postgres redis

# Ver logs
docker-compose logs -f backend
docker-compose logs -f admin

# Parar tudo
docker-compose down

# Limpeza completa
docker-compose down -v
```

### Docker Compose Profiles
```bash
# Apenas desenvolvimento
docker-compose --profile dev up

# Desenvolvimento + monitoramento
docker-compose --profile dev --profile monitoring up

# Produção
docker-compose --profile prod up
```

## 🔑 Configuração de APIs Externas

### 1. Mapbox
1. Criar conta em [mapbox.com](https://www.mapbox.com)
2. Gerar token de acesso
3. Configurar em:
   - `backend/.env`: `MAPBOX_TOKEN=pk.your-token`
   - `admin/.env.local`: `NEXT_PUBLIC_MAPBOX_TOKEN=pk.your-token`
   - `app/lib/config/env.dart`: `mapboxToken = 'pk.your-token'`

### 2. Mercado Pago
1. Criar conta de desenvolvedor em [developers.mercadopago.com](https://developers.mercadopago.com)
2. Obter credenciais de teste/produção
3. Configurar em:
   - `backend/.env`: 
     ```
     MERCADO_PAGO_ACCESS_TOKEN=your-access-token
     MERCADO_PAGO_PUBLIC_KEY=your-public-key
     ```

### 3. Firebase (para Flutter)
1. Criar projeto no [Firebase Console](https://console.firebase.google.com)
2. Configurar app Android/iOS
3. Baixar `google-services.json` (Android) e `GoogleService-Info.plist` (iOS)
4. Habilitar Cloud Messaging para push notifications

## 🧪 Testes

### Backend (Go)
```bash
cd backend

# Testes unitários
go test ./...

# Testes com cobertura
go test -coverprofile=coverage.out ./...

# Testes de integração
go test -tags=integration ./...

# Benchmark
go test -bench=. ./...
```

### Admin (Next.js)
```bash
cd admin

# Testes unitários
npm test

# Testes E2E
npm run test:e2e

# Testes com watch
npm run test:watch

# Cobertura
npm run test:coverage
```

### App (Flutter)
```bash
cd app

# Testes unitários
flutter test

# Testes com cobertura
flutter test --coverage

# Testes de integração
flutter test integration_test/

# Testes de widget
flutter test test/widget_test.dart
```

## 🏗️ Build e Deploy

### Builds Locais
```bash
# Backend
cd backend
go build -o bin/server main.go

# Admin
cd admin
npm run build

# Flutter
cd app
flutter build web
flutter build apk --release
```

### Deploy via CI/CD
O projeto usa GitHub Actions para CI/CD automático:

1. **Push para `develop`**: Deploy automático para staging
2. **Push para `main`**: Deploy automático para produção
3. **Pull Request**: Executa testes e validações

### Deploy Manual (AWS)
```bash
# Configurar AWS CLI
aws configure

# Deploy infraestrutura
cd infra/terraform
terraform init
terraform apply -var-file=environments/staging.tfvars

# Deploy aplicações
cd ../..
./scripts/deploy-staging.sh
```

## 📊 Monitoramento

### Logs
```bash
# Backend logs
docker-compose logs -f backend

# Admin logs
docker-compose logs -f admin

# Database logs
docker-compose logs -f postgres
```

### Métricas
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3001 (admin/admin)
- **Application metrics**: http://localhost:8080/metrics

### Health Checks
- **Backend API**: http://localhost:8080/health
- **Database**: http://localhost:8080/health/db
- **Redis**: http://localhost:8080/health/redis

## 🛠️ Ferramentas de Desenvolvimento

### IDEs Recomendadas
- **Backend**: VS Code com extensão Go
- **Admin**: VS Code com extensões React/TypeScript
- **Flutter**: VS Code ou Android Studio
- **Database**: DBeaver ou pgAdmin

### Extensões VS Code Úteis
```json
{
  "recommendations": [
    "golang.go",
    "bradlc.vscode-tailwindcss",
    "esbenp.prettier-vscode",
    "ms-typescript.vscode-typescript-next",
    "dart-code.flutter",
    "dart-code.dart-code",
    "ms-vscode.vscode-json"
  ]
}
```

### Git Hooks
```bash
# Instalar pre-commit
pip install pre-commit

# Configurar hooks
pre-commit install

# Executar manualmente
pre-commit run --all-files
```

## 🔍 Debugging

### Backend
```bash
# Executar com debugger
dlv debug main.go

# Debug com VS Code
# Usar configuração launch.json
```

### Admin
```bash
# Debug mode
npm run dev

# Usar dev tools do navegador
# React Developer Tools extension
```

### Flutter
```bash
# Debug mode
flutter run --debug

# Debug com VS Code
# Usar F5 ou configuração launch.json
```

## 📚 Recursos Adicionais

### Documentação
- [Backend API docs](./api/README.md)
- [Admin UI components](./admin/docs/components.md)
- [Flutter app architecture](./app/docs/architecture.md)

### Links Úteis
- [Fiber Documentation](https://docs.gofiber.io/)
- [Next.js 14 Docs](https://nextjs.org/docs)
- [Flutter Documentation](https://docs.flutter.dev/)
- [PostGIS Documentation](https://postgis.net/docs/)
- [Mapbox GL JS](https://docs.mapbox.com/mapbox-gl-js/)

### Comunidade
- [Discord](https://discord.gg/mapfactions)
- [GitHub Discussions](https://github.com/macielcr7/map-area-factions/discussions)
- [Stack Overflow Tag](https://stackoverflow.com/questions/tagged/map-area-factions)

## 🆘 Troubleshooting

### Problemas Comuns

#### 1. "Connection refused" ao PostgreSQL
```bash
# Verificar se container está rodando
docker-compose ps postgres

# Verificar logs
docker-compose logs postgres

# Recriar container
docker-compose down postgres
docker-compose up -d postgres
```

#### 2. "Module not found" no Go
```bash
# Limpar cache
go clean -modcache

# Reinstalar dependências
go mod download
go mod tidy
```

#### 3. "Command not found" no Flutter
```bash
# Verificar instalação
flutter doctor

# Adicionar ao PATH
export PATH="$PATH:$HOME/flutter/bin"
```

#### 4. Problemas de CORS
- Verificar configuração em `backend/.env`: `CORS_ORIGINS`
- Verificar URL da API em `admin/.env.local`: `NEXT_PUBLIC_API_URL`

### Logs de Debug
```bash
# Backend com debug
DEBUG=true go run main.go

# Admin com debug
DEBUG=true npm run dev

# Flutter com verbose
flutter run --verbose
```

### Suporte
Se os problemas persistirem:
1. Verificar [GitHub Issues](https://github.com/macielcr7/map-area-factions/issues)
2. Criar nova issue com detalhes do problema
3. Contactar no Discord da comunidade
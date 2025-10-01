# Workflow de Desenvolvimento - Map Area Factions

**Guia para desenvolvimento contínuo do projeto**

## 🔄 Fluxo de Trabalho Diário

### 1. Setup Inicial (Sempre executar primeiro)
```bash
# Navegar para o projeto
cd /Users/macielcr7/Desktop/dev/maciel/map-area-factions

# Verificar status dos containers
docker-compose ps

# Se necessário, iniciar todos os serviços
docker-compose up -d

# Verificar se tudo está funcionando
curl http://localhost:8080/health
curl http://localhost:3000
```

### 2. Desenvolvimento Backend
```bash
# Para mudanças no backend
cd backend

# Executar localmente (opcional)
go run main.go

# Ou usar Docker (recomendado)
docker-compose build backend
docker-compose up -d backend

# Verificar logs
docker-compose logs -f backend
```

### 3. Desenvolvimento Admin
```bash
# Para mudanças no admin
cd admin

# Executar localmente (opcional)
npm run dev

# Ou usar Docker (recomendado)
docker-compose build admin
docker-compose up -d admin

# Verificar logs
docker-compose logs -f admin
```

### 4. Desenvolvimento Mobile
```bash
# Para mudanças no mobile
cd app

# Executar Flutter
flutter run

# Ou usar emulador
flutter run -d android
flutter run -d ios
```

## 🛠️ Comandos de Desenvolvimento

### Backend (Go)
```bash
# Executar testes
cd backend
go test ./...

# Executar com coverage
go test -cover ./...

# Linting
golangci-lint run

# Build
go build -o bin/backend main.go

# Executar migrations
go run cmd/migrate/main.go

# Executar seeds
go run cmd/seed/main.go
```

### Admin (Next.js)
```bash
# Instalar dependências
cd admin
npm install

# Desenvolvimento
npm run dev

# Build
npm run build

# Testes
npm test

# Linting
npm run lint

# Type checking
npm run type-check
```

### Mobile (Flutter)
```bash
# Instalar dependências
cd app
flutter pub get

# Executar
flutter run

# Build
flutter build apk
flutter build ios

# Testes
flutter test

# Análise
flutter analyze
```

## 🐛 Debug e Troubleshooting

### Problemas Comuns

#### 1. Container não inicia
```bash
# Verificar logs
docker-compose logs <service>

# Rebuild container
docker-compose build <service>
docker-compose up -d <service>

# Reset completo
docker-compose down
docker-compose up -d
```

#### 2. Banco de dados com problemas
```bash
# Verificar PostgreSQL
docker-compose ps postgres

# Conectar ao banco
docker-compose exec postgres psql -U postgres -d mapfactions

# Resetar banco
docker-compose down postgres
docker-compose up -d postgres

# Executar migrations
docker-compose exec backend go run cmd/migrate/main.go
```

#### 3. Admin não carrega
```bash
# Verificar se admin está rodando
docker-compose ps admin

# Rebuild admin
docker-compose build admin
docker-compose up -d admin

# Verificar logs
docker-compose logs admin

# Testar localmente
cd admin
npm run dev
```

#### 4. API não responde
```bash
# Verificar backend
docker-compose ps backend

# Verificar logs
docker-compose logs backend

# Testar health check
curl http://localhost:8080/health

# Testar login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@mapfactions.com", "password": "admin123"}'
```

## 📊 Monitoramento

### Verificar Status dos Serviços
```bash
# Status geral
docker-compose ps

# Logs em tempo real
docker-compose logs -f

# Logs específicos
docker-compose logs -f backend
docker-compose logs -f admin
docker-compose logs -f postgres
```

### Métricas de Performance
```bash
# Uso de recursos
docker stats

# Espaço em disco
docker system df

# Limpeza
docker system prune
```

## 🔧 Configuração de Ambiente

### Variáveis de Ambiente
```bash
# Backend
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=mapfactions
export JWT_SECRET=your-secret-key

# Admin
export NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
export NEXTAUTH_SECRET=your-nextauth-secret
export NEXTAUTH_URL=http://localhost:3000
```

### Docker Compose
```yaml
# docker-compose.yml
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=password
      - DB_NAME=mapfactions
    depends_on:
      - postgres
      - redis

  admin:
    build: ./admin
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://backend:8080/api/v1
      - NEXTAUTH_SECRET=your-nextauth-secret
      - NEXTAUTH_URL=http://localhost:3000
    depends_on:
      - backend
```

## 🧪 Testes

### Backend Tests
```bash
cd backend

# Executar todos os testes
go test ./...

# Executar com coverage
go test -cover ./...

# Executar testes específicos
go test ./internal/handlers/...

# Executar com verbose
go test -v ./...
```

### Admin Tests
```bash
cd admin

# Executar testes
npm test

# Executar com coverage
npm run test:coverage

# Executar testes E2E
npm run test:e2e
```

### Mobile Tests
```bash
cd app

# Executar testes
flutter test

# Executar com coverage
flutter test --coverage

# Executar testes de integração
flutter test integration_test/
```

## 📝 Git Workflow

### Branching Strategy
```bash
# Criar branch para feature
git checkout -b feature/editor-maps

# Fazer commits
git add .
git commit -m "feat: implement map editor with Mapbox"

# Push para remote
git push origin feature/editor-maps

# Criar Pull Request
# Merge após review
```

### Commit Messages
```
feat: nova funcionalidade
fix: correção de bug
docs: documentação
style: formatação
refactor: refatoração
test: testes
chore: tarefas de manutenção
```

## 🚀 Deploy

### Desenvolvimento
```bash
# Build e deploy local
docker-compose up -d

# Verificar se tudo está funcionando
curl http://localhost:8080/health
curl http://localhost:3000
```

### Staging (Futuro)
```bash
# Deploy para staging
# (Implementar quando necessário)
```

### Produção (Futuro)
```bash
# Deploy para produção
# (Implementar quando necessário)
```

## 📋 Checklist de Desenvolvimento

### Antes de Começar
- [ ] Verificar se containers estão rodando
- [ ] Testar endpoints básicos
- [ ] Verificar logs por erros
- [ ] Fazer backup do banco (se necessário)

### Durante o Desenvolvimento
- [ ] Fazer commits frequentes
- [ ] Testar funcionalidades implementadas
- [ ] Verificar logs de erro
- [ ] Documentar mudanças importantes

### Após Implementar
- [ ] Testar funcionalidade completa
- [ ] Verificar se não quebrou outras funcionalidades
- [ ] Executar testes
- [ ] Atualizar documentação
- [ ] Fazer commit final

### Antes de Finalizar
- [ ] Testar em diferentes browsers/devices
- [ ] Verificar performance
- [ ] Limpar código desnecessário
- [ ] Atualizar README se necessário

## 🎯 Próximas Tarefas

### Prioridade Alta
1. **Editor de Mapas** - Implementar Mapbox GL JS
2. **CRUD Interfaces** - Completar formulários
3. **Mobile App** - Setup Flutter básico

### Prioridade Média
1. **WebSocket** - Tempo real
2. **Assinaturas** - Sistema de pagamento
3. **Notificações** - Push notifications

### Prioridade Baixa
1. **Deploy** - AWS infrastructure
2. **Monitoramento** - Grafana/Prometheus
3. **CI/CD** - GitHub Actions

## 📞 Suporte

### Logs Importantes
```bash
# Backend logs
docker-compose logs backend

# Admin logs
docker-compose logs admin

# Database logs
docker-compose logs postgres
```

### Comandos de Debug
```bash
# Verificar containers
docker-compose ps

# Verificar recursos
docker stats

# Verificar rede
docker network ls
```

**Este workflow garante desenvolvimento eficiente e sem problemas!**

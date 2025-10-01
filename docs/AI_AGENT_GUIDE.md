# Guia do Agente IA - Map Area Factions

**Para qualquer IA que for trabalhar neste projeto**

## 🎯 Status Atual do Projeto (Janeiro 2025)

### ✅ IMPLEMENTADO E FUNCIONAL
- **Backend Go + Fiber:** 90% completo
- **Admin Interface Next.js:** 85% completo
- **Sistema de autenticação:** Funcionando
- **Dashboard com dados reais:** Funcionando
- **Docker containers:** Operacionais

### ❌ NÃO IMPLEMENTADO
- **Mobile App Flutter:** 0% implementado
- **Editor de mapas:** 0% implementado
- **Features avançadas:** 0% implementado
- **Deploy produção:** 0% implementado

## 🚀 Como Executar o Projeto

### 1. Iniciar o Sistema
```bash
# Navegar para o diretório do projeto
cd /Users/macielcr7/Desktop/dev/maciel/map-area-factions

# Iniciar todos os serviços
docker-compose up -d

# Verificar status
docker-compose ps
```

### 2. Acessar as Interfaces
- **Admin Interface:** http://localhost:3000
- **Backend API:** http://localhost:8080
- **Health Check:** http://localhost:8080/health

### 3. Credenciais de Teste
- **Admin:** admin@mapfactions.com / admin123
- **Moderador:** moderator@mapfactions.com / moderator123
- **Colaborador:** collaborator@mapfactions.com / collaborator123

## 📁 Estrutura do Projeto

```
map-area-factions/
├── backend/                 # API Go + Fiber
│   ├── internal/
│   │   ├── auth/           ✅ JWT, middleware
│   │   ├── handlers/       ✅ API endpoints
│   │   ├── models/         ✅ GORM models
│   │   ├── repository/     ✅ Data access
│   │   └── services/       ✅ Business logic
│   ├── migrations/         ✅ SQL migrations
│   ├── seeds/             ✅ Initial data
│   └── main.go            ✅ Entry point
├── admin/                  # Next.js 14 Admin
│   ├── src/
│   │   ├── app/           ✅ App Router pages
│   │   ├── components/   ✅ UI components
│   │   ├── hooks/         ✅ Custom hooks
│   │   └── lib/           ✅ Utilities
│   └── auth.ts            ✅ NextAuth config
├── app/                    # Flutter Mobile (estrutura básica)
├── docs/                   # Documentação
└── docker-compose.yml     ✅ Services config
```

## 🎯 Próximos Passos Prioritários

### 1. COMPLETAR ADMIN INTERFACE (2-3 semanas)

#### A. Editor de Mapas (1 semana)
**Objetivo:** Implementar editor de mapas com Mapbox GL JS

**Tarefas:**
1. **Instalar dependências:**
   ```bash
   cd admin
   npm install mapbox-gl @mapbox/mapbox-gl-draw
   npm install @types/mapbox-gl @types/mapbox__mapbox-gl-draw
   ```

2. **Criar componente MapEditor:**
   ```typescript
   // admin/src/components/maps/map-editor.tsx
   - Integrar Mapbox GL JS
   - Adicionar Mapbox GL Draw
   - Implementar controles de desenho
   - Conectar com API backend
   ```

3. **Implementar funcionalidades:**
   - Desenho de polígonos e polilinhas
   - Edição de geometrias existentes
   - Controle de camadas por facção
   - Salvamento no backend
   - Visualização por cores

4. **Integrar com página Maps:**
   ```typescript
   // admin/src/app/(main)/maps/page.tsx
   - Usar componente MapEditor
   - Adicionar sidebar com propriedades
   - Implementar seleção de facção
   ```

#### B. CRUD Interfaces Completas (1 semana)
**Objetivo:** Completar interfaces de usuários e facções

**Tarefas:**
1. **Página de Usuários:**
   ```typescript
   // admin/src/app/(main)/users/page.tsx
   - Tabela com paginação
   - Formulários de criação/edição
   - Filtros e busca
   - Modais de confirmação
   ```

2. **Página de Facções:**
   ```typescript
   // admin/src/app/(main)/factions/page.tsx
   - Cards com color picker
   - Formulários completos
   - Validação de dados
   - Preview de cores
   ```

3. **Componentes reutilizáveis:**
   ```typescript
   // admin/src/components/
   - DataTable com paginação
   - FormModal para CRUD
   - ColorPicker para facções
   - SearchFilter para busca
   ```

#### C. Melhorar UX/UI (3 dias)
**Objetivo:** Polir interface e experiência do usuário

**Tarefas:**
1. **Loading states avançados**
2. **Error boundaries**
3. **Toast notifications**
4. **Validação de formulários com Zod**
5. **Confirmação de ações**

### 2. MOBILE APP CORE (3-4 semanas)

#### A. Setup Flutter (1 semana)
**Objetivo:** Estrutura base do app Flutter

**Tarefas:**
1. **Configurar Flutter:**
   ```bash
   cd app
   flutter create .
   flutter pub get
   ```

2. **Instalar dependências:**
   ```yaml
   # pubspec.yaml
   dependencies:
     flutter:
       sdk: flutter
     http: ^1.1.0
     mapbox_gl: ^0.16.0
     riverpod: ^2.4.0
     go_router: ^12.0.0
   ```

3. **Estrutura base:**
   ```dart
   app/lib/src/
   ├── models/           # Data models
   ├── services/         # API services
   ├── providers/        # Riverpod providers
   ├── screens/          # UI screens
   ├── widgets/          # Reusable widgets
   └── config/           # App configuration
   ```

#### B. Autenticação (1 semana)
**Objetivo:** Sistema de login no mobile

**Tarefas:**
1. **Tela de login:**
   ```dart
   // app/lib/src/screens/auth/login_screen.dart
   - Formulário de login
   - Validação de campos
   - Integração com API
   - Navegação após login
   ```

2. **Gerenciamento de estado:**
   ```dart
   // app/lib/src/providers/auth_provider.dart
   - Estado de autenticação
   - Persistência de tokens
   - Logout automático
   ```

#### C. Visualização de Mapas (1 semana)
**Objetivo:** Mapa com facções no mobile

**Tarefas:**
1. **Tela de mapa:**
   ```dart
   // app/lib/src/screens/map_screen.dart
   - Integração Mapbox
   - Renderização de facções
   - Controles de zoom
   - Botão "minha localização"
   ```

2. **Serviços de API:**
   ```dart
   // app/lib/src/services/api_service.dart
   - Cliente HTTP
   - Autenticação automática
   - Cache de dados
   ```

#### D. Sistema de Busca (1 semana)
**Objetivo:** Busca por regiões no mobile

**Tarefas:**
1. **Tela de busca:**
   ```dart
   // app/lib/src/screens/search_screen.dart
   - Barra de busca
   - Resultados em lista
   - Filtros por facção
   - Navegação para mapa
   ```

### 3. FEATURES AVANÇADAS (2-3 semanas)

#### A. WebSocket para Tempo Real (1 semana)
**Objetivo:** Updates em tempo real

**Tarefas:**
1. **Backend WebSocket:**
   ```go
   // backend/internal/realtime/websocket.go
   - Servidor WebSocket
   - Room management
   - Event broadcasting
   - Redis pub/sub
   ```

2. **Frontend WebSocket:**
   ```typescript
   // admin/src/lib/websocket.ts
   - Cliente WebSocket
   - Auto-reconnection
   - Event handling
   - Real-time updates
   ```

#### B. Sistema de Assinaturas (1 semana)
**Objetivo:** Monetização do sistema

**Tarefas:**
1. **Backend assinaturas:**
   ```go
   // backend/internal/subscription/
   - Modelos de assinatura
   - Integração Mercado Pago
   - Webhook handling
   - Middleware de paywall
   ```

2. **Frontend assinaturas:**
   ```typescript
   // admin/src/app/subscriptions/
   - Interface de gestão
   - Dashboard de billing
   - Métricas de conversão
   ```

#### C. Push Notifications (1 semana)
**Objetivo:** Notificações para mobile

**Tarefas:**
1. **Backend notifications:**
   ```go
   // backend/internal/notifications/
   - Firebase Admin SDK
   - Templates de notificação
   - Targeting por localização
   ```

2. **Flutter notifications:**
   ```dart
   // app/lib/src/services/notification_service.dart
   - Firebase Messaging
   - Local notifications
   - Background handling
   ```

## 🔧 Comandos Essenciais

### Desenvolvimento
```bash
# Backend
cd backend
go run main.go

# Admin
cd admin
npm run dev

# Mobile
cd app
flutter run
```

### Docker
```bash
# Iniciar todos os serviços
docker-compose up -d

# Rebuild específico
docker-compose build admin
docker-compose up -d admin

# Logs
docker-compose logs -f backend
```

### Banco de Dados
```bash
# Conectar ao PostgreSQL
docker-compose exec postgres psql -U postgres -d mapfactions

# Executar migrations
docker-compose exec backend go run cmd/migrate/main.go

# Executar seeds
docker-compose exec backend go run cmd/seed/main.go
```

## 🐛 Troubleshooting Comum

### Problema: Admin não carrega
```bash
# Verificar se admin está rodando
docker-compose ps admin

# Rebuild admin
docker-compose build admin
docker-compose up -d admin

# Verificar logs
docker-compose logs admin
```

### Problema: API não responde
```bash
# Verificar backend
docker-compose ps backend

# Verificar logs
docker-compose logs backend

# Testar health check
curl http://localhost:8080/health
```

### Problema: Banco de dados
```bash
# Verificar PostgreSQL
docker-compose ps postgres

# Resetar banco
docker-compose down postgres
docker-compose up -d postgres

# Executar migrations
docker-compose exec backend go run cmd/migrate/main.go
```

## 📊 Métricas de Sucesso

### Para Admin Interface
- [ ] Editor de mapas funcionando
- [ ] CRUD de usuários completo
- [ ] CRUD de facções completo
- [ ] Validação de formulários
- [ ] Responsividade mobile

### Para Mobile App
- [ ] Login funcionando
- [ ] Mapa renderizando facções
- [ ] Busca por regiões
- [ ] Navegação fluida
- [ ] Performance otimizada

### Para Features Avançadas
- [ ] WebSocket funcionando
- [ ] Assinaturas ativas
- [ ] Push notifications
- [ ] Cache offline

## 🎯 Critérios de Aceite

### Editor de Mapas
- ✅ Desenho de polígonos
- ✅ Edição de geometrias
- ✅ Controle de camadas
- ✅ Salvamento no backend
- ✅ Visualização por cores

### CRUD Interfaces
- ✅ Formulários completos
- ✅ Tabelas com paginação
- ✅ Filtros e busca
- ✅ Validação de dados
- ✅ Feedback visual

### Mobile App
- ✅ Autenticação funcionando
- ✅ Mapa com facções
- ✅ Busca por regiões
- ✅ Navegação intuitiva
- ✅ Performance otimizada

## 📝 Notas Importantes

1. **Sempre testar após mudanças** - Use `docker-compose up -d` para rebuild
2. **Verificar logs** - Use `docker-compose logs -f <service>` para debug
3. **Banco de dados** - Sempre executar migrations após mudanças no schema
4. **Mobile** - Testar em device real, não apenas emulador
5. **Performance** - Monitorar uso de memória e CPU dos containers

## 🚀 Próximo Desenvolvedor

Se você é uma IA trabalhando neste projeto:

1. **Leia este guia completamente**
2. **Execute o sistema** para entender o estado atual
3. **Foque nas próximas prioridades** listadas acima
4. **Teste cada funcionalidade** antes de prosseguir
5. **Documente mudanças** importantes

**O projeto está funcional e pronto para desenvolvimento das próximas funcionalidades!**

# Contexto Inicial para Agentes IA

**Use este contexto ao iniciar qualquer trabalho no projeto Map Area Factions**

## 🎯 Contexto do Projeto

Você está trabalhando no projeto **Map Area Factions**, um sistema completo para mapeamento de áreas dominadas por facções com interface administrativa web e aplicativo móvel.

## 📊 Status Atual (Janeiro 2025)

### ✅ IMPLEMENTADO E FUNCIONAL
- **Backend (Go + Fiber):** 90% completo - API REST funcional
- **Admin Interface (Next.js 14):** 85% completo - Dashboard com dados reais
- **Sistema de autenticação:** Funcionando com JWT
- **Docker containers:** Operacionais
- **Banco PostgreSQL + PostGIS:** Configurado com dados reais

### ❌ NÃO IMPLEMENTADO
- **Mobile App (Flutter):** 0% implementado
- **Editor de mapas:** 0% implementado
- **Features avançadas:** 0% implementado

## 🚀 Como Executar o Projeto

```bash
# Navegar para o projeto
cd /Users/macielcr7/Desktop/dev/maciel/map-area-factions

# Iniciar todos os serviços
docker-compose up -d

# Verificar status
docker-compose ps

# Acessar admin interface
open http://localhost:3000

# Login com credenciais de teste
Email: admin@mapfactions.com
Senha: admin123
```

## 📁 Estrutura do Projeto

```
map-area-factions/
├── backend/                 # API Go + Fiber (90% completo)
│   ├── internal/
│   │   ├── auth/           ✅ JWT, middleware
│   │   ├── handlers/       ✅ API endpoints
│   │   ├── models/         ✅ GORM models
│   │   ├── repository/     ✅ Data access
│   │   └── services/       ✅ Business logic
│   ├── migrations/         ✅ SQL migrations
│   ├── seeds/             ✅ Initial data
│   └── main.go            ✅ Entry point
├── admin/                  # Next.js 14 Admin (85% completo)
│   ├── src/
│   │   ├── app/           ✅ App Router pages
│   │   ├── components/   ✅ UI components
│   │   ├── hooks/         ✅ Custom hooks
│   │   └── lib/           ✅ Utilities
│   └── auth.ts            ✅ NextAuth config
├── app/                    # Flutter Mobile (0% implementado)
├── docs/                   # Documentação completa
└── docker-compose.yml     ✅ Services config
```

## 🎯 Próximas Prioridades

### 1. Editor de Mapas (1 semana) - PRIORIDADE ALTA
**Objetivo:** Implementar editor de mapas com Mapbox GL JS

**Tarefas:**
- Instalar dependências: `mapbox-gl`, `@mapbox/mapbox-gl-draw`
- Criar componente `MapEditor` em `admin/src/components/maps/`
- Integrar com API backend
- Implementar desenho de polígonos/polilinhas
- Controle de camadas por facção
- Salvamento de geometrias

### 2. CRUD Interfaces Completas (1 semana) - PRIORIDADE ALTA
**Objetivo:** Completar interfaces de usuários e facções

**Tarefas:**
- Completar página de usuários (`admin/src/app/(main)/users/page.tsx`)
- Completar página de facções (`admin/src/app/(main)/factions/page.tsx`)
- Implementar formulários de criação/edição
- Adicionar tabelas com paginação
- Implementar filtros e busca
- Validação de formulários com Zod

### 3. Mobile App Core (3-4 semanas) - PRIORIDADE MÉDIA
**Objetivo:** App Flutter básico funcionando

**Tarefas:**
- Setup Flutter com estrutura base
- Implementar autenticação
- Integrar visualização de mapas
- Sistema de busca por região

## 🔧 Comandos Essenciais

### Desenvolvimento
```bash
# Rebuild específico
docker-compose build admin
docker-compose up -d admin

# Verificar logs
docker-compose logs -f backend
docker-compose logs -f admin

# Reset banco se necessário
docker-compose down postgres
docker-compose up -d postgres
```

### Troubleshooting
```bash
# Se admin não carrega
docker-compose build admin
docker-compose up -d admin
docker-compose logs admin

# Se API não responde
docker-compose logs backend
curl http://localhost:8080/health

# Se banco com problemas
docker-compose exec postgres psql -U postgres -d mapfactions
```

## 📚 Documentação Importante

### 🎯 Para Desenvolvedores
- **[AI_AGENT_GUIDE.md](./docs/AI_AGENT_GUIDE.md)** - Guia completo para qualquer IA
- **[IMPLEMENTATION_STATUS.md](./docs/IMPLEMENTATION_STATUS.md)** - Status atual detalhado
- **[FUNCTIONAL_FEATURES.md](./docs/FUNCTIONAL_FEATURES.md)** - O que está funcionando
- **[DEVELOPMENT_WORKFLOW.md](./docs/DEVELOPMENT_WORKFLOW.md)** - Como desenvolver

### 📋 Planejamento
- **[development-plan.md](./docs/development-plan.md)** - Roadmap por fases
- **[.ai-config.md](./.ai-config.md)** - Configuração rápida

## 🎯 Funcionalidades Implementadas

### Backend (Go + Fiber)
- ✅ API REST completa com 15+ endpoints
- ✅ Autenticação JWT com refresh tokens
- ✅ CRUD completo: usuários, facções, geometrias, relatórios, auditoria
- ✅ Dashboard endpoints com métricas reais
- ✅ Sistema de auditoria completo
- ✅ Middleware de segurança (CORS, rate limiting, auth)
- ✅ Banco PostgreSQL + PostGIS configurado
- ✅ Seeds com dados reais
- ✅ Docker container funcional

### Admin Interface (Next.js 14)
- ✅ Next.js 14 com App Router
- ✅ NextAuth integrado com backend
- ✅ Layout responsivo com sidebar
- ✅ Dashboard com dados reais (não mock)
- ✅ Validação de login obrigatória
- ✅ API client com interceptors
- ✅ Hooks para dashboard implementados
- ✅ Páginas: usuários, facções, mapas, auditoria, relatórios, configurações

## 🚧 O que Precisa ser Implementado

### Editor de Mapas (PRIORIDADE 1)
- ❌ Integração Mapbox GL JS
- ❌ Ferramentas de desenho (Mapbox GL Draw)
- ❌ Controle de camadas por facção
- ❌ Salvamento de geometrias
- ❌ Edição de geometrias existentes

### CRUD Interfaces (PRIORIDADE 2)
- ❌ Formulários de criação/edição
- ❌ Tabelas com paginação
- ❌ Filtros e busca
- ❌ Modais de confirmação
- ❌ Validação de formulários

### Mobile App (PRIORIDADE 3)
- ❌ Setup Flutter
- ❌ Autenticação
- ❌ Visualização de mapas
- ❌ Sistema de busca

## 📊 Métricas de Sucesso

### Para Editor de Mapas
- [ ] Desenho de polígonos funcionando
- [ ] Edição de geometrias
- [ ] Controle de camadas
- [ ] Salvamento no backend
- [ ] Visualização por cores

### Para CRUD Interfaces
- [ ] Formulários completos
- [ ] Tabelas com paginação
- [ ] Filtros e busca
- [ ] Validação de dados
- [ ] Feedback visual

### Para Mobile App
- [ ] Autenticação funcionando
- [ ] Mapa com facções
- [ ] Busca por regiões
- [ ] Navegação intuitiva
- [ ] Performance otimizada

## ⚠️ Importante

1. **Sempre testar após mudanças** - Use `docker-compose up -d`
2. **Verificar logs** - Use `docker-compose logs -f <service>`
3. **Banco de dados** - Executar migrations após mudanças no schema
4. **Mobile** - Testar em device real, não apenas emulador
5. **Performance** - Monitorar uso de memória e CPU dos containers

## 🎯 Instruções para a IA

1. **Leia a documentação completa** antes de começar
2. **Execute o sistema** para entender o estado atual
3. **Foque nas próximas prioridades** listadas acima
4. **Teste cada funcionalidade** antes de prosseguir
5. **Documente mudanças** importantes

**O projeto está funcional e pronto para desenvolvimento das próximas funcionalidades!**

---

**Use este contexto como base para qualquer trabalho no projeto Map Area Factions.**

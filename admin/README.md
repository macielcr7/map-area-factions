# Admin Panel - Map Area Factions

Painel administrativo em Next.js 14 para gerenciamento do sistema de mapeamento de áreas por facção.

## 🏗️ Arquitetura

### Stack Tecnológico
- **Framework**: Next.js 14 (App Router)
- **UI**: Tailwind CSS + Shadcn/ui
- **Mapas**: Mapbox GL JS + Mapbox GL Draw
- **Estado**: Zustand
- **Formulários**: React Hook Form + Zod
- **HTTP**: Axios
- **WebSocket**: Socket.io-client
- **Auth**: NextAuth.js

### Funcionalidades Principais

#### 🗺️ Editor de Mapas
- Desenho e edição de polígonos/polilinhas
- Ferramentas: desenhar, editar, mover, dividir, mesclar
- Snap automático a vias
- Importação/exportação de GeoJSON
- Preview de rótulos grandes
- Camadas por facção e risco

#### 📊 Dashboard
- Métricas em tempo real
- Estatísticas por município
- Gráficos de atividade
- Alertas e notificações

#### 👥 Gestão de Usuários
- CRUD de usuários
- Controle de roles (admin, moderator, collaborator)
- Logs de atividade
- Auditoria de alterações

#### 🎨 Configurações
- Gestão de facções (nome, cor, prioridade)
- Estilos de mapa
- Configurações de legendas
- Integração com provedores de mapas

## 📁 Estrutura do Projeto

```
admin/
├── src/
│   ├── app/                    # App Router (Next.js 14)
│   │   ├── (auth)/            # Rotas de autenticação
│   │   ├── dashboard/         # Dashboard principal
│   │   ├── maps/              # Editor de mapas
│   │   ├── users/             # Gestão de usuários
│   │   ├── factions/          # Gestão de facções
│   │   ├── audit/             # Logs de auditoria
│   │   └── settings/          # Configurações
│   ├── components/            # Componentes reutilizáveis
│   │   ├── ui/               # Componentes base (shadcn)
│   │   ├── maps/             # Componentes de mapa
│   │   ├── forms/            # Formulários
│   │   └── layout/           # Layout e navegação
│   ├── lib/                  # Utilitários e configurações
│   ├── hooks/                # React hooks customizados
│   ├── types/                # TypeScript types
│   └── store/                # Estado global (Zustand)
├── public/                   # Assets estáticos
├── docs/                     # Documentação específica
└── package.json
```

## 🚀 Desenvolvimento

### Pré-requisitos
- Node.js 18+
- npm ou yarn

### Instalação
```bash
# Instalar dependências
npm install

# Configurar variáveis de ambiente
cp .env.local.example .env.local

# Iniciar servidor de desenvolvimento
npm run dev
```

### Variáveis de Ambiente
```bash
# API
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_WS_URL=ws://localhost:8080

# Mapbox
NEXT_PUBLIC_MAPBOX_TOKEN=pk.your-mapbox-token

# NextAuth
NEXTAUTH_SECRET=your-nextauth-secret
NEXTAUTH_URL=http://localhost:3000

# Uploads
NEXT_PUBLIC_MAX_FILE_SIZE=10485760
```

## 🎨 Componentes Principais

### MapEditor
Componente principal para edição de mapas com:
- Integração Mapbox GL JS
- Ferramentas de desenho (Mapbox GL Draw)
- Gestão de camadas
- Event handlers para tempo real

### Dashboard
- Cards de métricas
- Gráficos interativos
- Tabelas de dados
- Notificações em tempo real

### UserManagement
- Lista de usuários com filtros
- Formulários de criação/edição
- Controle de permissões
- Histórico de ações

## 🔒 Autenticação

### NextAuth.js Configuration
- JWT com refresh tokens
- Roles e permissões
- Session management
- Proteção de rotas

### Middleware de Autorização
```typescript
// middleware.ts
import { withAuth } from "next-auth/middleware"

export default withAuth(
  function middleware(req) {
    // Lógica de autorização
  },
  {
    callbacks: {
      authorized: ({ token, req }) => {
        // Verificar permissões
      },
    },
  }
)
```

## 🗺️ Integração com Mapas

### Mapbox GL JS Setup
```typescript
import mapboxgl from 'mapbox-gl'
import MapboxDraw from '@mapbox/mapbox-gl-draw'

// Configuração do mapa
const map = new mapboxgl.Map({
  container: 'map',
  style: 'mapbox://styles/mapbox/streets-v11',
  center: [-38.5267, -3.7172], // Fortaleza
  zoom: 10
})

// Ferramentas de desenho
const draw = new MapboxDraw({
  displayControlsDefault: false,
  controls: {
    polygon: true,
    line_string: true,
    trash: true
  }
})
```

### Camadas Dinâmicas
- Polígonos por facção
- Estilos baseados em risco
- Clusters de incidentes
- Heatmaps temporais

## 📡 Tempo Real

### WebSocket Integration
```typescript
import { io } from 'socket.io-client'

const socket = io(process.env.NEXT_PUBLIC_WS_URL!)

socket.on('geometry.updated', (data) => {
  // Atualizar mapa em tempo real
})

socket.on('incident.created', (data) => {
  // Mostrar notificação
})
```

### Server-Sent Events (Fallback)
```typescript
const eventSource = new EventSource('/api/events')

eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data)
  // Processar evento
}
```

## 🧪 Testes

```bash
# Testes unitários
npm run test

# Testes E2E
npm run test:e2e

# Testes de componentes
npm run test:components
```

## 📦 Build e Deploy

### Build de Produção
```bash
# Build otimizado
npm run build

# Iniciar servidor de produção
npm run start
```

### Docker
```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

FROM node:18-alpine AS runner
WORKDIR /app
COPY --from=builder /app/node_modules ./node_modules
COPY . .
RUN npm run build

EXPOSE 3000
CMD ["npm", "start"]
```

## 🎯 Páginas Principais

### `/dashboard`
- Overview geral do sistema
- Métricas principais
- Atividade recente
- Mapa resumo

### `/maps`
- Editor principal de mapas
- Ferramentas de desenho
- Gestão de camadas
- Importação/exportação

### `/maps/[regionId]`
- Editor específico por região
- Histórico de alterações
- Colaboração em tempo real

### `/users`
- Lista de usuários
- Filtros e busca
- Gestão de roles

### `/factions`
- CRUD de facções
- Configuração de cores
- Prioridades de exibição

### `/audit`
- Logs de auditoria
- Filtros por usuário/ação
- Exportação de relatórios

### `/settings`
- Configurações gerais
- Integração com APIs
- Estilos de mapa
- Notificações

## 📱 Responsividade

O painel é totalmente responsivo com:
- Layout adaptável para desktop/tablet
- Componentes otimizados para mobile
- Touch gestures no mapa
- Navigation drawer para mobile

## 🔧 Configurações Avançadas

### Performance
- Lazy loading de componentes
- Virtualização de listas grandes
- Debounce em operações de busca
- Caching de requisições

### Acessibilidade
- ARIA labels
- Keyboard navigation
- Focus management
- Screen reader support
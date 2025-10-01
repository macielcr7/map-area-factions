# Map Area Factions 🗺️

Um sistema completo para mapeamento de áreas dominadas por facções, com interface administrativa web e aplicativo móvel/desktop. O sistema oferece visualização em tempo real de polígonos geoespaciais, busca por localização, sistema de assinaturas e atualizações em tempo real.

## 🎯 Status do Projeto (Janeiro 2025)

- **Backend (Go + Fiber):** ✅ 90% COMPLETO - API REST funcional
- **Admin Interface (Next.js):** ✅ 85% COMPLETO - Dashboard com dados reais  
- **Mobile App (Flutter):** ❌ 0% IMPLEMENTADO
- **Features Avançadas:** ❌ 0% IMPLEMENTADO

**Ver documentação completa em:** [docs/](./docs/)

## 🏗️ Arquitetura

### Stack Tecnológico
- **Backend**: Go com Fiber framework
- **Banco de Dados**: PostgreSQL + PostGIS
- **Cache/Real-time**: Redis com Pub/Sub
- **Admin Web**: Next.js 14 (App Router)
- **Apps Mobile/Desktop**: Flutter
- **Mapas**: Mapbox GL JS / Mapbox GL Flutter
- **Infraestrutura**: Docker + docker-compose, AWS (Terraform)

### Estrutura do Monorepo
```
├── backend/          # API Go com Fiber
├── admin/           # Dashboard Next.js 14
├── app/             # Aplicativo Flutter
├── infra/           # Terraform/Docker configs
├── docs/            # Documentação técnica
└── scripts/         # Scripts de automação
```

## 🚀 Início Rápido

### Pré-requisitos
- Docker e docker-compose
- Go 1.21+
- Node.js 18+
- Flutter 3.16+

### Desenvolvimento Local
```bash
# Clonar repositório
git clone https://github.com/macielcr7/map-area-factions.git
cd map-area-factions

# Iniciar todos os serviços
make dev

# Ou individualmente:
make dev-backend
make dev-admin
make dev-app
```

### Configuração de Ambiente
1. Copie `.env.example` para `.env` em cada subprojeto
2. Configure suas chaves de API (Mapbox, Mercado Pago)
3. Execute as migrações do banco de dados

## 📊 Funcionalidades

### 🔧 Painel Administrativo
- ✅ Desenho e edição de polígonos/polilinhas no mapa
- ✅ Gestão de facções e níveis de risco
- ✅ Sistema de moderação e aprovação
- ✅ Histórico e auditoria de alterações
- ✅ Importação/exportação de GeoJSON
- ✅ Dashboard com métricas e estatísticas

### 📱 Aplicativo do Cidadão
- ✅ Visualização de mapas com áreas por facção
- ✅ Busca por bairro/cidade/CEP
- ✅ Localização "perto de mim"
- ✅ Sistema de assinaturas (Pix/Cartão)
- ✅ Atualizações em tempo real
- ✅ Cache offline por cidade
- ✅ Notificações push

### 🔒 Segurança e Compliance
- ✅ Autenticação JWT com refresh tokens
- ✅ RBAC (Role-Based Access Control)
- ✅ Rate limiting e proteção anti-abuso
- ✅ Conformidade com LGPD
- ✅ Auditoria completa de operações
- ✅ Backup automatizado

## 💰 Planos de Assinatura

| Plano | Preço | Recursos |
|-------|-------|----------|
| Gratuito | R$ 0,00 | Visualização limitada, atualizações diárias |
| Mensal | R$ 14,90 | Acesso completo, tempo real |
| Trimestral | R$ 29,90 | Acesso completo + histórico |
| Vitalício | R$ 79,90 | Todos os recursos + suporte prioritário |

## 🗺️ Provedores de Mapas Suportados

O sistema abstrai o provedor de mapas, permitindo fácil troca:
- **Mapbox** (recomendado)
- MapTiler
- Google Maps
- OpenStreetMap

## 📚 Documentação

### 🎯 Para Desenvolvedores
- **[Guia do Agente IA](./docs/AI_AGENT_GUIDE.md)** - Guia completo para qualquer IA trabalhar no projeto
- **[Status de Implementação](./docs/IMPLEMENTATION_STATUS.md)** - Status atual detalhado
- **[Funcionalidades Implementadas](./docs/FUNCTIONAL_FEATURES.md)** - O que está funcionando
- **[Workflow de Desenvolvimento](./docs/DEVELOPMENT_WORKFLOW.md)** - Como desenvolver

### 📋 Planejamento
- **[Plano de Desenvolvimento](./docs/development-plan.md)** - Roadmap por fases
- **[Arquitetura e ADRs](./docs/adr/)** - Decisões arquiteturais
- **[API Reference](./docs/api/)** - Documentação da API

### 📄 Legal
- **[Política de Privacidade](./docs/privacy-policy.md)
- **[Termos de Uso](./docs/terms-of-service.md)

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/amazing-feature`)
3. Commit suas mudanças (`git commit -m 'Add amazing feature'`)
4. Push para a branch (`git push origin feature/amazing-feature`)
5. Abra um Pull Request

## ⚖️ Legal

⚠️ **Disclaimer**: Os dados são colaborativos e podem conter imprecisões. Não substituem orientação oficial. Use por sua conta e risco.

## 📄 Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🆘 Suporte

- 📧 Email: suporte@mapareafactions.com
- 💬 Discord: [Comunidade MAP](https://discord.gg/mapareafactions)
- 🐛 Issues: [GitHub Issues](https://github.com/macielcr7/map-area-factions/issues)
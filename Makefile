# Makefile para Map Area Factions

.PHONY: help dev dev-backend dev-admin dev-app build test lint clean docker-build docker-up docker-down

# Configurações padrão
DOCKER_COMPOSE = docker-compose
GO_VERSION = 1.21
NODE_VERSION = 18
FLUTTER_VERSION = 3.16

help: ## Mostra esta mensagem de ajuda
	@echo "Map Area Factions - Sistema de Mapeamento de Áreas por Facção"
	@echo ""
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

##@ Desenvolvimento
dev: docker-up ## Inicia todos os serviços em modo desenvolvimento
	@echo "🚀 Iniciando todos os serviços..."
	@make -j3 dev-backend dev-admin dev-app

dev-backend: ## Inicia apenas o backend Go
	@echo "🔧 Iniciando backend..."
	@cd backend && go run main.go

dev-admin: ## Inicia apenas o admin Next.js
	@echo "🌐 Iniciando admin..."
	@cd admin && npm run dev

dev-app: ## Inicia apenas o app Flutter
	@echo "📱 Iniciando app Flutter..."
	@cd app && flutter run -d chrome

##@ Build
build: build-backend build-admin build-app ## Constrói todos os projetos

build-backend: ## Constrói o backend Go
	@echo "🔨 Construindo backend..."
	@cd backend && go build -o bin/server main.go

build-admin: ## Constrói o admin Next.js
	@echo "🔨 Construindo admin..."
	@cd admin && npm run build

build-app: ## Constrói o app Flutter
	@echo "🔨 Construindo app Flutter..."
	@cd app && flutter build web

##@ Testes
test: test-backend test-admin test-app ## Executa todos os testes

test-backend: ## Executa testes do backend
	@echo "🧪 Testando backend..."
	@cd backend && go test ./...

test-admin: ## Executa testes do admin
	@echo "🧪 Testando admin..."
	@cd admin && npm test

test-app: ## Executa testes do app
	@echo "🧪 Testando app..."
	@cd app && flutter test

##@ Linting
lint: lint-backend lint-admin lint-app ## Executa linting em todos os projetos

lint-backend: ## Executa linting do backend
	@echo "✨ Linting backend..."
	@cd backend && golangci-lint run

lint-admin: ## Executa linting do admin
	@echo "✨ Linting admin..."
	@cd admin && npm run lint

lint-app: ## Executa linting do app
	@echo "✨ Linting app..."
	@cd app && flutter analyze

##@ Docker
docker-build: ## Constrói as imagens Docker
	@echo "🐳 Construindo imagens Docker..."
	@$(DOCKER_COMPOSE) build

docker-up: ## Inicia os containers (banco, redis, etc.)
	@echo "🐳 Iniciando containers de infraestrutura..."
	@$(DOCKER_COMPOSE) up -d postgres redis

docker-down: ## Para todos os containers
	@echo "🐳 Parando containers..."
	@$(DOCKER_COMPOSE) down

docker-logs: ## Mostra logs dos containers
	@$(DOCKER_COMPOSE) logs -f

##@ Banco de Dados
db-migrate: ## Executa migrações do banco
	@echo "📊 Executando migrações..."
	@cd backend && go run cmd/migrate/main.go

db-seed: ## Popula o banco com dados de exemplo
	@echo "🌱 Populando banco com dados de exemplo..."
	@cd backend && go run cmd/seed/main.go

db-reset: ## Reseta o banco de dados
	@echo "⚠️  Resetando banco de dados..."
	@$(DOCKER_COMPOSE) down -v postgres
	@$(DOCKER_COMPOSE) up -d postgres
	@sleep 5
	@make db-migrate
	@make db-seed

##@ Instalação
install: install-backend install-admin install-app ## Instala dependências de todos os projetos

install-backend: ## Instala dependências do backend
	@echo "📦 Instalando dependências do backend..."
	@cd backend && go mod download

install-admin: ## Instala dependências do admin
	@echo "📦 Instalando dependências do admin..."
	@cd admin && npm install

install-app: ## Instala dependências do app
	@echo "📦 Instalando dependências do app..."
	@cd app && flutter pub get

##@ Limpeza
clean: ## Remove arquivos de build
	@echo "🧹 Limpando arquivos de build..."
	@cd backend && rm -rf bin/
	@cd admin && rm -rf .next/ out/
	@cd app && flutter clean

##@ Deploy
deploy-staging: ## Deploy para staging
	@echo "🚀 Deploy para staging..."
	@cd infra && terraform apply -var-file=staging.tfvars

deploy-prod: ## Deploy para produção
	@echo "🚀 Deploy para produção..."
	@cd infra && terraform apply -var-file=production.tfvars

##@ Utilitários
setup: ## Setup completo do projeto
	@./scripts/setup.sh

quick-setup: docker-up install db-migrate db-seed ## Setup rápido sem verificações
	@echo "✅ Setup rápido completo!"

check: lint test ## Executa verificações (lint + test)

docs: ## Gera documentação
	@echo "📚 Gerando documentação..."
	@cd docs && make build

fmt: ## Formata código
	@echo "💅 Formatando código..."
	@cd backend && go fmt ./...
	@cd admin && npm run format
	@cd app && dart format .

##@ Informações
version: ## Mostra versões das ferramentas
	@echo "Versões das ferramentas:"
	@echo "Go: $(shell go version)"
	@echo "Node: $(shell node --version)"
	@echo "Flutter: $(shell flutter --version | head -1)"
	@echo "Docker: $(shell docker --version)"

status: ## Mostra status dos serviços
	@echo "Status dos serviços:"
	@$(DOCKER_COMPOSE) ps
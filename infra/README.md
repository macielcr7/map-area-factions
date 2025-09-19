# Infrastructure - Map Area Factions

Infraestrutura como código usando Terraform para deploy na AWS.

## 🏗️ Arquitetura AWS

### Componentes Principais
- **ECS Fargate**: Containers para backend e admin
- **RDS PostgreSQL**: Banco de dados com PostGIS
- **ElastiCache Redis**: Cache e real-time
- **S3 + CloudFront**: Storage e CDN
- **Application Load Balancer**: Load balancing
- **Route 53**: DNS e health checks
- **CloudWatch**: Monitoring e logs

### Ambientes
- **Development**: Single AZ, instâncias menores
- **Staging**: Multi-AZ, réplicas de produção
- **Production**: Multi-AZ, alta disponibilidade

## 📁 Estrutura

```
infra/
├── terraform/
│   ├── modules/           # Módulos reutilizáveis
│   │   ├── ecs/          # ECS Fargate
│   │   ├── rds/          # PostgreSQL + PostGIS
│   │   ├── redis/        # ElastiCache Redis
│   │   ├── s3/           # Storage
│   │   ├── cloudfront/   # CDN
│   │   ├── alb/          # Load Balancer
│   │   └── vpc/          # Networking
│   ├── environments/     # Configurações por ambiente
│   │   ├── dev/
│   │   ├── staging/
│   │   └── production/
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   └── versions.tf
├── aws/                  # Scripts AWS CLI
├── monitoring/           # Prometheus/Grafana
└── scripts/              # Scripts de deploy
```

## 🚀 Deploy

### Pré-requisitos
- Terraform 1.6+
- AWS CLI configurado
- Docker para builds

### Comandos
```bash
# Inicializar Terraform
terraform init

# Planejar deploy
terraform plan -var-file=environments/staging.tfvars

# Aplicar mudanças
terraform apply -var-file=environments/staging.tfvars

# Destruir infraestrutura
terraform destroy -var-file=environments/staging.tfvars
```

## 🔧 Configuração

### Variables (staging.tfvars)
```hcl
# Environment
environment = "staging"
region      = "us-east-1"

# Networking
vpc_cidr = "10.1.0.0/16"
availability_zones = ["us-east-1a", "us-east-1b"]

# ECS
ecs_cpu    = 512
ecs_memory = 1024
ecs_desired_count = 2

# RDS
rds_instance_class = "db.t3.medium"
rds_allocated_storage = 100
rds_multi_az = true

# Redis
redis_node_type = "cache.t3.micro"

# Domain
domain_name = "staging.mapfactions.com"
```

## 🌐 Networking

### VPC Configuration
- Public subnets: ALB, NAT Gateway
- Private subnets: ECS, RDS, Redis
- Internet Gateway para acesso público
- NAT Gateway para saída privada

### Security Groups
- ALB: HTTP/HTTPS (80, 443)
- ECS: Apenas ALB (8080, 3000)
- RDS: Apenas ECS (5432)
- Redis: Apenas ECS (6379)

## 📊 Monitoring

### CloudWatch
- Métricas de ECS (CPU, Memory, Network)
- Métricas de RDS (Connections, CPU, Storage)
- Métricas de Redis (CPU, Memory, Connections)
- Application metrics via Prometheus

### Alarms
- High CPU usage (>80%)
- High memory usage (>85%)
- Database connections (>80% max)
- Response time (>2s)
- Error rate (>5%)

## 🔒 Segurança

### IAM Roles
- ECS Task Role: acesso mínimo necessário
- ECS Execution Role: pull de imagens
- Lambda Roles: funções específicas

### Secrets Management
- AWS Secrets Manager para credenciais
- Parameter Store para configurações
- Rotação automática de secrets

### WAF
- Rate limiting
- IP blocking
- SQL injection protection
- XSS protection

## 📈 Scaling

### Auto Scaling
- ECS Service Auto Scaling baseado em CPU/Memory
- RDS Read Replicas para leitura
- Redis Cluster Mode para alta disponibilidade

### Performance
- CloudFront para assets estáticos
- ElastiCache para queries frequentes
- Connection pooling no backend

## 💾 Backup & Recovery

### RDS
- Automated backups (7 dias)
- Manual snapshots antes de deploys
- Cross-region backup replication

### ECS
- Rolling deployments
- Blue/green deployments para prod
- Rollback automático em falhas

## 🔍 Costs Optimization

### Strategies
- Spot instances para desenvolvimento
- Reserved instances para produção
- S3 Intelligent Tiering
- CloudWatch Logs retention policy
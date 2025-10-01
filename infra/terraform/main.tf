terraform {
  required_version = ">= 1.6"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1"
    }
  }

  backend "s3" {
    # Configure via -backend-config or environment variables
    # bucket = "map-factions-terraform-state"
    # key    = "infrastructure/terraform.tfstate"
    # region = "us-east-1"
    # encrypt = true
    # dynamodb_table = "map-factions-terraform-locks"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "map-area-factions"
      Environment = var.environment
      ManagedBy   = "terraform"
      Owner       = "map-factions-team"
    }
  }
}

# Data sources
data "aws_availability_zones" "available" {
  state = "available"
}

data "aws_caller_identity" "current" {}

# Random suffix for unique resource names
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# VPC Module
module "vpc" {
  source = "./modules/vpc"

  name               = "${var.project_name}-${var.environment}"
  cidr               = var.vpc_cidr
  availability_zones = data.aws_availability_zones.available.names
  
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
  
  enable_nat_gateway = var.enable_nat_gateway
  enable_vpn_gateway = false

  tags = local.common_tags
}

# Security Groups
resource "aws_security_group" "alb" {
  name_prefix = "${var.project_name}-${var.environment}-alb-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-alb-sg"
  })
}

resource "aws_security_group" "ecs" {
  name_prefix = "${var.project_name}-${var.environment}-ecs-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description     = "Backend API"
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  ingress {
    description     = "Admin Frontend"
    from_port       = 3000
    to_port         = 3000
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-ecs-sg"
  })
}

resource "aws_security_group" "rds" {
  name_prefix = "${var.project_name}-${var.environment}-rds-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description     = "PostgreSQL"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs.id]
  }

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-rds-sg"
  })
}

resource "aws_security_group" "redis" {
  name_prefix = "${var.project_name}-${var.environment}-redis-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description     = "Redis"
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs.id]
  }

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-${var.environment}-redis-sg"
  })
}

# S3 Module
module "s3" {
  source = "./modules/s3"

  bucket_name = "${var.project_name}-${var.environment}-${random_string.suffix.result}"
  environment = var.environment
  
  enable_versioning = var.environment == "production"
  enable_encryption = true
  
  tags = local.common_tags
}

# RDS Module
module "rds" {
  source = "./modules/rds"

  identifier = "${var.project_name}-${var.environment}"
  
  engine         = "postgres"
  engine_version = "15.4"
  instance_class = var.rds_instance_class
  
  allocated_storage     = var.rds_allocated_storage
  max_allocated_storage = var.rds_max_allocated_storage
  storage_encrypted     = true
  
  db_name  = var.db_name
  username = var.db_username
  password = var.db_password
  
  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = module.vpc.database_subnet_group
  
  backup_retention_period = var.rds_backup_retention_period
  backup_window          = var.rds_backup_window
  maintenance_window     = var.rds_maintenance_window
  
  multi_az               = var.rds_multi_az
  publicly_accessible    = false
  
  deletion_protection = var.environment == "production"
  skip_final_snapshot = var.environment != "production"
  
  tags = local.common_tags
}

# Redis Module
module "redis" {
  source = "./modules/redis"

  cluster_id = "${var.project_name}-${var.environment}"
  
  node_type            = var.redis_node_type
  num_cache_nodes      = var.redis_num_cache_nodes
  parameter_group_name = "default.redis7"
  port                 = 6379
  
  subnet_group_name  = module.vpc.elasticache_subnet_group
  security_group_ids = [aws_security_group.redis.id]
  
  tags = local.common_tags
}

# ECS Module
module "ecs" {
  source = "./modules/ecs"

  cluster_name = "${var.project_name}-${var.environment}"
  
  # Backend service
  backend_image = var.backend_image
  backend_cpu   = var.ecs_backend_cpu
  backend_memory = var.ecs_backend_memory
  backend_desired_count = var.ecs_backend_desired_count
  
  # Admin service
  admin_image = var.admin_image
  admin_cpu   = var.ecs_admin_cpu
  admin_memory = var.ecs_admin_memory
  admin_desired_count = var.ecs_admin_desired_count
  
  # Networking
  vpc_id          = module.vpc.vpc_id
  private_subnets = module.vpc.private_subnets
  security_groups = [aws_security_group.ecs.id]
  
  # Environment variables
  backend_environment = [
    {
      name  = "ENV"
      value = var.environment
    },
    {
      name  = "DB_HOST"
      value = module.rds.endpoint
    },
    {
      name  = "REDIS_HOST"
      value = module.redis.primary_endpoint
    },
    {
      name  = "S3_BUCKET"
      value = module.s3.bucket_name
    }
  ]
  
  backend_secrets = [
    {
      name      = "DB_PASSWORD"
      valueFrom = aws_secretsmanager_secret.db_password.arn
    },
    {
      name      = "JWT_SECRET"
      valueFrom = aws_secretsmanager_secret.jwt_secret.arn
    }
  ]
  
  tags = local.common_tags
}

# Application Load Balancer
module "alb" {
  source = "./modules/alb"

  name = "${var.project_name}-${var.environment}"
  
  vpc_id  = module.vpc.vpc_id
  subnets = module.vpc.public_subnets
  security_groups = [aws_security_group.alb.id]
  
  # Backend target group
  backend_target_group_arn = module.ecs.backend_target_group_arn
  
  # Admin target group
  admin_target_group_arn = module.ecs.admin_target_group_arn
  
  # SSL certificate
  certificate_arn = var.ssl_certificate_arn
  
  tags = local.common_tags
}

# CloudFront Module
module "cloudfront" {
  source = "./modules/cloudfront"

  domain_name = var.domain_name
  
  # ALB as origin
  alb_domain_name = module.alb.dns_name
  
  # S3 for static assets
  s3_bucket_domain = module.s3.bucket_domain_name
  
  # SSL certificate for CloudFront
  acm_certificate_arn = var.cloudfront_certificate_arn
  
  tags = local.common_tags
}

# Secrets Manager
resource "aws_secretsmanager_secret" "db_password" {
  name_prefix = "${var.project_name}-${var.environment}-db-password-"
  description = "Database password for ${var.environment} environment"
  
  tags = local.common_tags
}

resource "aws_secretsmanager_secret_version" "db_password" {
  secret_id     = aws_secretsmanager_secret.db_password.id
  secret_string = var.db_password
}

resource "aws_secretsmanager_secret" "jwt_secret" {
  name_prefix = "${var.project_name}-${var.environment}-jwt-secret-"
  description = "JWT secret for ${var.environment} environment"
  
  tags = local.common_tags
}

resource "aws_secretsmanager_secret_version" "jwt_secret" {
  secret_id     = aws_secretsmanager_secret.jwt_secret.id
  secret_string = var.jwt_secret
}

# Locals
locals {
  common_tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}
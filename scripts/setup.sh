#!/bin/bash

# Map Area Factions - Setup Script
# This script sets up the complete development environment

set -e

echo "🗺️  Map Area Factions - Setup Script"
echo "=================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are installed
check_requirements() {
    print_status "Checking requirements..."
    
    local missing_tools=()
    
    if ! command -v docker &> /dev/null; then
        missing_tools+=("docker")
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        missing_tools+=("docker-compose")
    fi
    
    if ! command -v go &> /dev/null; then
        missing_tools+=("go")
    fi
    
    if ! command -v node &> /dev/null; then
        missing_tools+=("node")
    fi
    
    if ! command -v flutter &> /dev/null; then
        missing_tools+=("flutter")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        print_error "Missing required tools: ${missing_tools[*]}"
        print_status "Please install the missing tools and run this script again."
        exit 1
    fi
    
    print_success "All required tools are installed!"
}

# Setup environment files
setup_env_files() {
    print_status "Setting up environment files..."
    
    # Backend
    if [ ! -f "backend/.env" ]; then
        cp backend/.env.example backend/.env
        print_success "Created backend/.env"
    else
        print_warning "backend/.env already exists, skipping..."
    fi
    
    # Admin
    if [ ! -f "admin/.env.local" ]; then
        cp admin/.env.local.example admin/.env.local
        print_success "Created admin/.env.local"
    else
        print_warning "admin/.env.local already exists, skipping..."
    fi
    
    print_warning "⚠️  Please update the environment files with your API keys:"
    echo "   - backend/.env: Add your Mapbox token and other credentials"
    echo "   - admin/.env.local: Add your Mapbox token and other credentials"
}

# Start infrastructure services
start_infrastructure() {
    print_status "Starting infrastructure services (PostgreSQL + Redis)..."
    
    docker-compose up -d postgres redis
    
    print_status "Waiting for services to be ready..."
    sleep 10
    
    # Check if PostgreSQL is ready
    max_attempts=30
    attempt=1
    while [ $attempt -le $max_attempts ]; do
        if docker-compose exec -T postgres pg_isready -U postgres &> /dev/null; then
            print_success "PostgreSQL is ready!"
            break
        fi
        
        if [ $attempt -eq $max_attempts ]; then
            print_error "PostgreSQL failed to start within timeout"
            exit 1
        fi
        
        print_status "Waiting for PostgreSQL... (attempt $attempt/$max_attempts)"
        sleep 2
        ((attempt++))
    done
    
    # Check if Redis is ready
    attempt=1
    while [ $attempt -le $max_attempts ]; do
        if docker-compose exec -T redis redis-cli ping &> /dev/null; then
            print_success "Redis is ready!"
            break
        fi
        
        if [ $attempt -eq $max_attempts ]; then
            print_error "Redis failed to start within timeout"
            exit 1
        fi
        
        print_status "Waiting for Redis... (attempt $attempt/$max_attempts)"
        sleep 2
        ((attempt++))
    done
}

# Setup backend
setup_backend() {
    print_status "Setting up Go backend..."
    
    cd backend
    
    # Install dependencies
    print_status "Installing Go dependencies..."
    go mod download
    
    # Run migrations
    print_status "Running database migrations..."
    go run cmd/migrate/main.go
    
    # Seed database
    print_status "Seeding database with sample data..."
    go run cmd/seed/main.go
    
    cd ..
    print_success "Backend setup complete!"
}

# Setup admin
setup_admin() {
    print_status "Setting up Next.js admin panel..."
    
    cd admin
    
    # Install dependencies
    print_status "Installing npm dependencies..."
    npm install
    
    cd ..
    print_success "Admin setup complete!"
}

# Setup Flutter app
setup_flutter() {
    print_status "Setting up Flutter app..."
    
    cd app
    
    # Get dependencies
    print_status "Getting Flutter dependencies..."
    flutter pub get
    
    # Check Flutter setup
    print_status "Checking Flutter setup..."
    flutter doctor
    
    cd ..
    print_success "Flutter setup complete!"
}

# Main setup function
main() {
    echo ""
    print_status "Starting setup process..."
    echo ""
    
    check_requirements
    echo ""
    
    setup_env_files
    echo ""
    
    start_infrastructure
    echo ""
    
    setup_backend
    echo ""
    
    setup_admin
    echo ""
    
    setup_flutter
    echo ""
    
    print_success "🎉 Setup complete!"
    echo ""
    print_status "Next steps:"
    echo "1. Update your API keys in the environment files"
    echo "2. Start the development servers:"
    echo "   - Backend: make dev-backend (or cd backend && go run main.go)"
    echo "   - Admin: make dev-admin (or cd admin && npm run dev)"
    echo "   - Flutter: make dev-app (or cd app && flutter run -d chrome)"
    echo ""
    print_status "Useful commands:"
    echo "   - make dev: Start all services"
    echo "   - make test: Run all tests"
    echo "   - make lint: Run all linters"
    echo "   - docker-compose logs -f: View service logs"
    echo ""
    print_status "URLs:"
    echo "   - Backend API: http://localhost:8080"
    echo "   - Admin Panel: http://localhost:3000"
    echo "   - Flutter Web: http://localhost:8080 (when running)"
    echo "   - Prometheus: http://localhost:9090"
    echo "   - Grafana: http://localhost:3001 (admin/admin)"
    echo ""
    print_warning "Remember to read the documentation in docs/ for detailed information!"
}

# Run main function
main "$@"
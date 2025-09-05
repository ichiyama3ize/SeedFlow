#!/bin/bash

# SeedFlow Knowledge Management Tool - Startup Script
# This script initializes and starts the Docker containers

set -e

echo "🌱 SeedFlow Knowledge Management Tool - Docker Setup"
echo "=================================================="

# Check if .env file exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Creating from template..."
    cp .env.example .env
    echo "📝 Please edit .env file with your API keys before running this script again."
    echo "   Minimum required: OPENAI_API_KEY"
    exit 1
fi

# Source environment variables
source .env

# Validate required environment variables
if [ -z "$OPENAI_API_KEY" ]; then
    echo "❌ Error: OPENAI_API_KEY is required in .env file"
    exit 1
fi

# Create necessary directories
echo "📁 Creating necessary directories..."
mkdir -p data logs exports backups config

# Set proper permissions
echo "🔒 Setting directory permissions..."
chmod 700 data config
chmod 755 logs exports backups

# Check if Docker and Docker Compose are available
if ! command -v docker &> /dev/null; then
    echo "❌ Error: Docker is not installed"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "❌ Error: Docker Compose is not installed"
    exit 1
fi

# Stop existing containers if running
echo "🛑 Stopping existing containers..."
docker-compose down --remove-orphans || true

# Build and start containers
echo "🔨 Building and starting containers..."
docker-compose up --build -d

# Wait for services to be healthy
echo "⏳ Waiting for services to be ready..."
timeout=60
counter=0

while [ $counter -lt $timeout ]; do
    if docker-compose ps | grep -q "Up (healthy)"; then
        ai_health=$(docker-compose exec -T ai-service curl -s -o /dev/null -w "%{http_code}" http://localhost:8001/ai/health 2>/dev/null || echo "000")
        go_health=$(docker-compose exec -T go-app wget -q -O /dev/null -S --spider http://localhost:8080/api/health 2>&1 | grep "200 OK" || echo "")
        
        if [ "$ai_health" = "200" ] && [ ! -z "$go_health" ]; then
            echo "✅ All services are healthy!"
            break
        fi
    fi
    
    counter=$((counter + 5))
    sleep 5
    echo "   ... waiting ($counter/${timeout}s)"
done

if [ $counter -ge $timeout ]; then
    echo "⚠️  Timeout waiting for services to be ready"
    echo "   You can check the logs with: docker-compose logs"
else
    echo ""
    echo "🎉 SeedFlow is now running!"
    echo "   Web UI: http://localhost:8080"
    echo "   AI Service: http://localhost:8001"
    echo ""
    echo "📊 Container Status:"
    docker-compose ps
    echo ""
    echo "📋 Useful commands:"
    echo "   View logs: docker-compose logs -f"
    echo "   Stop: docker-compose down"
    echo "   Restart: docker-compose restart"
    echo "   Update: docker-compose pull && docker-compose up -d"
fi
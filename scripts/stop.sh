#!/bin/bash

# SeedFlow Knowledge Management Tool - Stop Script
# This script stops and cleans up Docker containers

set -e

echo "🛑 SeedFlow Knowledge Management Tool - Shutdown"
echo "=============================================="

# Stop containers
echo "⏹️  Stopping containers..."
docker-compose down

# Show final status
echo "✅ SeedFlow has been stopped."
echo ""
echo "📋 To completely remove everything:"
echo "   Remove containers and images: docker-compose down --rmi all --volumes"
echo "   Remove data (⚠️  This deletes your knowledge database): rm -rf data/"
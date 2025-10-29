#!/bin/bash

# Medical Reports System - Development Start Script

echo "🏥 Medical Reports System - Starting Development Environment"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

echo "📦 Building and starting services..."
echo ""

# Start Docker services
docker-compose up --build -d

echo ""
echo "⏳ Waiting for services to be ready..."
sleep 5

echo ""
echo "✅ Services are starting!"
echo ""
echo "📍 Access the application:"
echo "   🌐 Frontend:  http://localhost:3000"
echo "   🔌 Backend:   http://localhost:8080"
echo "   🗄️  Database:  localhost:5433"
echo ""
echo "📝 To view logs:"
echo "   docker-compose logs -f"
echo ""
echo "🛑 To stop all services:"
echo "   docker-compose down"
echo ""

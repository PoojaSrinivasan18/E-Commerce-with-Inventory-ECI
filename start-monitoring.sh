#!/bin/bash

# 📊 Quick Monitoring Setup for Demo
# This script starts Prometheus and Grafana dashboards alongside your e-commerce services

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 Starting E-Commerce Monitoring Stack${NC}"

# Check if main services are running
if ! docker ps | grep -q "catalog-service\|inventoryservice\|customerservice"; then
    echo -e "${YELLOW}⚠️  Main e-commerce services not detected. Starting them first...${NC}"
    echo "Command: docker-compose up -d"
    docker-compose up -d
    echo "Waiting 30 seconds for services to initialize..."
    sleep 30
fi

# Start monitoring services
echo -e "${GREEN}📊 Starting Prometheus and Grafana...${NC}"
if [ -f "monitoring/docker-compose.monitoring.yml" ]; then
    docker-compose -f monitoring/docker-compose.monitoring.yml up -d
else
    echo -e "${YELLOW}⚠️  Monitoring compose file not found. Creating network and starting basic monitoring...${NC}"
    
    # Create network if it doesn't exist
    docker network create e-commerce-network 2>/dev/null || true
    
    # Start Prometheus
    docker run -d --name prometheus \
        --network e-commerce-network \
        -p 9090:9090 \
        -v "$(pwd)/monitoring/prometheus.yml:/etc/prometheus/prometheus.yml" \
        prom/prometheus:latest \
        --config.file=/etc/prometheus/prometheus.yml \
        --storage.tsdb.path=/prometheus \
        --web.enable-lifecycle || echo "Prometheus already running"
    
    # Start Grafana
    docker run -d --name grafana \
        --network e-commerce-network \
        -p 3000:3000 \
        -e GF_SECURITY_ADMIN_PASSWORD=admin123 \
        grafana/grafana:latest || echo "Grafana already running"
fi

# Wait for services to be ready
echo "Waiting for monitoring services to be ready..."
sleep 15

# Check service health
echo -e "${GREEN}✅ Monitoring Stack Status:${NC}"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep -E "(prometheus|grafana|catalog|inventory|customer)"

echo ""
echo -e "${GREEN}🎯 Access Your Dashboards:${NC}"
echo -e "  📈 Prometheus: ${BLUE}http://localhost:9090${NC}"
echo -e "  📊 Grafana:    ${BLUE}http://localhost:3000${NC} (admin/admin123)"
echo ""
echo -e "${GREEN}🔍 Key Demo Points:${NC}"
echo "  • Prometheus → Status → Targets (show service discovery)"
echo "  • Prometheus → Graph → Query: up{job=\"catalog-service\"}"  
echo "  • Grafana → Dashboards → E-Commerce Metrics"
echo "  • Show real-time metrics during inter-service calls"
echo ""
echo -e "${YELLOW}💡 Pro Tip: Run 'python3 interservice-workflow.py' in another terminal to generate metrics${NC}"
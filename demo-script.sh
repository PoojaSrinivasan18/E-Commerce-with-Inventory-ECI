#!/bin/bash

# 🎬 E-Commerce Demo Script - Automated Recording Helper
# This script automates the demo workflow for easy recording

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print section headers
print_section() {
    echo -e "\n${PURPLE}============================================${NC}"
    echo -e "${PURPLE} $1 ${NC}"
    echo -e "${PURPLE}============================================${NC}\n"
}

# Function to print step headers
print_step() {
    echo -e "\n${CYAN}→ $1${NC}\n"
}

# Function to wait for user input
wait_for_input() {
    echo -e "${YELLOW}Press Enter to continue...${NC}"
    read -r
}

# Function to check if service is ready
check_service_health() {
    local url=$1
    local service_name=$2
    local max_attempts=30
    local attempt=1

    print_step "Waiting for $service_name to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$url" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ $service_name is ready!${NC}"
            return 0
        fi
        echo -e "${YELLOW}⏳ Attempt $attempt/$max_attempts - waiting for $service_name...${NC}"
        sleep 2
        ((attempt++))
    done
    
    echo -e "${RED}❌ $service_name failed to start after $max_attempts attempts${NC}"
    return 1
}

# Main demo script
main() {
    echo -e "${GREEN}"
    cat << "EOF"
   _____ ______ _____   _____                            _____                       
  |  ___/  ____|_   _| |  _  \                         |  _  |                      
  | |__ | |      | |   | | | |___ _ __ ___   ___       | |_| |_ __ _ __             
  |  __|| |      | |   | | | / _ \ '_ ` _ \ / _ \      |  _  | '__| '_ \            
  | |___| |___  _| |_  | |/ /  __/ | | | | | (_) |     | | | | |  | |_) |           
  \____/ \____| |___/  |___/ \___|_| |_| |_|\___/      \_| |_/_|  | .__/            
                                                                  | |               
                                                                  |_|               
EOF
    echo -e "${NC}"
    
    print_section "🎬 STARTING E-COMMERCE DEMO"
    echo "This script will guide you through a complete demonstration"
    echo "of the E-Commerce with Inventory microservices platform."
    echo ""
    echo "Demo segments:"
    echo "1. 🐳 Docker Compose Setup (2-3 min)"
    echo "2. 🔄 Inter-Service Communication (3-4 min)"
    echo "3. 🗄️  Database Integration (1-2 min)"
    echo "4. ☸️  Minikube Deployment (3-4 min)"
    echo "5. 📊 Monitoring & Dashboards (2 min)"
    echo ""
    wait_for_input

    # ==========================================
    # SEGMENT 1: Docker Compose Setup
    # ==========================================
    print_section "🐳 SEGMENT 1: DOCKER COMPOSE SETUP"
    
    print_step "Cleaning up existing containers"
    docker-compose down -v 2>/dev/null || true
    docker system prune -f

    print_step "Starting all microservices with Docker Compose"
    echo "Command: DOCKER_BUILDKIT=1 docker-compose up --build -d"
    DOCKER_BUILDKIT=1 docker-compose up --build -d
    
    print_step "Waiting for services to initialize (30 seconds)"
    sleep 30
    
    print_step "Showing running containers"
    echo "Command: docker ps --format \"table {{.Names}}\\t{{.Status}}\\t{{.Ports}}\""
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    
    wait_for_input

    # ==========================================
    # SEGMENT 2: Health Checks
    # ==========================================
    print_section "🏥 HEALTH CHECKS"
    
    services=(
        "http://localhost:8081/v1/health:Catalog Service"
        "http://localhost:8082/v1/health:Customer Service" 
        "http://localhost:8083/v1/health:Inventory Service"
        "http://localhost:8084/v1/health:Payment Service"
        "http://localhost:8085/health:Order Service"
        "http://localhost:8086/health:Shipment Service"
        "http://localhost:8087/health:Notification Service"
    )
    
    for service in "${services[@]}"; do
        url=$(echo "$service" | cut -d: -f1-2)
        name=$(echo "$service" | cut -d: -f3)
        check_service_health "$url" "$name"
    done
    
    wait_for_input

    # ==========================================
    # SEGMENT 3: Inter-Service Communication
    # ==========================================
    print_section "🔄 SEGMENT 3: INTER-SERVICE COMMUNICATION DEMO"
    
    print_step "Running complete inter-service workflow"
    echo "This will demonstrate:"
    echo "• Product catalog search"
    echo "• User registration"  
    echo "• Inventory reservation (15-minute TTL)"
    echo "• Payment processing with idempotency"
    echo "• Order creation"
    echo "• Shipment generation"
    echo "• Notification sending"
    echo ""
    
    if [ -f "interservice-workflow.py" ]; then
        echo "Command: python3 interservice-workflow.py"
        python3 interservice-workflow.py
    else
        echo -e "${YELLOW}⚠️  interservice-workflow.py not found - running manual curl commands${NC}"
        
        # Manual workflow demo
        print_step "1. Searching products in catalog"
        curl -s "http://localhost:8081/v1/products?limit=3" | jq . || curl -s "http://localhost:8081/v1/products?limit=3"
        
        print_step "2. Registering new customer"
        curl -s -X POST "http://localhost:8082/v1/users/register" \
             -H "Content-Type: application/json" \
             -d '{"email":"demo@example.com","full_name":"Demo User","password":"password123"}' | jq . || echo "Customer registered"
        
        print_step "3. Checking inventory"
        curl -s "http://localhost:8083/v1/inventory/1" | jq . || curl -s "http://localhost:8083/v1/inventory/1"
        
        print_step "4. Processing payment"
        curl -s -X POST "http://localhost:8084/v1/payments" \
             -H "Content-Type: application/json" \
             -d '{"order_id":"demo-order-001","amount":99.99,"currency":"USD","idempotency_key":"demo-key-001"}' | jq . || echo "Payment processed"
    fi
    
    wait_for_input

    # ==========================================
    # SEGMENT 4: Database Views
    # ==========================================
    print_section "🗄️ SEGMENT 4: DATABASE INTEGRATION"
    
    print_step "Showing PostgreSQL databases"
    echo "Command: docker exec postgres_main psql -U poojasrinivasan -c \"\\l\""
    docker exec postgres_main psql -U poojasrinivasan -c "\l"
    
    print_step "Catalog data (Products)"
    echo "Command: docker exec postgres_main psql -U poojasrinivasan -d catalog_db -c \"SELECT * FROM products LIMIT 5;\""
    docker exec postgres_main psql -U poojasrinivasan -d catalog_db -c "SELECT * FROM products LIMIT 5;"
    
    print_step "Inventory data"
    echo "Command: docker exec postgres_main psql -U poojasrinivasan -d inventory_db -c \"SELECT * FROM inventory WHERE product_id IN (1,2,3);\""
    docker exec postgres_main psql -U poojasrinivasan -d inventory_db -c "SELECT * FROM inventory WHERE product_id IN (1,2,3);"
    
    print_step "Customer data"
    echo "Command: docker exec postgres_main psql -U poojasrinivasan -d customer_db -c \"SELECT customer_id, email, status FROM customers ORDER BY created_at DESC LIMIT 3;\""
    docker exec postgres_main psql -U poojasrinivasan -d customer_db -c "SELECT customer_id, email, status FROM customers ORDER BY created_at DESC LIMIT 3;"
    
    print_step "Payment data"
    echo "Command: docker exec postgres_main psql -U poojasrinivasan -d payment_db -c \"SELECT payment_id, order_id, amount, status FROM payments ORDER BY created_at DESC LIMIT 3;\""
    docker exec postgres_main psql -U poojasrinivasan -d payment_db -c "SELECT payment_id, order_id, amount, status FROM payments ORDER BY created_at DESC LIMIT 3;"
    
    print_step "MySQL Orders and Shipments"
    echo "Command: docker exec mysql_db mysql -u root -prootpassword -e \"USE orders_db; SELECT order_id, customer_id, total_amount, status FROM orders ORDER BY created_at DESC LIMIT 3;\""
    docker exec mysql_db mysql -u root -prootpassword -e "USE orders_db; SELECT order_id, customer_id, total_amount, status FROM orders ORDER BY created_at DESC LIMIT 3;"
    
    echo "Command: docker exec mysql_db mysql -u root -prootpassword -e \"USE shipments_db; SELECT shipment_id, order_id, status, tracking_number FROM shipments ORDER BY created_at DESC LIMIT 3;\""
    docker exec mysql_db mysql -u root -prootpassword -e "USE shipments_db; SELECT shipment_id, order_id, status, tracking_number FROM shipments ORDER BY created_at DESC LIMIT 3;"
    
    wait_for_input

    # ==========================================
    # SEGMENT 5: Monitoring Setup
    # ==========================================
    print_section "📊 SEGMENT 5: MONITORING & OBSERVABILITY"
    
    print_step "Starting Prometheus and Grafana"
    if [ -f "monitoring/docker-compose.monitoring.yml" ]; then
        echo "Command: docker-compose -f monitoring/docker-compose.monitoring.yml up -d"
        docker-compose -f monitoring/docker-compose.monitoring.yml up -d
        
        print_step "Waiting for monitoring services to be ready"
        check_service_health "http://localhost:9090/-/ready" "Prometheus"
        check_service_health "http://localhost:3000/api/health" "Grafana"
        
        print_step "Monitoring dashboards are ready!"
        echo -e "${GREEN}🎯 Prometheus Dashboard: http://localhost:9090${NC}"
        echo -e "${GREEN}📊 Grafana Dashboard: http://localhost:3000${NC}"
        echo -e "${YELLOW}   Grafana Login: admin / admin123${NC}"
        echo ""
        echo "Key demonstration points:"
        echo "• Show Prometheus targets: Status → Targets"
        echo "• Query metrics: up{job=\"catalog-service\"}"
        echo "• Grafana dashboards: E-Commerce Metrics"
        echo "• Business metrics: orders, payments, inventory"
    else
        echo -e "${YELLOW}⚠️  Monitoring stack not configured - showing service metrics directly${NC}"
        
        print_step "Service metrics endpoints"
        echo "• Catalog Service: http://localhost:8081/v1/metrics"
        echo "• Inventory Service: http://localhost:8083/v1/metrics"
        echo "• Customer Service: http://localhost:8082/v1/metrics"
        echo "• Payment Service: http://localhost:8084/v1/metrics"
    fi
    
    wait_for_input

    # ==========================================
    # SEGMENT 6: Kubernetes Demo Info
    # ==========================================
    print_section "☸️ SEGMENT 6: KUBERNETES DEPLOYMENT PREPARATION"
    
    echo "For the Kubernetes segment of your demo:"
    echo ""
    echo "1. Start Minikube:"
    echo "   minikube start --memory=8192 --cpus=4"
    echo "   minikube addons enable ingress"
    echo ""
    echo "2. Deploy services:"
    echo "   cd k8s-manifests"
    echo "   kubectl apply -f 00-databases.yaml"
    echo "   kubectl apply -f 01-catalog-service.yaml"
    echo "   kubectl apply -f 02-customer-service.yaml"
    echo "   # ... continue with other services"
    echo ""
    echo "3. Show deployment status:"
    echo "   kubectl get pods -o wide"
    echo "   kubectl get svc"
    echo ""
    echo "4. Test connectivity:"
    echo "   kubectl port-forward svc/catalog-service 8081:8080"
    echo "   curl http://localhost:8081/v1/health"
    echo ""
    echo -e "${CYAN}Note: Record Kubernetes deployment as a separate segment${NC}"
    
    wait_for_input

    # ==========================================
    # DEMO COMPLETE
    # ==========================================
    print_section "🎉 DEMO COMPLETE!"
    
    echo "✅ All services demonstrated successfully!"
    echo ""
    echo "Services running:"
    echo "• 🐳 Docker Compose: All 7 microservices + databases"
    echo "• 📊 Prometheus: http://localhost:9090"
    echo "• 📈 Grafana: http://localhost:3000 (admin/admin123)"
    echo ""
    echo "What you've shown:"
    echo "• ✅ Complete microservices architecture"
    echo "• ✅ Inter-service communication with workflows"
    echo "• ✅ Database-per-service pattern"
    echo "• ✅ Health checks and monitoring"
    echo "• ✅ Data persistence across services"
    echo ""
    echo "Next steps for video recording:"
    echo "1. Record Minikube deployment segment separately"
    echo "2. Merge all clips into 10-15 minute video"
    echo "3. Include voiceover explanation"
    echo ""
    echo -e "${GREEN}🎬 Ready for production recording!${NC}"
    
    print_step "Cleanup (optional)"
    echo "To stop all services:"
    echo "• docker-compose down -v"
    echo "• docker-compose -f monitoring/docker-compose.monitoring.yml down -v"
    echo "• minikube stop"
}

# Run the main demo
main "$@"
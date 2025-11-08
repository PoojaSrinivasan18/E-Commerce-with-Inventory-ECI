#!/bin/bash

# E-commerce Monitoring Stack Deployment
# Deploys Prometheus and Grafana for monitoring

set -e

echo "📊 Deploying E-commerce Monitoring Stack"
echo "======================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed. Please install kubectl first."
    exit 1
fi

# Deploy monitoring components
print_status "Deploying Prometheus configuration..."
kubectl apply -f monitoring/prometheus-config.yaml

print_status "Deploying Prometheus..."
kubectl apply -f monitoring/prometheus.yaml

print_status "Deploying Grafana configuration..."
kubectl apply -f monitoring/grafana-dashboard.yaml
kubectl apply -f monitoring/grafana.yaml

# Wait for monitoring services to be ready
print_status "Waiting for monitoring services to be ready..."
kubectl wait --for=condition=available deployment/prometheus -n monitoring --timeout=300s
kubectl wait --for=condition=available deployment/grafana -n monitoring --timeout=300s

print_success "Monitoring stack deployed successfully!"

# Get Minikube IP for URLs
if command -v minikube &> /dev/null; then
    MINIKUBE_IP=$(minikube ip)
    
    echo -e "\n📊 Monitoring Access Information:"
    echo "================================="
    echo "Minikube IP: $MINIKUBE_IP"
    echo ""
    echo "🔗 Monitoring URLs:"
    echo "  • Prometheus:  http://$MINIKUBE_IP:30090"
    echo "  • Grafana:     http://$MINIKUBE_IP:30030"
    echo ""
    echo "🔑 Grafana Login:"
    echo "  • Username: admin"
    echo "  • Password: admin123"
    echo ""
    echo "📈 Available Metrics:"
    echo "  • orders_placed_total - Total orders placed"
    echo "  • payments_completed_total - Successful payments"
    echo "  • payments_failed_total - Failed payments"
    echo "  • inventory_reservations_total - Inventory reservations"
    echo "  • inventory_available_stock - Available stock per product"
    echo "  • http_requests_total - HTTP request metrics"
    echo "  • up - Service availability"
    echo ""
    echo "🚨 Configured Alerts:"
    echo "  • ServiceDown - When a service goes down for >2min"
    echo "  • HighPaymentFailureRate - When payment success <90%"
    echo "  • LowInventoryAlert - When stock <10 units"
    
else
    print_warning "Minikube not found. Use 'kubectl port-forward' to access services:"
    echo ""
    echo "kubectl port-forward service/prometheus 9090:9090 -n monitoring"
    echo "kubectl port-forward service/grafana 3000:3000 -n monitoring"
fi

print_success "Monitoring stack is ready! 🎉"
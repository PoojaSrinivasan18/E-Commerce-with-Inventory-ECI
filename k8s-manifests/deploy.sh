#!/bin/bash

# E-commerce Kubernetes Deployment Script
# Deploys all microservices to Minikube cluster

set -e

echo "🚀 Starting E-commerce Kubernetes Deployment"
echo "============================================="

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

# Check if minikube is installed
if ! command -v minikube &> /dev/null; then
    print_error "Minikube is not installed. Please install minikube first."
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed. Please install kubectl first."
    exit 1
fi

# Start Minikube if not running
print_status "Checking Minikube status..."
if ! minikube status &> /dev/null; then
    print_status "Starting Minikube..."
    minikube start --driver=docker --memory=4096 --cpus=2
    print_success "Minikube started successfully"
else
    print_success "Minikube is already running"
fi

# Enable necessary addons
print_status "Enabling Minikube addons..."
minikube addons enable ingress
minikube addons enable metrics-server
print_success "Addons enabled"

# Set kubectl context to minikube
kubectl config use-context minikube

# Load Docker images into Minikube
print_status "Loading Docker images into Minikube..."
eval $(minikube docker-env)

# Check if images exist locally
IMAGES=(
    "e-commerce-with-inventory-eci-catalogservice:latest"
    "e-commerce-with-inventory-eci-inventoryservice:latest"
    "e-commerce-with-inventory-eci-customerservice:latest"
    "e-commerce-with-inventory-eci-payment_service:latest"
    "e-commerce-with-inventory-eci-order_service:latest"
    "e-commerce-with-inventory-eci-shipment_service:latest"
    "e-commerce-with-inventory-eci-notification_service:latest"
)

print_status "Checking for required Docker images..."
missing_images=()
for image in "${IMAGES[@]}"; do
    if ! docker images | grep -q "${image%:*}"; then
        missing_images+=("$image")
    fi
done

if [ ${#missing_images[@]} -gt 0 ]; then
    print_warning "The following images are missing and need to be built:"
    for image in "${missing_images[@]}"; do
        echo "  - $image"
    done
    
    read -p "Do you want to build missing images now? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Building missing Docker images..."
        cd ..
        docker-compose build
        cd k8s-manifests
        print_success "Docker images built successfully"
    else
        print_error "Cannot proceed without Docker images. Exiting."
        exit 1
    fi
fi

# Apply Kubernetes manifests
print_status "Applying Kubernetes manifests..."

# Apply in order (databases first, then services)
kubectl apply -f 01-postgres-init.yaml
kubectl apply -f 00-databases.yaml

# Wait for databases to be ready
print_status "Waiting for databases to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres-main -n ecommerce --timeout=300s
kubectl wait --for=condition=ready pod -l app=mysql-db -n ecommerce --timeout=300s
print_success "Databases are ready"

# Apply services
kubectl apply -f 02-catalog-service.yaml
kubectl apply -f 03-inventory-service.yaml
kubectl apply -f 04-customer-service.yaml
kubectl apply -f 05-payment-service.yaml
kubectl apply -f 06-order-service.yaml
kubectl apply -f 07-shipment-service.yaml
kubectl apply -f 08-notification-service.yaml

# Wait for services to be ready
print_status "Waiting for services to be ready..."
kubectl wait --for=condition=available deployment --all -n ecommerce --timeout=300s
print_success "All services are ready"

# Apply ingress
kubectl apply -f 09-ingress.yaml
print_success "Ingress controller configured"

# Display deployment status
print_status "Deployment Summary:"
echo "==================="

# Show pod status
kubectl get pods -n ecommerce

# Show service status  
echo -e "\nServices:"
kubectl get services -n ecommerce

# Show ingress status
echo -e "\nIngress:"
kubectl get ingress -n ecommerce

# Get Minikube IP for external access
MINIKUBE_IP=$(minikube ip)
print_success "Deployment completed successfully!"

echo -e "\n📋 Access Information:"
echo "======================"
echo "Minikube IP: $MINIKUBE_IP"
echo ""
echo "🔗 Service URLs (NodePort access):"
echo "  • Catalog Service:      http://$MINIKUBE_IP:30001"
echo "  • Inventory Service:    http://$MINIKUBE_IP:30002"
echo "  • Customer Service:     http://$MINIKUBE_IP:30003"
echo "  • Payment Service:      http://$MINIKUBE_IP:30004"
echo "  • Order Service:        http://$MINIKUBE_IP:30005"
echo "  • Shipment Service:     http://$MINIKUBE_IP:30006"
echo "  • Notification Service: http://$MINIKUBE_IP:30007"
echo ""
echo "🌐 Ingress URLs (add to /etc/hosts):"
echo "  echo '$MINIKUBE_IP ecommerce.local' | sudo tee -a /etc/hosts"
echo "  • API Gateway: http://ecommerce.local/api/v1/"
echo ""
echo "🔧 Useful Commands:"
echo "  • View pods:           kubectl get pods -n ecommerce"
echo "  • View services:       kubectl get services -n ecommerce"
echo "  • View logs:           kubectl logs -f deployment/<service-name> -n ecommerce"
echo "  • Scale service:       kubectl scale deployment <service-name> --replicas=3 -n ecommerce"
echo "  • Port forward:        kubectl port-forward service/<service-name> 8080:8080 -n ecommerce"

print_success "E-commerce platform is now running on Kubernetes! 🎉"
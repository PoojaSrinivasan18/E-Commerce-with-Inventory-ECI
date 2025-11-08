#!/bin/bash

# E-commerce Kubernetes Cleanup Script
# Removes all deployed resources

set -e

echo "🧹 E-commerce Kubernetes Cleanup"
echo "==============================="

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

# Confirm cleanup
echo "This will remove all E-commerce services and data from Kubernetes."
read -p "Are you sure you want to continue? (y/N): " -n 1 -r
echo

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_warning "Cleanup cancelled."
    exit 0
fi

print_status "Starting cleanup..."

# Delete E-commerce resources
if kubectl get namespace ecommerce &> /dev/null; then
    print_status "Removing E-commerce services..."
    kubectl delete namespace ecommerce
    print_success "E-commerce namespace deleted"
else
    print_warning "E-commerce namespace not found"
fi

# Delete Monitoring resources (if exists)
if kubectl get namespace monitoring &> /dev/null; then
    print_status "Removing monitoring stack..."
    kubectl delete namespace monitoring
    print_success "Monitoring namespace deleted"
else
    print_warning "Monitoring namespace not found"
fi

# Clean up any remaining resources
print_status "Cleaning up remaining resources..."
kubectl delete ingress --all --all-namespaces 2>/dev/null || true
kubectl delete pv --all 2>/dev/null || true

# Optional: Stop Minikube
read -p "Do you want to stop Minikube? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_status "Stopping Minikube..."
    minikube stop
    print_success "Minikube stopped"
    
    read -p "Do you want to delete Minikube cluster? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Deleting Minikube cluster..."
        minikube delete
        print_success "Minikube cluster deleted"
    fi
fi

print_success "Cleanup completed! 🧹"

echo -e "\n📝 What was cleaned up:"
echo "  ✅ E-commerce namespace and all services"
echo "  ✅ Monitoring namespace (if deployed)"
echo "  ✅ Persistent volumes and claims"
echo "  ✅ Ingress resources"
echo ""
echo "🔧 To redeploy:"
echo "  ./deploy.sh"
echo "  ./deploy-monitoring.sh  # (optional)"
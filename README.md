# E-commerce with Inventory (ECI) - Complete Documentation

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Prerequisites](#prerequisites)
4. [Quick Start](#quick-start)
5. [Docker Deployment](#docker-deployment)
6. [Kubernetes Deployment](#kubernetes-deployment)
7. [API Documentation](#api-documentation)
8. [Monitoring & Observability](#monitoring--observability)
9. [Testing & Validation](#testing--validation)
10. [Demo Instructions](#demo-instructions)
11. [Troubleshooting](#troubleshooting)

## Project Overview

This project implements a complete **E-commerce platform with Inventory Management** using microservices architecture.

### Completed Features

#### 1. Microservices (7 Services)
- **Catalog Service** - Product CRUD, search, pricing
- **Inventory Service** - Stock management, reservations, warehouse operations
- **Customer Service** - Customer management and authentication
- **Order Service** - Order creation, management, totals calculation
- **Payment Service** - Payment processing with idempotency
- **Shipping Service** - Shipment creation, tracking, status updates
- **Notification Service** - Email/SMS notifications for order lifecycle

#### 2. Database-Per-Service Architecture
- **PostgreSQL**: catalog_db, customer_db, inventory_db, payment_db, notifications
- **MySQL**: order_db (shared with shipments)
- **No cross-database joins** - each service owns its data
- **API-based data sharing** with minimal replication

#### 3. Advanced Inter-Service Workflows
- **Place Order Workflow**: Reserve → Pay → Ship → Notify
- **Reservation System**: 15-minute TTL with automatic cleanup
- **Idempotency**: Prevents duplicate orders and payments
- **Error Handling**: Proper rollback on failures
- **Warehouse Allocation**: Single-warehouse first strategy

#### 4. Production-Ready Features
- **API Versioning**: /v1 endpoints across all services
- **Health Checks**: Comprehensive health monitoring
- **Resource Management**: CPU/memory limits and requests
- **Observability**: Prometheus metrics + Grafana dashboards
- **Scalability**: Horizontal pod autoscaling ready

## Architecture

### Microservices Overview

```
Client Applications
        |
    API Gateway/Ingress
        |
    ┌─────────────────────────────────────────────────────────┐
    │                 Microservices                           │
    ├─────────────┬─────────────┬─────────────┬─────────────┤
    │ Catalog     │ Inventory   │ Customer    │ Payment     │
    │ Service     │ Service     │ Service     │ Service     │
    │ :8081       │ :8083       │ :8082       │ :8084       │
    ├─────────────┼─────────────┼─────────────┼─────────────┤
    │ Order       │ Shipment    │ Notification│             │
    │ Service     │ Service     │ Service     │             │
    │ :8085       │ :8086       │ :8087       │             │
    └─────────────┴─────────────┴─────────────┴─────────────┘
        |             |             |             |
    ┌─────────────────────────────────────────────────────────┐
    │                  Databases                              │
    ├─────────────┬─────────────┬─────────────┬─────────────┤
    │ PostgreSQL  │ PostgreSQL  │ PostgreSQL  │ PostgreSQL  │
    │ catalog_db  │ inventory_db│ customer_db │ payment_db  │
    ├─────────────┼─────────────┼─────────────┼─────────────┤
    │ MySQL       │ PostgreSQL  │             │             │
    │ orders_db   │notifications│             │             │
    │ shipments_db│             │             │             │
    └─────────────┴─────────────┴─────────────┴─────────────┘
```

### Database Schema per Service

| Service | Database | Tables | Purpose |
|---------|----------|--------|---------|
| Catalog | catalog_db (PG) | Products | Product catalog, pricing, categories |
| Inventory | inventory_db (PG) | Inventory, Reservations | Stock levels, warehouse operations |
| Customer | customer_db (PG) | Customers, Auth | Customer data, authentication |
| Payment | payment_db (PG) | Payments, Idempotency | Payment processing, fraud prevention |
| Order | order_db (MySQL) | Orders, OrderItems | Order management, line items |
| Shipment | order_db (MySQL) | Shipments | Shipping, tracking, delivery |
| Notification | notifications (PG) | NotificationLog | Audit trail for notifications |

## Prerequisites

### Software Requirements
- **Docker** (v20.10+)
- **Docker Compose** (v2.0+)  
- **Kubernetes/Minikube** (v1.24+)
- **Go** (v1.19+) for service development
- **Python** (v3.8+) for demo scripts
- **PostgreSQL Client** (psql)
- **MySQL Client** (mysql)

### System Resources
- **Memory**: 8GB+ RAM recommended
- **CPU**: 4+ cores recommended  
- **Storage**: 10GB+ available space
- **Network**: Internet access for image pulls

## Quick Start

### 1. Clone Repository
```bash
git clone https://github.com/PoojaSrinivasan18/E-Commerce-with-Inventory-ECI.git
cd E-Commerce-with-Inventory-ECI
```

### 2. Start All Services
```bash
# Clean any existing containers
docker-compose down -v
docker system prune -f

# Start all services
DOCKER_BUILDKIT=1 docker-compose up --build -d

# Wait for initialization (30 seconds)
sleep 30

# Check service status
docker ps
```

### 3. Verify Health
```bash
# Check all service health endpoints
curl http://localhost:8081/v1/health  # Catalog
curl http://localhost:8082/v1/health  # Customer  
curl http://localhost:8083/v1/health  # Inventory
curl http://localhost:8084/v1/health  # Payment
curl http://localhost:8085/health     # Order
curl http://localhost:8086/health     # Shipment
curl http://localhost:8087/health     # Notification
```

### 4. Run Demo Workflow
```bash
# Execute complete inter-service workflow
python3 interservice-workflow.py
```

## Docker Deployment

### Service Ports
- **Catalog Service**: 8081
- **Customer Service**: 8082  
- **Inventory Service**: 8083
- **Payment Service**: 8084
- **Order Service**: 8085
- **Shipment Service**: 8086
- **Notification Service**: 8087
- **PostgreSQL**: 5432
- **MySQL**: 3306

### Environment Variables
```bash
# Database Configuration
POSTGRES_USER=${POSTGRES_USER}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}

# Service Configuration
CATALOG_DB_HOST=postgres_main
INVENTORY_DB_HOST=postgres_main
CUSTOMER_DB_HOST=postgres_main
PAYMENT_DB_HOST=postgres_main
ORDER_DB_HOST=mysql_db
SHIPMENT_DB_HOST=mysql_db
```

### Build and Run
```bash
# Build all services
docker-compose build

# Start specific services
docker-compose up -d postgres_main mysql_db
docker-compose up -d catalog-service inventory-service
docker-compose up -d customer-service payment-service
docker-compose up -d order-service shipment-service notification_service

# View logs
docker-compose logs -f [service-name]

# Scale services
docker-compose up -d --scale catalog-service=3
```

## Kubernetes Deployment

### Prerequisites
```bash
# Start Minikube
minikube start --memory=8192 --cpus=4
minikube addons enable ingress
minikube addons enable metrics-server
```

### Deploy Services
```bash
# Navigate to Kubernetes manifests
cd k8s-manifests

# Deploy in order (databases first)
kubectl apply -f 00-databases.yaml
kubectl apply -f 01-catalog-service.yaml  
kubectl apply -f 02-customer-service.yaml
kubectl apply -f 03-inventory-service.yaml
kubectl apply -f 04-payment-service.yaml
kubectl apply -f 05-order-service.yaml
kubectl apply -f 06-shipment-service.yaml
kubectl apply -f 07-notification-service.yaml
kubectl apply -f 08-configmaps.yaml
kubectl apply -f 09-ingress.yaml

# Check deployment status
kubectl get pods -o wide
kubectl get svc
kubectl get ingress
```

### Access Services
```bash
# Port forward for testing
kubectl port-forward svc/catalog-service 8081:8080
kubectl port-forward svc/inventory-service 8083:8080

# Test connectivity  
curl http://localhost:8081/v1/health
curl http://localhost:8083/v1/health
```

## API Documentation

### Catalog Service (/v1)
- `GET /v1/products` - List products
- `GET /v1/products/{id}` - Get product details
- `POST /v1/products` - Create product
- `PUT /v1/products/{id}` - Update product
- `DELETE /v1/products/{id}` - Delete product
- `GET /v1/health` - Health check

### Inventory Service (/v1)  
- `GET /v1/inventory/{product_id}` - Get stock level
- `POST /v1/inventory/reserve` - Reserve inventory
- `POST /v1/inventory/release` - Release reservation
- `POST /v1/inventory/ship` - Mark as shipped
- `GET /v1/health` - Health check

### Customer Service (/v1)
- `POST /v1/users/register` - Register customer
- `POST /v1/users/login` - Customer login
- `GET /v1/users/{id}` - Get customer details
- `PUT /v1/users/{id}` - Update customer
- `GET /v1/health` - Health check

### Payment Service (/v1)
- `POST /v1/payments` - Process payment
- `GET /v1/payments/{id}` - Get payment status
- `POST /v1/payments/refund` - Process refund
- `GET /v1/health` - Health check

### Order Service
- `POST /v1/orders` - Create order
- `GET /v1/orders/{id}` - Get order details
- `PUT /v1/orders/{id}/status` - Update order status
- `GET /health` - Health check

### Shipment Service
- `POST /v1/shipments` - Create shipment
- `GET /v1/shipments/{id}` - Get shipment details  
- `PUT /v1/shipments/{id}/status` - Update shipment status
- `GET /health` - Health check

### Notification Service
- `POST /v1/notifications` - Send notification
- `GET /v1/notifications/{id}` - Get notification status
- `GET /health` - Health check

## Monitoring & Observability

### Prometheus Metrics
- Service health (`up`)
- HTTP request rates (`http_requests_total`)
- Response times (`http_request_duration_seconds`)
- Business metrics (`orders_created_total`, `payments_successful_total`)
- Database connections (`pg_stat_activity_count`)

### Grafana Dashboards
- **Service Health Status** - Real-time service availability
- **Request Rate & Latency** - Performance monitoring  
- **Business Metrics** - Orders, payments, inventory levels
- **Database Monitoring** - Connection pools and query performance

### Access Monitoring
```bash
# Deploy monitoring stack
kubectl apply -f monitoring/

# Port forward dashboards
kubectl port-forward -n monitoring svc/prometheus 9090:9090 &
kubectl port-forward -n monitoring svc/grafana 3000:3000 &

# Access URLs
# Prometheus: http://localhost:9090
# Grafana: http://localhost:3000 (admin/admin)
```

## Testing & Validation

### Health Checks
```bash
# Automated health check script
./scripts/health-check.sh

# Manual verification
curl -f http://localhost:8081/v1/health || echo "Catalog service down"
curl -f http://localhost:8082/v1/health || echo "Customer service down"  
curl -f http://localhost:8083/v1/health || echo "Inventory service down"
curl -f http://localhost:8084/v1/health || echo "Payment service down"
```

### Database Validation
```bash
# Verify PostgreSQL databases
docker exec postgres_main psql -U poojasrinivasan -c "\l"

# Check MySQL databases  
docker exec mysql_db mysql -u root -prootpassword -e "SHOW DATABASES;"

# Verify data persistence
docker exec postgres_main psql -U poojasrinivasan -d catalog_db -c "SELECT COUNT(*) FROM products;"
```

## Troubleshooting

### Common Issues

**Services not starting:**
```bash
# Check logs
docker-compose logs [service-name]

# Restart specific service
docker-compose restart [service-name]

# Full restart
docker-compose down -v && docker-compose up -d
```

**Database connection errors:**
```bash  
# Verify database containers
docker ps | grep -E "(postgres|mysql)"

# Check database logs
docker logs postgres_main
docker logs mysql_db

# Verify database creation
docker exec postgres_main psql -U poojasrinivasan -l
```

**Port conflicts:**
```bash
# Find processes using ports
lsof -ti:8081,8082,8083,8084,8085,8086,8087 | xargs kill -9

# Use different ports if needed
export CATALOG_PORT=9081
docker-compose up -d
```

**Kubernetes issues:**
```bash
# Reset Minikube
minikube delete && minikube start --memory=8192 --cpus=4

# Check pod status
kubectl describe pod [pod-name]

# View pod logs
kubectl logs -f [pod-name]
```

### Performance Tuning
```bash
# Increase memory limits
docker-compose --compatibility up -d

# Monitor resource usage
docker stats

# Kubernetes resource monitoring
kubectl top nodes
kubectl top pods
```

---

## Architecture Decisions

### Database Selection
- **PostgreSQL**: ACID compliance for financial data (payments, orders)
- **MySQL**: High performance for order processing
- **Database per service**: Complete data isolation and independence

### Service Communication
- **REST APIs**: Synchronous communication for immediate responses
- **Health checks**: Kubernetes-native health monitoring
- **Circuit breakers**: Resilience patterns for service failures

### Container Strategy
- **Multi-stage builds**: Optimized image sizes
- **Health checks**: Docker and Kubernetes native monitoring
- **Resource limits**: Prevent resource contention

### Security Considerations
- **Idempotency keys**: Prevent duplicate transactions
- **Input validation**: Comprehensive request validation
- **Database isolation**: No cross-database queries
- **Container security**: Non-root users and minimal base images

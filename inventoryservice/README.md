# inventoryservice

Exposes a REST API for CRUD operations on inventory items and supports persistence. Designed for containerized deployment and easy local development.

## Features
- CRUD operations for inventory items
- Simple JSON REST API
- Validation and error responses
- Configurable data store (e.g., PostgreSQL)
- K8s Friendly

## Prerequisites
- Git
- K8s (optional for containerized runs)
- A database such as PostgreSQL if not using in-memory store

## Quickstart (local)
1. Clone the repo:
    ```bash
    git clone <repo-url> inventoryservice
    cd inventoryservice
    ```
2. Configure environment variables (see Configuration to update the Postgres DB).
    ```bash
    check under deploy folder: postgres-db.yaml
    sudo iptables -I INPUT -p tcp -s 0.0.0.0/0 --dport 30001 -j ACCEPT
    oc port-forward postgres-558cf69f8b-n4g7g -n assignment 30001:5432 --address <<10.187.168.189>>
    ```
3. Run the service (example for a generic runtime):
    ```bash
    # example: go run main.go
    # or create image using docker build. Refer Quick Start
    ```

## Quickstart (Docker)
Build and run with Docker: Refer Docker File
```bash
docker build -t inventoryservice .

kubectl create -f postgres-db.yaml
sudo iptables -I INPUT -p tcp -s 0.0.0.0/0 --dport 30001 -j ACCEPT
oc port-forward postgres-558cf69f8b-n4g7g -n assignment 30001:5432 --address <<10.187.168.189>>
```

## Configuration
Common environment variables:
- Fill the details in the dbConfig.yaml file

## API (example)
- Base URL: http://localhost:8080
- Refer the Postman Collection under postman folder.


# ECI Shipments Service

A single-purpose Shipments microservice (FastAPI). It stores shipments in a MySQL database named `shipping_db` and can bootstrap its schema and data from CSV files.

## Repository layout

- `main.py` — FastAPI application exposing `/shipments` endpoints
- `db_utils.py` — helper to read `.env` and create MySQL connections
- `db_setup.py` — creates `shipping_db` and loads CSV data
- `requirements.txt` — Python dependencies
- `.env` — service environment (excluded from git)
- `csv_files/` — included sample CSV `Shipments.csv`
- `.gitignore` — ignores `.env`, pycache, etc.
- `README.md` — this file

## Overview

- Purpose: manage Shipments (create, list, fetch-by-id).
- DB: MySQL database `shipping_db` (created by `db_setup.py` if missing).
- CSV loader: `db_setup.py` will read CSVs and insert rows into tables.

## Requirements

- Python 3.10+
- MySQL server accessible with credentials in `.env`
- Install dependencies:

```bash
pip install -r requirements.txt
```

## Configuration

Edit the service-local `.env` file (already present in this folder). Typical variables:

- `DB_HOST` (default: 127.0.0.1)
- `DB_PORT` (default: 3306)
- `DB_USER` (default: root)
- `DB_PASSWORD` (change from the placeholder)
- `CSV_DIR` (optional: defaults to `./csv_files` in this folder)

db_utils.py loads the `.env` automatically from this service folder.

## Load CSVs into MySQL (bootstrap)

1. Ensure MySQL is running and `.env` contains valid credentials.
2. From this folder run:

```bash
python db_setup.py
```

This will:
- create `shipping_db` (if missing)
- create table `Shipments` (if missing)
- bulk-insert the CSV from the directory referenced by `CSV_DIR` (defaults to `./csv_files`)

If you want to use a different CSV location, set `CSV_DIR` in the `.env` or the environment before running.

## Docker & Local container instructions

The repository includes helper files so you can run this service in Docker with MySQL. Files to support Docker are:

- `Dockerfile` — builds the shipment_service image
- `entrypoint.sh` — waits for DB, runs `db_setup.py` and starts uvicorn
- `wait_for_db.py` — helper to wait for DB TCP readiness
- `docker-compose.yml` — local stack (MySQL + service)
- `.dockerignore` — excludes large or secret files from the image

Example `.env` (place under `shipment_service/.env`):

```
DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=rootpassword
DB_NAME=shipping_db
MYSQL_ROOT_PASSWORD=rootpassword
CSV_DIR=./csv_files
```

### Build & run with docker-compose (recommended)

From the `shipment_service/` folder:

```bash
cd shipment_service
docker-compose up --build
```

Run in background:

```bash
docker-compose up -d --build
```

Watch logs:

```bash
docker-compose logs -f shipment_service
```

Open the docs:

```text
http://127.0.0.1:8001/docs
```

### Initialize DB (wait for MySQL then run db_setup)

Before running the commands below, wait for MySQL to be ready (check `docker-compose logs mysql` or `docker logs mysql_db_shipments`).

Option A — run `db_setup.py` inside the running service container (recommended when using `docker-compose`):

```bash
docker-compose exec shipment_service python db_setup.py
```

Option B — run `db_setup.py` in a temporary container attached to the same network (useful if you only built the image):

```bash
docker run --rm --network $(docker network ls --filter name=shipment_service_default -q) --env-file .env <your-dockerhub-username>/shipment_service:latest python db_setup.py
```

Or, if you started the container manually by name:

```bash
docker exec -it shipment_service_container python db_setup.py
```

Make the entrypoint executable locally before building (macOS / zsh):

```bash
cd shipment_service
chmod +x entrypoint.sh
```

Tip: when running via `docker-compose` set `DB_HOST` to the compose service name (`mysql`) so the app uses the Docker network to reach the database.

### Manual Docker build & run (no compose)

Build image:

```bash
cd shipment_service
docker build -t <your-dockerhub-username>/shipment_service:latest .
```

Create network and run MySQL:

```bash
docker network create shipment_net || true
docker run -d --name mysql_db_shipments --network shipment_net \
	-e MYSQL_ROOT_PASSWORD=rootpassword \
	-e MYSQL_DATABASE=shipping_db \
	mysql:8.0
```

Run service (using `.env`):

```bash
docker run -d --name shipment_service_container --network shipment_net --env-file .env -p 8001:8001 <your-dockerhub-username>/shipment_service:latest
```

### Push image to Docker Hub / GHCR

Follow the same push steps as described in `order_service/README.md` — build, tag, login and push. Example:

```bash
docker tag <your-dockerhub-username>/shipment_service:latest <your-dockerhub-username>/shipment_service:1.0.0
docker push <your-dockerhub-username>/shipment_service:1.0.0
```

### CI: automatic build & publish (GHCR)

There's a GitHub Actions workflow included in `.github/workflows/docker-publish-shipment.yml` that will build and push the `shipment_service` image to GitHub Container Registry (GHCR) when you push to `main`.

Notes:
- The workflow uses `GITHUB_TOKEN` to authenticate to GHCR. You may need to allow `packages: write` for `GITHUB_TOKEN` in repository Settings → Actions → General. Alternatively, replace the login step to use a PAT stored in `secrets.CR_PAT` if you prefer.



## Run the API (development)

Start the app with uvicorn (default port 8001):

```bash
uvicorn main:app --reload --host 127.0.0.1 --port 8001
```

OpenAPI docs: `http://127.0.0.1:8001/docs`

## API Endpoints

Shipments

- `GET /shipments?limit={n}` — list shipments (default `limit=10`)
- `GET /shipments/{shipment_id}` — get shipment details
- `POST /shipments` — insert a new shipment (payload matches `Shipment` Pydantic model in `main.py`)
 - `PUT /shipments/{shipment_id}` — update a shipment (all fields supported)
 - `DELETE /shipments/{shipment_id}` — delete a shipment by id

## Examples

Minimal `POST /shipments` payload:

```json
{
	"shipment_id": 1,
	"order_id": 123,
	"carrier": "DHL",
	"status": "SHIPPED",
	"tracking_no": "TRACK123"
}
```

Example `PUT /shipments/{shipment_id}` payload:

```json
{
	"shipment_id": 1,
	"order_id": 123,
	"carrier": "UPS",
	"status": "IN_TRANSIT",
	"tracking_no": "NEWTRACK123",
	"shipped_at": "2025-01-01T12:00:00",
	"delivered_at": null
}
```

Endpoint examples (curl)

1) List shipments (GET /shipments)

```bash
curl "http://127.0.0.1:8001/shipments?limit=5"
```

Sample response (array):

```json
[
	{ "shipment_id": 1, "order_id": 128, "carrier": "DHL", "status": "DELIVERED", "tracking_no": "TRK6188", "shipped_at": "2023-09-25 09:07:12", "delivered_at": "2023-03-07 16:56:37" }
]
```

2) Get shipment by id (GET /shipments/{shipment_id})

```bash
curl "http://127.0.0.1:8001/shipments/1"
```

Sample response (object):

```json
{
	"shipment_id": 1,
	"order_id": 128,
	"carrier": "DHL",
	"status": "DELIVERED",
	"tracking_no": "TRK6188",
	"shipped_at": "2023-09-25 09:07:12",
	"delivered_at": "2023-03-07 16:56:37"
}
```

3) Create a new shipment (POST /shipments)

```bash
curl -X POST "http://127.0.0.1:8001/shipments" \
	-H "Content-Type: application/json" \
	-d '{"shipment_id": 1, "order_id": 123, "carrier": "DHL", "status":"SHIPPED", "tracking_no":"TRACK123" }'
```

Successful response:

```json
{ "message": "✅ Shipment inserted successfully" }
```

4) Update shipment (PUT /shipments/{shipment_id})

```bash
curl -X PUT "http://127.0.0.1:8001/shipments/1" \
	-H "Content-Type: application/json" \
	-d '{"shipment_id":1, "order_id":123, "carrier":"UPS", "status":"IN_TRANSIT", "tracking_no":"NEWTRACK123", "shipped_at":"2025-01-01T12:00:00", "delivered_at": null }'
```

Successful response:

```json
{ "message": "✅ Shipment updated successfully" }
```

5) Delete shipment (DELETE /shipments/{shipment_id})

```bash
curl -X DELETE "http://127.0.0.1:8001/shipments/1"
```

Successful response:

```json
{ "message": "🗑️ Shipment 1 deleted successfully" }
```

## Notes & Caveats

- CSV loading uses pandas and inserts rows directly; very large CSVs may take time.
- Primary keys (`shipment_id`) must be unique. Duplicate keys will produce DB errors.
- The service reads `.env` from the service folder — do not commit real secrets. `.gitignore` already excludes `.env`.
- Error responses include DB error messages to help debugging in development; avoid exposing such details in production.

## Teardown (drop) SQL

If you need to drop the table or the database during development, the SQL snippets are included as commented code inside `db_setup.py` under the comment "OPTIONAL: teardown SQL". Use them only when you intentionally want to remove data.

## Useful files

- `main.py` — API implementation
- `db_utils.py` — connection helper (loads `.env` in the service folder)
- `db_setup.py` — create/load schema and data; contains commented drop SQL
- `csv_files/` — Shipments sample CSV

## License

No license specified. Add a `LICENSE` file if you want to make this open source.

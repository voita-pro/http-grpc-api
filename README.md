# Proto3, gRPC, HTTP, OpenAPI, swagger, PostgreSQL, sqlc, goose

## gRPC, HTTP, OpenAPI template swagger generation from proto3 schemas, 
## PostgreSQL with sqlc generation and goose migrations

A template project for building API services in Go with gRPC, HTTP (grpc-gateway), OpenAPI/Swagger, JWT/Firebase authentication, and PostgreSQL.

Use it as a foundation for your microservices: define protobuf contracts, get a gRPC server and an automatically generated HTTP gateway with documentation.

## Key features
- gRPC server + REST gateway via grpc-gateway
- OpenAPI generation and Swagger UI (in Debug mode)
- Authentication: Firebase Auth or local JWT fallback
- PostgreSQL integration (docker-compose for local development)
- Clean architecture: adapters → usecase → domain → repository
- Built-in example endpoints: auth and a Countries catalog

## Architecture and flow
- gRPC: service `ServicePB` (see `api/proto/service.proto`).
- HTTP: proxied through grpc-gateway → the same methods are available as REST.
- Swagger: OpenAPI static files in `api/openapi` are embedded via `api/openapi.go`; UI available at `/swagger/` when `APP_DEBUG=true`.
- Auth: uses Firebase when configured; otherwise falls back to internal JWT provider (`internal/helpers/jwt`).

Main packages:
- `internal/app` — service bootstrap and HTTP/gRPC servers
- `internal/adapters/grpc` — gRPC handler + server
- `internal/adapters/http` — grpc-gateway + Swagger mounting
- `internal/usecase` — business logic (e.g., `SignIn`, `Countries`, `CountrySave`)
- `internal/repository` — data access (PostgreSQL)
- `internal/config` — configuration via ENV
- `api/proto` — protobuf contracts
- `pkg/pb` — generated gRPC stubs

## Requirements
- Go 1.22+
- make
- Docker + docker-compose (for local dev with DB)

Recommended tooling (codegen):
- buf, protoc plugins, grpc-gateway, protoc-gen-go, protoc-gen-go-grpc, protoc-gen-grpc-gateway
  (see `make install-deps`)

## Quick start
1) Copy this repo as a template and go to the project folder.

2) Create `.env` in the project root (see example below).

3) Run locally via Docker:
```
docker compose up --build
```
Service will be available at:
- HTTP: http://localhost:8989
- gRPC: localhost:5001
- PostgreSQL: localhost:5435 (forwarded to 5432 inside the container)

4) Sanity check:
- Healthcheck: `GET http://localhost:8989/healthz` → `200 ok`
- Swagger UI: `http://localhost:8989/swagger/` (when `APP_DEBUG=true`)

### Run locally without Docker
1) Run PostgreSQL yourself (or use an existing instance).
2) Adjust `.env` for your DB.
3) Install codegen deps (optional): `make install-deps`.
4) Start the service:
```
go run ./cmd/main
```
By default HTTP listens on `:8080`, gRPC on `:5001`.

## Environment variables (.env)
All variables are read with the `APP_` prefix (see `internal/config`). Example:
```
# General
APP_ENV=.env
APP_DEBUG=true
APP_SHUTDOWN_TIMEOUT=5s

# HTTP
APP_HTTP_HOST=0.0.0.0
APP_HTTP_PORT=8080
APP_HTTP_SECRET=1122334455
APP_HTTP_EXP_TOKEN=1h

# gRPC
APP_GRPC_HOST=0.0.0.0
APP_GRPC_PORT=5001
APP_GRPC_SECRET=5544332211

# DB (used by docker-compose)
APP_DB_HOST=application_db
APP_DB_PORT=5432
APP_DB_USER=postgres
APP_DB_PASSWORD=postgres
APP_DB_NAME=app_db
APP_DB_SSLMODE=disable

# Firebase (optional; when empty falls back to JWT)
APP_FIREBASE_PROJECT_ID=
APP_FIREBASE_CREDENTIALS_FILE=/app/firebase.json
APP_FIREBASE_WEB_API_KEY=


# Migrations
GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=./db/migrations
GOOSE_DBSTRING=postgres://${APP_DB_USER}:${APP_DB_PASSWORD}@${APP_DB_HOST}:${APP_DB_PORT}/${APP_DB_NAME}?sslmode=${APP_DB_SSLMODE}

```
Notes:
- For local runs on the host without Docker, set `APP_DB_HOST` to your Postgres (e.g., `localhost`) and port `5435` if you use compose for DB only.
- If `APP_FIREBASE_*` are set correctly, Firebase auth will be used; otherwise the internal JWT provider is used.

## API (example endpoints)
Contracts: `api/proto/service.proto`.
HTTP routes (via grpc-gateway):
- POST `/api/login` — authentication, body:
  ```json
  {"email": "user@example.com", "password": "secret"}
  ```
  Response:
  ```json
  {"token": "<JWT or Firebase custom token>"}
  ```
- GET `/api/countries` — list countries (open in the template)
- POST `/api/country` — add a country (requires Bearer token)
- PUT `/api/country/{id}` — update a country (requires Bearer token)

gRPC methods:
- `Login(LoginIN) returns (LoginOUT)`
- `Countries(google.protobuf.Empty) returns (CountriesOUT)`
- `CountryAdd(CountryIN) returns (Country)`
- `CountrySave(Country) returns (Country)`

### Authentication
- Get a token via `POST /api/login`.
- Pass it in the `Authorization: Bearer <token>` header for protected endpoints.
- In JWT fallback the token is HMAC-signed based on `APP_GRPC_SECRET` and uid; TTL is controlled by `APP_HTTP_EXP_TOKEN` (see `internal/helpers/jwt`).
- In the example `Countries` is open, while `CountryAdd/Save` are protected and use token data.

### Request examples (HTTP)
```
# Login
curl -sS http://localhost:8989/api/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"secret"}'

# Get countries
curl -sS http://localhost:8989/api/countries

# Add country (requires token)
TOKEN="<paste_your_token>"
curl -sS http://localhost:8989/api/country \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Spain","code":"ES","iso_code":"ESP"}'

# Update country
curl -sS -X PUT http://localhost:8989/api/country/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"id":1,"title":"España","code":"ES","iso_code":"ESP"}'
```

### Examples (gRPC)
Using grpcurl:
```
# Login
grpcurl -plaintext localhost:5001 pb.ServicePB/Login \
  -d '{"email":"user@example.com","password":"secret"}'

# Countries
grpcurl -plaintext localhost:5001 pb.ServicePB/Countries
```

## Database migrations
The project includes a structure for migrations (`db/migrations`) and SQL queries (`db/sql`).
Makefile has basic goose targets (assuming goose is installed in your system):
```
make migration-status
make migration-add name=create_countries
make migration-up
make migration-down
```

## Generate the PostgreSQL repository with sqlc
This project uses `sqlc` to generate the data-access layer under `internal/repository/pgdb` from SQL queries in `db/sql` using schemas from `db/migrations`.

- Config file: `sqlc.yaml`
- Engine: PostgreSQL
- Queries source: `db/sql`
- Schema for analysis: `db/migrations`
- Output package: `internal/repository/pgdb` with package name `pgdb`
- Driver: `pgx/v5`
- Extras: interfaces are emitted, JSON and DB tags, empty slices instead of nil, exported queries, pointers for NULL types, and overrides for `numeric` → `shopspring/decimal.Decimal` and `uuid` → `google/uuid`

### Install sqlc
- Go install:
  ```bash
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  ```
  Make sure your `GOBIN` is on `PATH`.
- Or run via Docker (no install needed):
  ```bash
  docker run --rm -v "$PWD":/src -w /src sqlc/sqlc:latest version
  ```

### Generate code
From the project root:
```bash
# Local binary
sqlc generate

# Or via Docker
docker run --rm -v "$PWD":/src -w /src sqlc/sqlc:latest generate
```
The generated files will appear under `internal/repository/pgdb` (e.g., `querier.go`, models, and CRUD methods). The generated interface can be used in your use cases for dependency inversion.

### Adding or changing queries
- Put your `.sql` files in `db/sql`. Use `-- name:` annotations to name queries and return types, for example:
  ```sql
  -- name: GetCountry :one
  SELECT id, title, code, iso_code
  FROM countries
  WHERE id = $1;

  -- name: ListCountries :many
  SELECT id, title, code, iso_code
  FROM countries
  ORDER BY id;

  -- name: InsertCountry :one
  INSERT INTO countries (title, code, iso_code)
  VALUES ($1, $2, $3)
  RETURNING id, title, code, iso_code;
  ```
- Use `:one`, `:many`, `:exec`, or `:execrows` depending on the expected result.
- For nullable columns, sqlc will generate pointer fields when `emit_pointers_for_null_types` is enabled (as in this repo).
- After changes, re-run `sqlc generate`.

## Protobuf, gateway, and Swagger codegen
Install the tools:
```
make install-deps
```
Then run buf (see `buf.yaml`, `buf.gen.yaml`):
```
$(go env GOPATH)/bin/buf generate
```
Generated files go to `pkg/pb` and `api/openapi`.

## Project structure (high level)
- `cmd/main` — entry point
- `internal/app` — initialization and service startup
- `internal/adapters/grpc` — gRPC server and handlers
- `internal/adapters/http` — grpc-gateway and Swagger
- `internal/usecase` — business logic
- `internal/domain` — domain interfaces and models
- `internal/repository` — data access (PostgreSQL)
- `api/proto` — protobuf contracts
- `api/openapi` — OpenAPI/Swagger JSON
- `pkg/pb` — generated gRPC clients/servers

## Useful links
- gRPC Gateway: https://github.com/grpc-ecosystem/grpc-gateway
- buf: https://buf.build/
- Firebase Auth: https://firebase.google.com/docs/auth
- sqlc: https://docs.sqlc.dev/
- pgx: https://github.com/jackc/pgx

## License
See LICENSE.

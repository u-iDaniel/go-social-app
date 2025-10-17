# Go Social App — Backend

A production-ready Go backend for a simple social network: users can register and activate via email, authenticate with JWT, create and manage posts, follow/unfollow, and fetch a personalized feed. It uses PostgreSQL, Redis (optional), SendGrid/Mailtrap for email, and Chi for HTTP routing with Swagger docs.

This README covers backend setup, configuration, running locally and via Docker, migrations/seed, API overview, auth, rate limiting, caching, and testing.


## Tech Stack

- Language: Go 1.24
- HTTP: Chi v5, CORS middleware, expvar debug
- Auth: JWT (github.com/golang-jwt/jwt/v5), Basic Auth for /debug/vars
- Database: PostgreSQL (github.com/lib/pq)
- Cache: Redis v8 (optional) for user fetches
- Email: SendGrid (default) or Mailtrap
- Validation: go-playground/validator
- Docs: Swagger (swaggo)
- Logging: Uber Zap


## Architecture at a Glance

- Entry point: `cmd/api/main.go` initializes config, DB, optional Redis, mailer, JWT, and rate limiter, then mounts routes in `cmd/api/api.go`.
- Domain storage: `internal/store` with concrete stores for users, posts, comments, followers, and roles. Query timeouts and optimistic locking for posts.
- Auth: `internal/auth` and JWT middleware in `cmd/api/middleware.go`.
- Rate limiting: `internal/ratelimiter` (fixed window, in-memory).
- Email: `internal/mailer` with templates in `internal/mailer/templates/`.
- Caching: `internal/store/cache` (Redis) — optional.
- DB helpers and seeding: `internal/db` and `cmd/migrate/seed`.
- OpenAPI: `docs/swagger.yaml` and served under `/v1/swagger/*`.


## Quick Start

Choose Docker or Local setup. By default the API listens on `:8080` and exposes Swagger at `http://localhost:8080/v1/swagger/index.html`.

### 1) Run backing services with Docker Compose

This repo includes a minimal `docker-compose.yml` for PostgreSQL and Redis.

```bash
docker compose up -d
```

Services:
- Postgres: localhost:5431 (db: social, user: admin, password: adminpassword)
- Redis: localhost:6379
- Redis Commander (optional): http://localhost:8081

### 2) Configure environment

Create a `.env` file at the repo root if you want to override defaults used by the app. Common variables:

```env
# Server
ADDR=:8080
ENV=development
EXTERNAL_URL=localhost:8080        # used for Swagger host
FRONTEND_URL=http://localhost:5173 # CORS allowlist and activation link base

# Database
DB_ADDR=postgres://admin:adminpassword@localhost:5431/social?sslmode=disable
DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=30
DB_MAX_IDLE_TIME=15m

# Redis (optional)
REDIS_ENABLED=false
REDIS_ADDR=localhost:6379
REDIS_DB=0
REDIS_PW=

# Auth
AUTH_BASIC_USER=admin
AUTH_BASIC_PASS=adminpassword
AUTH_TOKEN_SECRET=example-secret-key

# Mail
FROM_EMAIL=you@example.com
SENDGRID_API_KEY=your-sendgrid-key
# or Mailtrap
MAILTRAP_API_KEY=your-mailtrap-api-key

# Rate Limiting
RATELIMITER_ENABLED=true
RATELIMITER_REQUESTS_COUNT=20 # per 5s window
```

Notes:
- If `REDIS_ENABLED=true`, the Users cache layer is activated for `GET /users/{userID}`.
- Either `SENDGRID_API_KEY` or `MAILTRAP_API_KEY` should be provided alongside `FROM_EMAIL` for welcome/activation emails.
- `FRONTEND_URL` is used to build the activation URL sent to new users: `${FRONTEND_URL}/confirm/{token}`.

### 3) Run migrations

This project uses the `migrate` CLI. Targets are provided in the `Makefile` and expect `DB_ADDR` to be set (from `.env`).

```bash
make migrate-up              # apply all migrations
# make migrate-down N        # rollback N steps
# make migrate-force VERSION # force set migration version (recovery)
```

Migrations live in `cmd/migrate/migrations/` and include users, posts, comments, roles, followers, and supporting indexes.

### 4) Start the API (local)

```bash
go run ./cmd/api
```

Or build and run the binary:

```bash
go build -o bin/main ./cmd/api
./bin/main
```

Swagger UI: http://localhost:8080/v1/swagger/index.html

### Hot Reload (Air)

Install Air (one-time):

```bash
go install github.com/air-verse/air@latest
```

Run the API with hot reload:

```bash
make dev
```

This uses `.air.toml` to rebuild `./cmd/api` to `./tmp/api` on file changes and restart automatically.

### Optional: Seed Development Data

Seeding creates sample users, posts, and comments.

```bash
make seed
```

The seeder uses the same `DB_ADDR`. Ensure your DB is empty or truncate tables before seeding if needed.


## Docker Image

The `Dockerfile` uses a multi-stage build:
- Builder: golang:1.24-alpine compiles the `cmd/api` package into a static binary
- Runner: `scratch` with CA certs copied for outbound TLS (email providers)

Expose 8080 and run `./api`.

Build and run:

```bash
docker build -t go-social-app:local .
docker run --rm -p 8080:8080 --env-file .env go-social-app:local
```


## API Overview

Base path: `/v1` — full OpenAPI spec in `docs/swagger.yaml` and served at `/v1/swagger/*`.

Auth scheme:
- JWT Bearer tokens via `/v1/authentication/token`
- Add header: `Authorization: Bearer <token>`
- Some endpoints require Basic Auth for admin/debug: `/v1/debug/vars`

Common responses:
- JSON envelope on success: `{ "data": ... }`
- JSON envelope on error: `{ "error": "message" }`

Key endpoints (see Swagger for all details):

- Auth
	- POST `/authentication/user` — Register; sends activation email and returns an invitation token
	- POST `/authentication/token` — Obtain JWT using email/password

- Users
	- PUT `/users/activate/{token}` — Activate account via invitation token
	- GET `/users/{userID}` — Fetch profile (JWT)
	- PUT `/users/{userID}/follow` — Follow user (JWT)
	- PUT `/users/{userID}/unfollow` — Unfollow user (JWT)
	- GET `/users/feed` — Personalized feed with pagination, tags, and search (JWT)

- Posts (JWT required)
	- POST `/posts` — Create
	- GET `/posts/{id}` — Get (includes comments)
	- PATCH `/posts/{id}` — Update (optimistic locking by version)
	- DELETE `/posts/{id}` — Delete

- Ops
	- GET `/health` — Health check
	- GET `/debug/vars` — expvar (Basic Auth)

Roles and permissions:
- Post update: owner or role level ≥ moderator
- Post delete: owner or role level ≥ admin


## Rate Limiting

Fixed-window in-memory limiter with configurable requests per 5s window. Controlled via env:

- `RATELIMITER_ENABLED` (default true)
- `RATELIMITER_REQUESTS_COUNT` (default 20)

When exceeded, returns 429 with `Retry-After` header.


## Caching (Redis, optional)

If enabled, user lookups are cached for 1 minute under keys like `user-{id}`. Configure with `REDIS_*` envs.


## Email Providers

- SendGrid: default client with sandbox toggle for non-production. Requires `SENDGRID_API_KEY` and `FROM_EMAIL`.
- Mailtrap: alternative SMTP client using `MAILTRAP_API_KEY` and `FROM_EMAIL`.

Template: `internal/mailer/templates/user_invitation.tmpl` must include `subject` and `body` blocks. Activation link uses `FRONTEND_URL`.


## Development and Testing

- Run tests:

```bash
make test
```

- Generate Swagger docs (if you update annotations):

```bash
make gen-docs
```

Tips:
- For local dev, run `docker compose up -d` to start DB/Redis first.
- Ensure `.env` includes `DB_ADDR`; the Makefile reads it for migrations and seed.
- If migrations get a dirty state, use `make migrate-force <version>` to recover.


## Troubleshooting

- Can’t connect to Postgres: ensure Docker is running and port 5431 is free, or update `DB_ADDR`.
- Email sending fails: verify `FROM_EMAIL` and API keys. In development, consider Mailtrap or SendGrid sandbox mode.
- 401 Unauthorized: make sure you’re sending `Authorization: Bearer <token>` with a valid JWT.
- CORS issues: update `FRONTEND_URL` to match your frontend dev server URL.
- Swagger host incorrect: set `EXTERNAL_URL` to the externally reachable host:port.

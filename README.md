# Virtus

Backend API for the Virtus fulfillment platform. Built with Go, PostgreSQL, and Redis.

## Prerequisites

- [Go](https://go.dev/dl/) 1.22+
- [Docker](https://www.docker.com/) + Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI
- [sqlc](https://sqlc.dev/) (only needed if you edit SQL queries)
- [Stripe CLI](https://stripe.com/docs/stripe-cli) (only needed for local payment testing)

## Setup

### 1. Clone and install dependencies

```bash
git clone https://github.com/eulerbutcooler/virtus
cd virtus
go mod download
```

### 2. Configure environment

```bash
cp .env.example .env
```

Open `.env` and fill in the required fields:

| Variable | Required | Notes |
|---|---|---|
| `DB_PASSWORD` | ✅ | Password for postgres |
| `JWT_SECRET` | ✅ | Any long random string |
| `STRIPE_SECRET_KEY` | Only for payments | From your Stripe dashboard |
| `STRIPE_WEBHOOK_SECRET` | Only for payments | From `stripe listen` output |

Everything else has sensible defaults and works as-is for local dev.

### 3. Start infrastructure

```bash
docker compose up -d
```

This starts PostgreSQL on port `5432` and Redis on port `6379`.

### 4. Run migrations

```bash
migrate -path ./migrations -database "postgres://virtus:YOUR_DB_PASSWORD@localhost:5432/virtus?sslmode=disable" up
```

Replace `YOUR_DB_PASSWORD` with the value you set in `.env`.

### 5. Run the API server

```bash
go run ./cmd/api
```

Server starts at `http://localhost:8080`. Check it's alive:

```bash
curl http://localhost:8080/healthz
```

### 6. Run the background worker (optional, separate terminal)

```bash
go run ./cmd/worker
```

The worker handles queue funding, delivery watchdog, and impact reminders on a schedule.

---

## Local Stripe webhooks (optional)

If you want to test payments locally, you need to forward Stripe events to your running server:

```bash
stripe listen --forward-to localhost:8080/webhooks/stripe \
  --events payment_intent.succeeded,payment_intent.payment_failed,payment_intent.canceled
```

Copy the webhook signing secret it prints and set it as `STRIPE_WEBHOOK_SECRET` in your `.env`, then restart the server.

---

## Project structure

```
cmd/
  api/        → HTTP server entrypoint
  worker/     → Background jobs entrypoint
internal/
  config/     → Env config loading
  domain/     → Entities, interfaces, sentinel errors
  handler/    → HTTP router, middleware, request handlers
  repository/ → Postgres (sqlc) + Redis implementations
  service/    → Business logic
  worker/     → Periodic task implementations
migrations/   → SQL migration files (up + down)
pkg/          → Shared utilities (crypto, logger, pagination, stripe)
```

---

## Regenerating SQL queries

If you edit any file in `internal/repository/postgres/queries/`, regenerate the Go code:

```bash
sqlc generate
```

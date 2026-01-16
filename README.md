# URL Shortener

![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go)

Service converts long URL into short, shareable link and provides click analytics.

## Features

- User registration and JWT-based authentication
- Creation and management of short links
- Public redirects by short code (`GET /{code}`)
- Asynchronous links click tracking
- Per-link analytics
- Automatic cleanup of expired links via background workers
- PostgreSQL storage with migrations
- Swagger UI

---

## Project structure

```
cmd/app               # main.go
internal/config       # env config
internal/feature/*    # business features (handlers/services/repos)
internal/storage      # postgres repositories
internal/server       # router + middleware
migrations            # SQL migrations
openapi               # OpenAPI specification 
```

## Quick start

### 1) Prerequisites

- Go 1.25+ (or newer)
- PostgreSQL 14+ (or newer)

### 2) Configure env

For local development the service loads config from file `.env` when `APP_ENV=dev`.

![app_env.png](assets/app_env.png)

### 3) Run database migrations

This project uses **golang-migrate**.

Install:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Apply migrations:

```bash
migrate -path migrations -database "postgres://postgres:admin@localhost:5433/url_shortener?sslmode=disable" up
```

---

## API overview

Swagger UI is available at http://localhost:8080/swagger/

![swagger.png](assets/swagger.png)

### Auth

#### Register

`POST /api/v1/auth/register`

Body:

```json
{ "email": "user@example.com", "password": "secret1" }
```

---

#### Login

`POST /api/v1/auth/login`

Body:

```json
{ "email": "user@example.com", "password": "secret1" }
```

**Success (200 OK)**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
---

### Short links

#### Create (auth required)

`POST /api/v1/urls`

Header:

```
Authorization: Bearer <JWT>
```

Body:

```json
{ "long_url": "https://example.com" }
```

**Success (201 Created)**

```json
{
  "code": "kP3sA2",
  "short_url": "https://short.ly/kP3sA2",
  "expires_at": "2026-01-14T00:00:00Z"
}
```
---

#### Redirect (public)

`GET /{code}` → `302 Found`

**Response**

```
Location: https://example.com
```

---

### Analytics

#### Link stats (auth required)

`GET /api/v1/links/{code}/stats?days=30`

Header:

```
Authorization: Bearer <JWT>
```

**Success (200 OK)**

```json
{
  "total_clicks": 1284,
  "unique_clicks": 732,
  "daily": [
    { "date": "2026-01-10", "clicks": 120 },
    { "date": "2026-01-11", "clicks": 240 },
    { "date": "2026-01-12", "clicks": 360 },
    { "date": "2026-01-13", "clicks": 564 }
  ]
}
```

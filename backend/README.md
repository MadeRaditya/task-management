# Backend

REST API Task Management dengan Go + Gin + MySQL + Redis.

## Requirement

- Go 1.21+
- MySQL 8 (Docker disarankan)
- Redis 7 (Docker disarankan)

## Setup

### 1. Jalankan MySQL & Redis

```bash
docker run --name mysql-task -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=taskdb -p 3306:3306 -d mysql:8
docker run --name redis-task -p 6379:6379 -d redis:7
```

Atau dengan Docker Compose dari root project:
```bash
docker compose up -d
```

### 2. Konfigurasi environment

Copy `.env.example` → `.env` dan sesuaikan:
```
APP_PORT=8080
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=taskdb
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 3. Jalankan migration

```powershell
.\migrate.ps1
```

Atau manual:
```bash
Get-Content .\migrations\000001_create_tasks.up.sql | docker exec -i mysql-task mysql -uroot -proot taskdb
```

### 4. Jalankan server

```bash
go mod download
go run ./cmd/server
```

Server akan listening di `http://localhost:8080`.

## Test

```bash
go test ./... -v
```

Test yang tersedia:
- `internal/repository` — search/filter query, soft delete.
- `internal/service` — update, duplicate title (409), cache invalidation.
- `internal/cache` — set/get, cache key per query, invalidation, TTL.

## Arsitektur

```
Handler (HTTP) → Service (business logic) → Repository (SQL)
                     ↓
                  Cache (Redis)
```

- **Handler** — parsing request, formatting response.
- **Service** — validasi, aturan bisnis, cache invalidation.
- **Repository** — query SQL, soft delete.
- **Cache** — Redis dengan versioning untuk invalidation.

## Cache Strategy

Key format:
```
tasks:list:v={version}|status={status}|keyword={keyword}|assignee={assignee}|page={page}|limit={limit}|sort_by={sort_by}|sort_dir={sort_dir}
```

Invalidation dilakukan dengan `INCR tasks:version`. Semua key lama otomatis tidak match dengan versi baru. TTL key = 60 detik.

## Contoh Request

```bash
# Create
curl.exe -X POST http://localhost:8080/api/tasks -H "Content-Type: application/json" -d '{\"title\":\"Belajar Go\"}'

# List dengan filter
curl.exe "http://localhost:8080/api/tasks?status=todo&keyword=go&page=1&limit=10&sort_by=created_at&sort_dir=desc"

# Update
curl.exe -X PUT http://localhost:8080/api/tasks/1 -H "Content-Type: application/json" -d '{\"status\":\"done\"}'

# Delete
curl.exe -X DELETE http://localhost:8080/api/tasks/1
```

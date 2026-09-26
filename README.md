# Task Management App

Fullstack Task Management application built with Go (Gin) + MySQL + Redis on the backend, and React Native (Expo) + TypeScript on the frontend.

Dibuat sebagai bagian dari technical assessment untuk posisi Full Stack Developer.

## Fitur

### Backend
- REST API untuk CRUD task (`GET`, `POST`, `PUT`, `DELETE`).
- Filtering: `status`, `keyword`, `assignee`.
- Pagination: `page`, `limit`.
- Sorting: `sort_by`, `sort_dir`.
- Soft delete (`deleted_at`).
- Redis cache untuk `GET /api/tasks` (TTL 60 detik) dengan cache key yang menyertakan semua query parameter.
- Cache invalidation otomatis setelah create/update/delete (pakai versioning).
- Error response konsisten dengan format `{ "error": { "code", "message" } }`.
- Duplicate title mengembalikan HTTP 409, bukan 500.

### Frontend
- Daftar task dengan search (debounced), filter status, dan pagination.
- Edit modal untuk update task.
- Loading state.
- Refresh otomatis setelah update/hapus.

## Struktur

```
task-management/
├── backend/          # Go (Gin)
│   ├── cmd/server/       # entry point
│   ├── internal/         # config, database, model, repository, service, handler, cache
│   ├── migrations/       # SQL migration
│   └── migrate.ps1       # helper script
├── frontend/         # React Native (Expo)
│   ├── src/              # api, components, hooks, types, utils
│   └── __tests__/        # component tests
└── docker-compose.yml
```

## Menjalankan

Lihat README di masing-masing folder:
- [backend/README.md](backend/README.md)
- [frontend/README.md](frontend/README.md)

## API Endpoints

| Method | Endpoint | Deskripsi |
|---|---|---|
| GET | `/health` | Cek koneksi MySQL + Redis |
| GET | `/api/tasks` | List task (filter, pagination, sort) |
| POST | `/api/tasks` | Buat task baru |
| GET | `/api/tasks/:id` | Detail task |
| PUT | `/api/tasks/:id` | Update task |
| DELETE | `/api/tasks/:id` | Soft delete task |

### Query params untuk `GET /api/tasks`

- `status` — `todo` | `in_progress` | `done`
- `keyword` — cari di title & description
- `assignee` — filter assignee
- `page` — default 1
- `limit` — default 10, max 100
- `sort_by` — `id` | `title` | `status` | `assignee` | `created_at` | `updated_at`
- `sort_dir` — `asc` | `desc`

## Testing

Backend:
```bash
cd backend
go test ./... -v
```

Frontend:
```bash
cd frontend
npm test
```

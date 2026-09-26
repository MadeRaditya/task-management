# Frontend

Task Management UI dengan React Native (Expo) + TypeScript.

## Requirement

- Node.js 20+
- npm

## Setup

### 1. Install dependency

```bash
npm install
```

> File `.npmrc` sudah diset `legacy-peer-deps=true` untuk menghindari konflik peer dependency di stack Expo SDK 57.

### 2. Konfigurasi environment

Copy `.env.example` → `.env`:
```
EXPO_PUBLIC_API_URL=http://localhost:8080
```

### 3. Jalankan

Web (disarankan untuk dev):
```bash
npx expo start --web
```

Atau mobile (scan QR dengan Expo Go):
```bash
npx expo start
```

Buka `http://localhost:8081` di browser.

Pastikan backend sudah jalan di `http://localhost:8080`.

## Test

```bash
npm test
```

Test yang tersedia:
- `SearchBar.test.tsx` — input search memanggil `onChangeText`.
- `TaskList.test.tsx` — render task, tombol edit/hapus, empty state.

## Struktur

```
src/
├── api/           # axios client + tasks endpoints
├── components/    # SearchBar, StatusFilter, Pagination, TaskList, EditTaskModal, LoadingState
├── hooks/         # useTasks (fetch + filter + pagination state)
├── types/         # TypeScript types
└── utils/         # useDebounce
```

## Fitur

- **Search** — input keyword dengan debounce 400ms. Mencari di title & description.
- **Filter status** — chip untuk `Todo`, `In Progress`, `Done`, `Semua`.
- **Pagination** — Prev/Next, `limit` bisa diubah di `App.tsx`.
- **Edit modal** — form untuk update title, description, assignee, status.
- **Loading state** — ActivityIndicator saat fetch.
- **Refresh otomatis** — list refresh setelah update/delete.

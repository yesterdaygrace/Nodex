# crudin

A small **PostgreSQL CRUD API** (Go + Gin + GORM) with a Vue 3 frontend — built as a
learning project. Run it locally, hit the REST endpoints, and read the code: the
handlers are deliberately small and commented, and the test suite locks in the
API contract.

## Stack

| Layer | Tech |
|---|---|
| Backend | Go 1.25, Gin, GORM (Postgres driver) |
| Database | PostgreSQL (local container recommended) |
| Frontend | Vue 3 + Vite + Tailwind, Axios |

## Prerequisites

- Go 1.25+
- Docker (only needed to run Postgres) — or any Postgres reachable on `127.0.0.1:5435`
- Node.js + npm (only for the frontend)

## Quick start (backend)

1. Start PostgreSQL with the credentials the app expects:

   ```bash
   docker run -d --name crudin-postgres \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=test_db \
     -p 5435:5432 \
     postgres:16
   ```

2. Run the server from this directory:

   ```bash
   go run .            # server on http://localhost:3001
   ```

   ```bash
   # or build and run the binary
   go build -o crudin-server .
   ./crudin-server
   ```

3. Verify it is up:

   ```bash
   curl http://localhost:3001/
   # {"message":"Hello World!"}
   ```

**Database configuration** — the connection string is read from the
`DATABASE_URL` environment variable, falling back to the local default above.
Override it without recompiling:

```bash
export DATABASE_URL="host=127.0.0.1 user=postgres password=postgres dbname=test_db port=5435 sslmode=disable"
go run .
```

On startup the app retries the connection for a few seconds (5 attempts × 2s)
so it boots cleanly when Postgres is starting up at the same time, then panics
if the database is still unreachable. The schema is auto-migrated from the
`Post` model.

## Frontend (development mode)

```bash
cd frontend
npm install
npm run dev
```

Vite serves the UI on `http://localhost:5174` and proxies `/api` to the Go
server on `:3001` (see `frontend/vite.config.js`), so the UI and API work
together with no extra setup.

## API reference

Successful responses and error responses use the same envelope:

```json
{ "success": true, "message": "...", "data": ... }
```

The list endpoint additionally returns pagination metadata (`page`, `limit`,
`total`, `total_pages`) as top-level fields.

| Method | Path | Purpose | Success | Errors |
|---|---|---|---|---|
| `GET` | `/` | Hello world / health | `200` | — |
| `GET` | `/api/posts` | List posts (paginated) | `200` | `400` invalid `page`/`limit` |
| `POST` | `/api/posts` | Create a post | `201` | `400` validation |
| `GET` | `/api/posts/:id` | Fetch one post | `200` | `404` not found |
| `PUT` | `/api/posts/:id` | Update a post | `200` | `400` validation, `404` not found |
| `DELETE` | `/api/posts/:id` | Delete a post | `200` | `404` not found |

### List examples

Posts come back **newest-first** (`created_at DESC, id DESC`).

```bash
# defaults: page=1, limit=20
curl "http://localhost:3001/api/posts"

# explicit window
curl "http://localhost:3001/api/posts?page=2&limit=10"
```

```json
{
  "success": true,
  "message": "Lists Data Posts",
  "data": [
    { "id": 2, "created_at": "2026-08-20T01:23:15.294245+07:00", "title": "beta", "content": "c-beta" },
    { "id": 1, "created_at": "2026-08-20T01:23:15.285296+07:00", "title": "alpha", "content": "c-alpha" }
  ],
  "page": 1,
  "limit": 20,
  "total": 2,
  "total_pages": 1
}
```

Query parameters:

| Param | Default | Bounds |
|---|---|---|
| `page` | `1` | `>= 1` |
| `limit` | `20` | `1..100` (capped) |

### Create example

```bash
curl -X POST http://localhost:3001/api/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"Hello","content":"world"}'
```

Validation failures return `400` with per-field errors:

```json
{ "errors": [ { "field": "Title", "message": "This field is required" } ] }
```

### Error shapes

- `404` — record not found:
  ```json
  { "success": false, "message": "Post not found", "data": null }
  ```
- `400` — invalid pagination params:
  ```json
  { "success": false, "message": "invalid page/limit", "data": null }
  ```
- `500` — database failure (e.g. `"Failed to fetch posts"`).

## Model

| Field | Type | Notes |
|---|---|---|
| `id` | int | primary key, auto-increment |
| `created_at` | timestamp | auto-populated by GORM; drives newest-first ordering |
| `title` | string | required |
| `content` | string | required |

## Testing

The test suite is a real integration suite: it runs against the Postgres
container on `:5435` (the same default DSN the app uses), exercising the Gin
router end to end.

```bash
go test ./controllers/ -v
```

**Note:** `TestMain` truncates the `posts` table (`TRUNCATE ... RESTART
IDENTITY`) before and after the run, so results are deterministic and the
table is left clean. Don't run tests against a database you care about.

Covered behavior (9 tests): empty list, invalid pagination params, create,
create-validation, get-by-id not found, update, update not found, delete,
delete not found.

## Project layout

```text
crudin/
├── main.go                        # routing + slog setup
├── controllers/
│   ├── postsController.go         # handlers + pagination
│   └── postsController_test.go    # integration tests
├── models/
│   ├── post.go                    # Post model
│   └── setup.go                   # DB connection (env-driven + retry)
└── frontend/                      # Vue 3 UI
    └── src/
        ├── App.vue                # state + pagination wiring
        ├── components/            # PostForm / PostList / PostItem
        └── services/api.js        # Axios client
```

## Design notes

- **Consistent response envelope** — successes and errors share the
  `{success, message, data}` shape (see `jsonError`), so the frontend can read
  `err.response.data.message` in one place instead of juggling key names.
- **Pagination** — offset-based (`LIMIT/OFFSET`) for simplicity, newest-first.
  The `FindPosts` handler carries a `TODO` to move to *cursor-based* pagination
  (seek by `created_at + id` from the last row) for write-heavy workloads
  where offset shifts rows between page fetches.
- **Config** — connection settings come from `DATABASE_URL` with a local
  default (12-factor style), so nothing is baked into the binary.
- **Startup resilience** — a short retry-with-backoff tolerates Postgres booting
  at the same time as the app.
<!-- auto test 2 2026-08-30T16:51:34Z -->

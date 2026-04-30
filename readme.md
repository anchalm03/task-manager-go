# task-manager-go

A small Task Manager API in Go (Gin + GORM + Postgres + JWT).

## Setup

1. Postgres running locally on `5432` with a user `taskuser` / password `taskpass` and db `taskdb`.
2. Tables already exist (users, projects, tasks, comments) — no migrations run by the app.
3. Copy `.env` (present in repo) and adjust if needed:
   ```
   DB_URL=postgres://taskuser:taskpass@localhost:5432/taskdb?sslmode=disable
   JWT_SECRET=supersecret
   PORT=8080
   ```
4. Run:
   ```bash
   go run .
   ```

## Auth flow

- `POST /register` → create user
- `POST /login` → returns `{ token }`
- Send `Authorization: Bearer <token>` on every `/api/*` request

JWT claims: `sub` = user id, `role`, `exp` = 24h.

## Endpoints

### Public
| Method | Path | Purpose |
|---|---|---|
| GET | `/ping` | health check |
| POST | `/register` | create user (name, email, password) |
| POST | `/login` | returns JWT |

### Authenticated (`/api/*`, Bearer token required)

**User**
| Method | Path | Purpose |
|---|---|---|
| GET | `/api/me` | caller's profile |

**Projects**
| Method | Path | Rule |
|---|---|---|
| POST | `/api/projects` | any user; owner = caller |
| GET | `/api/projects` | projects owned by caller |
| GET | `/api/projects/:id` | owner only |
| PATCH | `/api/projects/:id` | owner only (name, description) |
| DELETE | `/api/projects/:id` | owner only |

**Tasks**
| Method | Path | Rule |
|---|---|---|
| POST | `/api/projects/:id/tasks` | project owner only |
| GET | `/api/projects/:id/tasks` | project owner only |
| GET | `/api/me/tasks` | tasks assigned to caller |
| PATCH | `/api/tasks/:id` | owner or assignee (only owner can reassign) |
| DELETE | `/api/tasks/:id` | project owner only |

**Comments**
| Method | Path | Rule |
|---|---|---|
| POST | `/api/tasks/:id/comments` | project owner or task assignee |
| GET | `/api/tasks/:id/comments` | project owner or task assignee |
| DELETE | `/api/comments/:id` | comment author or project owner |

## Response envelope

All handlers return:
```json
{
  "success": true,
  "errorcode": 0,
  "error": "",
  "data": { ... }
}
```
Error codes: 0 ok, 1 internal, 2 bad request, 3 unauthorized, 4 not found, 5 forbidden, 6 conflict.

## Postman

Import `task-manager.postman_collection.json`. Run **Auth → Login** first — it saves `{{token}}` automatically. Subsequent requests reuse `{{projectId}}`, `{{taskId}}`, `{{commentId}}` via test scripts.

## What's left

- Logout (stateless JWT — frontend discard, or add a blocklist)
- CORS middleware (only needed once a frontend exists)
- A basic frontend (Vite/React) — see old readme notes
- Deployment

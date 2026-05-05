# Task Assignment — Frontend Integration Guide

Audience: frontend engineer building the "assign a task to a user" feature.
Everything below reflects the **current** state of the backend after the recent
changes. Anything marked **NEW** did not exist in the previous revision.

---

## 1. Mental model

- A **Project** has exactly one **owner** (`owner_id` on the project).
- A **Task** belongs to exactly one project and optionally has one **assignee**
  (`assigned_to`, nullable).
- There is no project-members table yet. Any registered user can be an assignee.
- Only the **project owner** can assign or reassign tasks.
- The **assignee** and the **owner** can both edit a task's content
  (title / description / status / due date), but **only the owner can change
  who it is assigned to**.

---

## 2. Endpoints that relate to assignment

All `/api/*` endpoints require the `Authorization: Bearer <jwt>` header.
Responses use the standard envelope:

```json
{
  "success": true,
  "errorcode": 0,
  "error": "",
  "data": { ... }
}
```

`errorcode` values you'll see on this feature:

| code | meaning           | HTTP |
|------|-------------------|------|
| 0    | NoError           | 200  |
| 1    | InternalServerError | 500 |
| 2    | BadRequest        | 400  |
| 3    | Unauthorized      | 401  |
| 4    | NotFound          | 404  |
| 5    | Forbidden         | 403  |

### 2.1 Create a task (optionally pre-assigned) — existed before

`POST /api/projects/:id/tasks`

Who: **project owner only.**

Body:
```json
{
  "title": "Write launch email",
  "description": "...",
  "status": "todo",
  "due_date": "2026-06-01T00:00:00Z",
  "assigned_to": "1c3c9e2a-..."   // optional user UUID
}
```

**NEW:** if `assigned_to` is provided, the backend now verifies that the user
actually exists. You'll get `errorcode: 2` with
`"assignee user does not exist"` if the UUID doesn't resolve to a real user.

Response `data` is a `TaskView`:
```json
{
  "id": "uuid",
  "title": "...",
  "description": "...",
  "status": "todo",
  "due_date": "2026-06-01T00:00:00Z",
  "assigned_to": "uuid-or-null",
  "project_id": "uuid",
  "created_at": "...",
  "updated_at": "..."
}
```

### 2.2 Assign / reassign / unassign a task — **NEW endpoint**

`POST /api/tasks/:id/assign`

Who: **project owner only.** Non-owners get `403 Forbidden` with message
`"only project owner can assign tasks"`.

Body (assign or reassign):
```json
{ "user_id": "1c3c9e2a-..." }
```

Body (unassign — set to no one):
```json
{ "user_id": null }
```

Failure modes:

| Scenario                                 | errorcode | message                         |
|------------------------------------------|-----------|---------------------------------|
| Not logged in                            | 3         | "unauthorized access"           |
| `:id` not a UUID                         | 2         | "invalid task id"               |
| Task not found                           | 4         | "resource not found"            |
| Caller is not the project owner          | 5         | "only project owner can assign tasks" |
| `user_id` is not a UUID                  | 2         | "invalid user_id"               |
| `user_id` is a UUID but user doesn't exist | 2       | "assignee user does not exist"  |
| DB failure                               | 1         | "internal server error"         |

On success the response `data` is the updated `TaskView` (same shape as create).

### 2.3 Update a task's content — existed, but **assignment removed**

`PATCH /api/tasks/:id`

Who: project owner **or** current assignee.

Body — any subset of:
```json
{
  "title": "...",
  "description": "...",
  "status": "in_progress",
  "due_date": "2026-06-05T00:00:00Z"
}
```

**BREAKING CHANGE:** this endpoint no longer accepts `assigned_to`. If your
existing UI was using PATCH to reassign, switch it to
`POST /api/tasks/:id/assign`. The field is simply ignored if sent, which means
it will appear to succeed but nothing will change — watch out for stale
client code.

### 2.4 Read tasks — existed, unchanged

- `GET /api/projects/:id/tasks` — tasks in a project. Owner only (for now).
- `GET /api/me/tasks` — tasks assigned to the caller.
- Both return `TaskView[]` where `assigned_to` is the user UUID or `null`.

### 2.5 What's **missing** and will affect your UI

- **No "list users" endpoint yet.** `GET /api/me` returns only the caller. To
  populate an assignee picker you'd normally hit something like
  `GET /api/users?search=...` — it doesn't exist. Options until it does:
  1. Ask the backend to add a lightweight user search endpoint.
  2. Make the user type/paste an email and add a tiny
     `GET /api/users/lookup?email=...` endpoint.
  3. Short-term hack: collect `assigned_to` UUIDs the frontend has already
     seen (e.g. from existing tasks) and show those as a recent-assignees list.
- **No expansion of `assigned_to` to a user object.** `TaskView.assigned_to` is
  just a UUID. If you want to render "Assigned to: Priya Sharma", the frontend
  currently has to fetch the user separately. There's also no
  `GET /api/users/:id` endpoint yet — another ask for the backend.

---

## 3. Suggested frontend flows

### 3.1 Task detail view
1. Fetch task via whichever list endpoint you're on.
2. If `assigned_to` is null → show "Unassigned" + an "Assign" button
   (only enabled for the project owner).
3. If `assigned_to` is set → show the UUID (or resolved name once user
   lookup exists) + a "Reassign" / "Unassign" button (owner-only).

### 3.2 Assigning
1. Owner clicks "Assign" → picker opens.
2. User picks an assignee → frontend sends
   `POST /api/tasks/:id/assign` with `{ "user_id": "<uuid>" }`.
3. On success, replace the local task with the returned `TaskView`.
4. On `errorcode: 2` + `"assignee user does not exist"`, show a toast — usually
   means a stale / wrong UUID in the picker.
5. On `errorcode: 5`, the caller isn't the owner; disable the button in the UI
   going forward (this is a defensive check — the button shouldn't have been
   clickable).

### 3.3 Unassigning
Same as above, body `{ "user_id": null }`.

### 3.4 Owner detection on the client
There is no `/api/projects/:id/members` and no `is_owner` flag on project
responses today — the frontend infers ownership by comparing
`project.owner_id === me.id` (both available: `owner_id` on the project
response, `id` on `GET /api/me`). Use that to show/hide the assign controls.

---

## 4. Quick diff: what changed in this pass

| Area                              | Before                                       | After                                           |
|-----------------------------------|----------------------------------------------|-------------------------------------------------|
| Create task with `assigned_to`    | Accepted any UUID, even garbage              | Rejects unknown users with `BadRequest`         |
| Reassigning via `PATCH /tasks/:id`| Allowed (owner only)                         | **Removed.** Use the new assign endpoint        |
| Assign endpoint                   | Did not exist                                | `POST /api/tasks/:id/assign` (owner only)       |
| Unassign semantics                | PATCH with `assigned_to: ""`                 | POST assign with `user_id: null`                |

---

## 5. Example curl calls

Assign:
```sh
curl -X POST https://<host>/api/tasks/<task-id>/assign \
  -H "Authorization: Bearer $JWT" \
  -H "Content-Type: application/json" \
  -d '{"user_id":"1c3c9e2a-1111-2222-3333-444455556666"}'
```

Unassign:
```sh
curl -X POST https://<host>/api/tasks/<task-id>/assign \
  -H "Authorization: Bearer $JWT" \
  -H "Content-Type: application/json" \
  -d '{"user_id":null}'
```

Update content (no assignment here anymore):
```sh
curl -X PATCH https://<host>/api/tasks/<task-id> \
  -H "Authorization: Bearer $JWT" \
  -H "Content-Type: application/json" \
  -d '{"status":"in_progress"}'
```

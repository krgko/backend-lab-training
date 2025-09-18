# Copilot instructions for backend-lab

This repository is a small Go HTTP service (Fiber + GORM + SQLite) implementing auth (register/login -> JWT) and a user profile. These notes capture project-specific patterns and commands so an AI agent can become productive quickly.

Before each modification, you need to read `docs` folder first.
- User API - `docs/detail.md`

High-level architecture

- `main.go`: app entrypoint. Calls `database.Init("app.db")`, runs `database.DB.AutoMigrate(...)`, registers routes and starts Fiber.
- `database/`: GORM + sqlite init (`Init(path)` and global `DB`).
- `models/`: GORM models. `models/user.go` embeds `gorm.Model` and now contains profile fields (first_name, last_name, phone, member_code, membership_level, points).
- `controllers/`: HTTP handlers. `auth.go` (register, login), `profile.go` (GET/PUT `/api/profile`).
- `routes/`: route wiring and simple JWT middleware (`AuthRequired`) that parses Bearer tokens and sets `c.Locals("user")`.
- `docs/`: `openapi.json` (spec) and `index.html` (Swagger UI + small helper to preauthorize JWT).

Important code patterns & conventions (project-specific)

- GORM models embed `gorm.Model` (ID, CreatedAt, UpdatedAt). Add new model files in `models/` and include them in AutoMigrate in `main.go`.
- Passwords are bcrypt-hashed (`golang.org/x/crypto/bcrypt`) and model field uses `json:"-"` to avoid leaking hashes.
- JWTs use `github.com/golang-jwt/jwt/v5`. Secret read from `JWT_SECRET` env var; fallback `secret`. Claims include `sub` (user ID) and `email` and an `exp` timestamp.
- Routes: group under `/api`. Auth routes are `/api/auth/*`. Profile endpoints are `/api/profile` and are protected by `AuthRequired` middleware. The middleware expects Authorization: Bearer <token> and loads the DB user into `c.Locals("user")`.
- MemberCode is generated during register (current code: `LBK` + timestamp). MembershipLevel defaults to `Basic` and Points default to 0.

Docs & API surface

- `docs/openapi.json` defines Register, Login, Profile schemas and BearerAuth security scheme. Keep it in sync when endpoints change.
- `docs/index.html` loads Swagger UI and contains a small input to preauthorize a JWT. Serve UI at `/swagger` (static) — see `routes.Setup`.

Developer workflows (commands)

- Run locally: `go run .` (run from repo root after `go mod tidy`).
- Build: `go build -o bin/app .`
- Dependencies: `go mod tidy` or explicit `go get` commands for gorm, sqlite driver, jwt, bcrypt packages.
- DB file: `app.db` created in project root by default. To reset schema during development, edit models and rerun AutoMigrate or delete `app.db` (dev only).

Where to edit for common tasks

- Add a model: create `models/<name>.go` and add the type to the AutoMigrate call in `main.go`.
- Add handler: create function in `controllers/` and wire the route in `routes/routes.go` under the `/api` group.
- Protect endpoints: use `profile.Use(AuthRequired)` or call `AuthRequired` in route setup to require JWTs.
- Update API docs: edit `docs/openapi.json` and optionally `docs/index.html` for UI helpers.

Examples (concrete curl requests)

- Register:
  curl -X POST http://localhost:3000/api/auth/register -H "Content-Type: application/json" -d '{"email":"u@e.com","password":"p"}'
- Login:
  curl -X POST http://localhost:3000/api/auth/login -H "Content-Type: application/json" -d '{"email":"u@e.com","password":"p"}'
- Get profile (replace <token>):
  curl -H "Authorization: Bearer <token>" http://localhost:3000/api/profile
- Update profile:
  curl -X PUT -H "Authorization: Bearer <token>" -H "Content-Type: application/json" -d '{"first_name":"Somchai","last_name":"Jai","phone":"081-234-5678"}' http://localhost:3000/api/profile

Behavioral notes for edits

- Never return `Password` in JSON (keep `json:"-"`).
- Controller errors should map to Fiber status codes (use `fiber.StatusBadRequest`, `fiber.StatusUnauthorized`, etc.).
- Prefer adding API routes under `/api` so versioning can be applied later.
- When changing models that require DB schema changes, update `main.go` AutoMigrate and consider dev migration steps (delete or migrate app.db in dev).

Coding style guidelines

- Formatting:
  - Run `gofmt -w .` and `goimports -w .` before commits. Keep code formatted by the Go toolchain.
  - Run `go vet` and (optionally) `golangci-lint run` in CI.
- Imports: group standard library, third-party, then local project imports. Avoid unused imports.
- Packages and files:
  - Use short, lowercase package names (no underscores). One responsibility per file.
  - Keep `controllers/`, `models/`, `routes/`, and `database/` focused and small.
- Error handling:
  - Return errors and handle them at the boundary (handlers). Map errors to appropriate HTTP status codes using Fiber constants.
  - Avoid panics in request handlers. Use `fmt.Errorf` to wrap context when needed.
- Logging:
  - Use a consistent logger. Include contextual fields (user id, request path) for important events and errors.
- Security & secrets:
  - Read secrets (e.g. `JWT_SECRET`) from environment variables. Never commit secrets to source.
  - Hash passwords with bcrypt and never expose the hash via JSON (`json:"-"`).
  - Validate and sanitize input payloads.
- GORM usage:
  - Embed `gorm.Model` for standard fields. Add indexes for frequently queried columns (e.g. email).
  - Use explicit column tags where needed and prefer readable field names.
- JSON / API:
  - Use snake_case JSON tags to match API conventions (e.g. `first_name`).
  - Keep API docs in `docs/openapi.json` up to date when changing request/response shapes.
- Testing:
  - Write unit tests in `_test.go` files. Use table-driven tests for handlers and model logic.
  - Mock or use an in-memory sqlite DB for tests where appropriate.
- Comments & docs:
  - Add package comments and godoc-style comments for exported functions/types.
  - Keep README and docs updated for developer workflows.
- Git & commits:
  - Use imperative, short commit messages. Run linters and formatters before committing.

If anything is missing or unclear, provide the exact runtime error or the file you want updated and I will adjust these instructions.

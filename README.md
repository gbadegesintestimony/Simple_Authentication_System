# JWT Authentication (Gin + GORM)

This repository provides a minimal JWT-based authentication example using Gin (HTTP), GORM (Postgres), and golang-jwt.

Structure

- `cmd/main.go` — application entry, loads .env, connects DB, migrates models, starts Gin server.
- `database/` — database connection helpers.
- `models/` — GORM models and request/response DTOs.
- `controllers/` — HTTP handlers (auth & user handlers).
- `middleware/` — JWT middleware.
- `routes/` — registers HTTP routes.

Environment
Copy `.env.example` to `.env` and fill in DB credentials and JWT secret.

Run (PowerShell)

```powershell
# ensure go modules are downloaded
go mod tidy

# run the server
go run ./cmd/main.go
```

Quick test (PowerShell)

```powershell
# register
$body = @{ name='Test User'; email='you@example.com'; password='password123' } | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/auth/register -Body $body -ContentType 'application/json'

# login
$body = @{ email='you@example.com'; password='password123' } | ConvertTo-Json
$resp = Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/auth/login -Body $body -ContentType 'application/json'
$token = $resp.token

# change password
$body = @{ current_password='password123'; new_password='newpass456' } | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/change-password -Body $body -ContentType 'application/json' -Headers @{ Authorization = "Bearer $token" }
```

Notes

- Use a secure `JWT_SECRET` in production. Do not commit `.env` with secrets.
- Consider refresh tokens or short lived access tokens for better security.
- Add rate limiting and input validation for production readiness.

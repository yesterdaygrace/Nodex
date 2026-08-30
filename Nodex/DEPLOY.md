# DEPLOY.PRD
## Production Deployment Specification
### Go Application + Dokploy + VPS + Public Live Demo

**Document Version:** 1.0  
**Status:** Production Ready Specification  
**Target:** Public HR-facing live demo  
**Application:** Go Backend / Web Application  
**Deployment Platform:** Dokploy  
**Infrastructure:** Linux VPS  
**Container Runtime:** Docker  
**Reverse Proxy:** Dokploy Proxy / Traefik  
**Database:** PostgreSQL  
**Cache:** Redis, optional  
**Source Control:** GitHub  
**Public Access:** HTTPS + Custom Domain

---

# 1. Deployment Objective

Deploy the Go application as a publicly accessible production-like live demo.

The deployment must allow an HR reviewer to open a single HTTPS URL from the CV without requiring:

- VPS access
- Dokploy access
- SSH
- VPN
- Port numbers
- IP addresses
- Development tooling
- Local environment setup

### Target user experience

```text
CV
 │
 │ Click "Live Demo"
 ▼
https://demo.example.com
 │
 ▼
Dokploy Proxy
 │
 ▼
Go Application
 │
 ├── PostgreSQL
 └── Redis (optional)
```

The deployment should behave like a small production environment rather than simply exposing a development server.

---

# 2. Goals

## 2.1 Primary Goals

1. Deploy the Go application through Dokploy.
2. Automatically deploy from GitHub.
3. Expose the application through a custom HTTPS domain.
4. Keep PostgreSQL private.
5. Keep Redis private if used.
6. Provide health checks.
7. Support automatic restart.
8. Support rollback to a previous deployment.
9. Persist database data.
10. Keep secrets outside Git.
11. Provide production configuration through environment variables.
12. Make the application suitable for HR evaluation.
13. Minimize downtime during redeployment.
14. Provide a clean public API/web URL.

---

# 3. Non-Goals

This deployment is NOT intended to provide:

- Multi-region infrastructure
- Kubernetes
- Horizontal autoscaling
- Enterprise disaster recovery
- Multi-VPS clustering
- Global CDN architecture
- High-frequency financial trading infrastructure
- Unlimited traffic

This is a **portfolio-grade production deployment**, not AWS pretending to be NASA.

---

# 4. Technology Stack

## Application

```text
Language: Go
Runtime: Linux
Architecture: amd64
API: REST
Protocol: HTTP/HTTPS
```

Recommended Go version:

```text
Go 1.24+
```

Use the version actually required by the project.

---

## Infrastructure

```text
VPS
├── Docker
├── Dokploy
├── Dokploy Proxy
├── Go Application
├── PostgreSQL
└── Redis (optional)
```

---

## Source Control

```text
GitHub
```

Repository structure:

```text
repository/
├── cmd/
├── internal/
├── pkg/
├── migrations/
├── tests/
├── Dockerfile
├── .dockerignore
├── go.mod
├── go.sum
├── README.md
└── .env.example
```

---

# 5. Deployment Architecture

## 5.1 High-Level Architecture

```text
                         INTERNET
                             │
                             │ HTTPS :443
                             ▼
                    ┌─────────────────┐
                    │  Custom Domain  │
                    │ demo.example.com│
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Dokploy Proxy   │
                    │ Reverse Proxy   │
                    │ TLS Termination │
                    └────────┬────────┘
                             │
                             │ HTTP
                             ▼
                    ┌─────────────────┐
                    │   Go Service    │
                    │    :8080        │
                    └───────┬─────────┘
                            │
                    ┌───────┴────────┐
                    │                │
                    ▼                ▼
             ┌─────────────┐   ┌─────────────┐
             │ PostgreSQL  │   │    Redis    │
             │   Private   │   │   Private   │
             └─────────────┘   └─────────────┘
```

---

# 6. Network Requirements

The Go application MUST listen on:

```text
0.0.0.0:8080
```

Do NOT bind the application only to:

```text
127.0.0.1:8080
```

because the application is running inside a container.

Example:

```go
server := &http.Server{
    Addr: ":8080",
}
```

The application container exposes:

```text
8080
```

Dokploy Proxy handles public access.

---

# 7. Port Strategy

Only the reverse proxy should be publicly accessible.

### Public

```text
80/tcp
443/tcp
```

### Internal

```text
8080    Go application
5432    PostgreSQL
6379    Redis
```

Do NOT expose:

```text
0.0.0.0:5432
0.0.0.0:6379
0.0.0.0:8080
```

unless there is a specific infrastructure reason.

The public request path should be:

```text
Internet
   ↓
443
   ↓
Dokploy Proxy
   ↓
Go container :8080
```

---

# 8. Domain Configuration

Use a dedicated subdomain.

Example:

```text
demo.example.com
```

DNS:

```text
Type: A
Name: demo
Value: <VPS_PUBLIC_IP>
TTL: 300
```

Result:

```text
demo.example.com
        ↓
     VPS IP
        ↓
Dokploy Proxy
        ↓
Go Application
```

---

# 9. HTTPS Requirements

The application MUST be publicly accessible through:

```text
https://demo.example.com
```

HTTP requests should redirect to HTTPS.

Required:

```text
http://demo.example.com
        ↓
301/308
        ↓
https://demo.example.com
```

TLS certificate should be managed through the Dokploy reverse proxy / supported certificate mechanism.

Never put TLS private keys into GitHub.

---

# 10. GitHub Deployment Flow

Deployment source:

```text
GitHub
```

Recommended branch:

```text
main
```

Flow:

```text
Developer
    │
    ▼
git push origin main
    │
    ▼
GitHub
    │
    ▼
Dokploy Webhook
    │
    ▼
Build Docker Image
    │
    ▼
Run Health Check
    │
    ▼
Deploy Container
    │
    ▼
Public Application
```

---

# 11. Dokploy Application Configuration

Create:

```text
Project:
go-live-demo
```

Application:

```text
go-api
```

Repository:

```text
GitHub repository
```

Branch:

```text
main
```

Build method:

```text
Dockerfile
```

Container port:

```text
8080
```

Protocol:

```text
HTTP
```

Domain:

```text
demo.example.com
```

---

# 12. Dockerfile Requirements

Use a multi-stage Docker build.

Recommended structure:

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o /app/server ./cmd/server

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/server /app/server

USER app

EXPOSE 8080

ENTRYPOINT ["/app/server"]
```

The exact Go version and Alpine version should match the application's compatibility requirements.

---

# 13. Docker Image Requirements

The production image MUST:

- Use a multi-stage build.
- Exclude source build dependencies from runtime.
- Run as a non-root user.
- Expose only the required application port.
- Avoid storing secrets.
- Avoid development tools.
- Avoid unnecessary packages.

Target:

```text
Build Image
    ↓
Go compiler
    ↓
Static binary
    ↓
Minimal runtime image
```

---

# 14. .dockerignore

Required:

```text
.git
.github
.env
.env.*
README.md
Dockerfile
docker-compose.yml
tmp/
logs/
coverage/
*.log
```

Do NOT exclude files required by the application at runtime.

---

# 15. Environment Variables

Secrets MUST NOT be committed to GitHub.

Production variables should be configured in Dokploy.

Example:

```env
APP_ENV=production
APP_PORT=8080

DATABASE_URL=postgres://...
REDIS_URL=redis://...

JWT_SECRET=...

CORS_ALLOWED_ORIGINS=https://demo.example.com
```

---

# 16. Required Configuration Categories

## Application

```env
APP_ENV
APP_PORT
APP_NAME
APP_URL
```

## Database

```env
DATABASE_URL
DB_HOST
DB_PORT
DB_NAME
DB_USER
DB_PASSWORD
```

## Authentication

```env
JWT_SECRET
JWT_EXPIRATION
```

## CORS

```env
CORS_ALLOWED_ORIGINS
```

## Redis

```env
REDIS_URL
```

Redis variables are only required if Redis is used.

---

# 17. Production Configuration Rules

Production configuration MUST:

- Disable debug mode.
- Use production logging.
- Use secure cookies.
- Restrict CORS.
- Use strong secrets.
- Disable unnecessary development endpoints.
- Disable verbose stack traces.
- Avoid returning database errors directly to users.

Never:

```text
DEBUG=true
```

in production.

---

# 18. PostgreSQL Deployment

PostgreSQL should run as a separate Dokploy service.

Architecture:

```text
Dokploy Project
│
├── go-api
│
└── postgres
```

The Go application connects using the internal service/network address.

Example:

```text
go-api
   │
   │ TCP 5432
   ▼
postgres
```

PostgreSQL should NOT be publicly exposed.

---

# 19. PostgreSQL Persistence

The database MUST use persistent storage.

Example:

```text
PostgreSQL Container
       │
       ▼
Persistent Volume
       │
       ▼
VPS Disk
```

Without persistent storage:

```text
Container deleted
       ↓
Database deleted
       ↓
HR discovers that your demo has achieved enlightenment
```

Persistent volume is mandatory.

---

# 20. Database Migration

Application startup should NOT blindly perform destructive migrations.

Preferred workflow:

```text
Deploy
   ↓
Run migration
   ↓
Verify migration
   ↓
Start application
```

Migration command example:

```bash
./server migrate
```

or:

```bash
go run ./cmd/migrate
```

depending on the project architecture.

Recommended separate migration binary:

```text
cmd/
├── server/
│   └── main.go
│
└── migrate/
    └── main.go
```

---

# 21. Health Check

The application MUST expose:

```text
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

HTTP status:

```text
200 OK
```

---

# 22. Readiness Check

Recommended:

```text
GET /ready
```

This should verify required dependencies.

Example:

```text
GET /ready

200 OK
```

when:

```text
Go application ✓
PostgreSQL ✓
Redis ✓
```

If PostgreSQL is unavailable:

```text
503 Service Unavailable
```

---

# 23. Health Check Docker Configuration

Recommended:

```dockerfile
HEALTHCHECK --interval=30s \
            --timeout=5s \
            --start-period=10s \
            --retries=3 \
            CMD wget -qO- http://127.0.0.1:8080/health || exit 1
```

If `wget` is not available in the runtime image, use an alternative strategy or install the minimal required package.

---

# 24. Graceful Shutdown

The Go server MUST handle:

```text
SIGTERM
SIGINT
```

Graceful shutdown sequence:

```text
SIGTERM
   ↓
Stop accepting requests
   ↓
Finish active requests
   ↓
Close Redis
   ↓
Close PostgreSQL
   ↓
Exit
```

Recommended shutdown timeout:

```text
10 seconds
```

Example concept:

```go
ctx, cancel := signal.NotifyContext(
    context.Background(),
    syscall.SIGINT,
    syscall.SIGTERM,
)
defer cancel()
```

---

# 25. HTTP Server Requirements

Production server should configure:

```text
ReadTimeout
WriteTimeout
IdleTimeout
ReadHeaderTimeout
```

Example:

```text
ReadHeaderTimeout: 5s
ReadTimeout:       15s
WriteTimeout:      15s
IdleTimeout:       60s
```

Do not use an unlimited timeout for public HTTP services.

---

# 26. Logging

Application logs should be sent to stdout/stderr.

Preferred:

```text
stdout
stderr
```

NOT:

```text
/app/logs/application.log
```

inside the container.

Reason:

```text
Docker
  ↓
Dokploy
  ↓
Container logs
```

makes logs easier to inspect.

Recommended structured logging:

```json
{
  "level": "info",
  "message": "request completed",
  "method": "GET",
  "path": "/api/products",
  "status": 200,
  "duration_ms": 14
}
```

---

# 27. Logging Requirements

Every HTTP request should ideally provide:

```text
Request ID
Method
Path
Status
Duration
```

Errors should provide:

```text
Request ID
Error type
Internal context
```

Never log:

```text
Passwords
JWT secrets
Database passwords
API keys
Session tokens
```

---

# 28. Security Requirements

## Mandatory

- HTTPS.
- Strong application secrets.
- Non-root container.
- Private PostgreSQL.
- Private Redis.
- Restricted CORS.
- Secure cookies if authentication is used.
- Input validation.
- SQL parameterization.
- Rate limiting for sensitive endpoints.
- No production credentials in Git.

---

# 29. CORS

Do not use:

```text
Access-Control-Allow-Origin: *
```

for an authenticated production API unless explicitly required.

Prefer:

```env
CORS_ALLOWED_ORIGINS=https://demo.example.com
```

If the frontend is hosted separately:

```env
CORS_ALLOWED_ORIGINS=https://demo.example.com,https://www.example.com
```

---

# 30. Authentication

If the project requires authentication:

```text
POST /api/auth/login
POST /api/auth/register
POST /api/auth/refresh
POST /api/auth/logout
```

Passwords MUST be hashed using a password hashing algorithm such as:

```text
Argon2id
```

or:

```text
bcrypt
```

Never store plaintext passwords.

---

# 31. Demo Account

For HR usability, provide a dedicated demo account.

Example:

```text
Email:
demo@example.com

Password:
<generated-demo-password>
```

Do NOT use:

```text
admin/admin
```

Do NOT use production administrative credentials.

Demo account should have limited permissions.

Recommended:

```text
Role:
Demo User

Permissions:
READ
CREATE
UPDATE

No:
DELETE critical data
Manage users
Change system configuration
```

---

# 32. Demo Data

Production demo should contain realistic seed data.

Example:

```text
Users
Products
Categories
Orders
Customers
Transactions
```

depending on project domain.

Do not use real personal information.

Use:

```text
John Doe
demo@example.com
Demo Company
```

instead of actual people's data.

---

# 33. Demo Data Reset

Recommended endpoint:

```text
POST /admin/demo/reset
```

BUT this endpoint must NOT be publicly accessible without strong authorization.

Alternative preferred approach:

```text
Scheduled database reset
```

or manual Dokploy operation.

---

# 34. Deployment Strategy

Preferred:

```text
Build
 ↓
Deploy new container
 ↓
Health check
 ↓
Route traffic
 ↓
Stop old container
```

If Dokploy configuration supports zero-downtime/rolling behavior, use it.

Otherwise:

```text
Stop old
 ↓
Start new
 ↓
Health check
```

is acceptable for a portfolio demo.

---

# 35. Deployment Trigger

Preferred:

```text
Git push → main
```

Deployment:

```text
Developer
   ↓
git push
   ↓
GitHub
   ↓
Webhook
   ↓
Dokploy
   ↓
Build
   ↓
Deploy
```

Avoid manually SSH-ing into the server every time.

The whole point of deployment automation is to prevent you from becoming the deployment pipeline.

---

# 36. CI/CD Validation

Before Dokploy deployment:

```text
go fmt
go vet
go test
go build
```

Recommended GitHub Actions:

```text
Push
 ↓
go fmt check
 ↓
go vet
 ↓
go test
 ↓
Docker build
 ↓
Deploy Dokploy
```

Deployment should ideally occur only when tests pass.

---

# 37. GitHub Actions Example

Recommended workflow:

```text
.github/
└── workflows/
    └── ci.yml
```

Pipeline:

```text
checkout
   ↓
setup Go
   ↓
download dependencies
   ↓
go fmt
   ↓
go vet
   ↓
go test ./...
   ↓
docker build
```

Dokploy can then handle the actual deployment.

---

# 38. Database Backup

Minimum requirement:

```text
Daily PostgreSQL backup
```

Recommended retention:

```text
7 daily backups
```

Better:

```text
7 daily
4 weekly
3 monthly
```

For a portfolio project, 7 daily backups is sufficient.

Backups should NOT live only inside the same PostgreSQL container.

---

# 39. Disaster Recovery

If VPS fails:

```text
New VPS
   ↓
Install Dokploy
   ↓
Deploy application
   ↓
Restore PostgreSQL backup
   ↓
Configure DNS
   ↓
HTTPS
```

Recovery target:

```text
RTO: < 4 hours
RPO: < 24 hours
```

These are reasonable portfolio targets.

Do not advertise these as enterprise guarantees.

---

# 40. VPS Minimum Specification

Recommended starting point:

```text
CPU: 2 vCPU
RAM: 4 GB
Storage: 40+ GB SSD
OS: Ubuntu 24.04 LTS
Docker: Latest supported stable
```

For:

```text
Go API
PostgreSQL
Redis
Dokploy
```

4 GB RAM provides a much more comfortable starting point than trying to make everything survive on 1 GB through sheer optimism.

---

# 41. VPS Firewall

Allow:

```text
22/tcp
80/tcp
443/tcp
```

SSH:

```text
22/tcp
```

should preferably be restricted by IP or protected with SSH keys.

Do not expose:

```text
5432
6379
8080
```

publicly.

---

# 42. SSH Security

Use:

```text
SSH key authentication
```

Disable password authentication when practical.

Disable root login when practical:

```text
PermitRootLogin no
```

Use a dedicated sudo user.

---

# 43. Resource Limits

Application should have reasonable resource limits.

Example:

```text
CPU:
1 CPU

Memory:
512 MB - 1 GB
```

PostgreSQL:

```text
Memory:
512 MB - 1 GB
```

Actual limits should depend on VPS capacity and workload.

---

# 44. Restart Policy

The Go container should restart automatically.

Expected behavior:

```text
Application crash
      ↓
Docker restart
      ↓
Health check
      ↓
Application available
```

Recommended:

```text
unless-stopped
```

or equivalent Dokploy configuration.

---

# 45. Observability

Minimum:

```text
Application logs
Health check
Container status
CPU usage
RAM usage
Disk usage
```

Recommended future stack:

```text
Prometheus
Grafana
OpenTelemetry
```

Do not install an entire observability zoo just to impress an HR person. Add it when the project actually benefits from it.

---

# 46. Error Handling

Public API errors should return structured responses.

Example:

```json
{
  "error": {
    "code": "PRODUCT_NOT_FOUND",
    "message": "Product not found",
    "request_id": "req_123"
  }
}
```

Do not return:

```json
{
  "error": "pq: password authentication failed for user postgres"
}
```

Internal infrastructure errors belong in logs, not in the browser.

---

# 47. API Documentation

Expose API documentation if relevant.

Recommended:

```text
/api/docs
```

or:

```text
/docs
```

Use:

```text
OpenAPI
Swagger UI
```

Documentation should show:

- Authentication
- Endpoints
- Request body
- Response body
- Error codes
- Example requests

This is particularly valuable for an HR/technical reviewer.

---

# 48. Public Routes

Recommended public structure:

```text
/
├── /health
├── /ready
├── /api
│   ├── /auth
│   ├── /users
│   ├── /products
│   └── /orders
└── /docs
```

Health endpoints should be lightweight.

---

# 49. Frontend + Go Backend

If the project has a separate frontend:

```text
                    INTERNET
                       │
          ┌────────────┴────────────┐
          ▼                         ▼
   frontend.example.com       api.example.com
          │                         │
          ▼                         ▼
      Frontend                  Go API
                                    │
                         ┌──────────┴──────────┐
                         ▼                     ▼
                    PostgreSQL              Redis
```

Example:

```text
Frontend:
https://demo.example.com

API:
https://api.demo.example.com
```

CORS:

```env
CORS_ALLOWED_ORIGINS=https://demo.example.com
```

---

# 50. Single-Domain Alternative

For simpler projects, prefer:

```text
https://demo.example.com
```

with:

```text
/demo
/api
```

Example:

```text
https://demo.example.com
https://demo.example.com/api/products
https://demo.example.com/docs
```

This avoids unnecessary CORS complexity.

For a portfolio project, **single-domain architecture is preferable unless there is a reason to separate frontend and backend.**

---

# 51. Deployment Environment

Use:

```text
Development
    ↓
GitHub
    ↓
Production Demo
```

Minimum environment separation:

```text
.env.local
.env.example
Production environment variables in Dokploy
```

Never commit:

```text
.env
```

---

# 52. Production Secrets

Secrets must be generated securely.

Examples:

```text
JWT_SECRET
DATABASE_PASSWORD
SESSION_SECRET
API_KEY
```

Use long random values.

Minimum recommendation:

```text
32+ bytes
```

Do not use:

```text
secret
password
123456
myprojectsecret
```

---

# 53. Database Connection Pool

Configure PostgreSQL pool limits.

Example:

```text
MaxOpenConns: 25
MaxIdleConns: 10
ConnMaxLifetime: 30m
```

Actual values should be adjusted based on workload and PostgreSQL capacity.

For a small HR demo, excessive connection counts are pointless.

---

# 54. Redis Requirements

Redis is optional.

Use Redis for:

```text
Caching
Sessions
Rate limiting
Queues
Temporary state
```

Do NOT add Redis simply because the architecture diagram looks more impressive.

If not needed:

```text
Go
 ↓
PostgreSQL
```

is perfectly valid.

---

# 55. Rate Limiting

Recommended for:

```text
POST /api/auth/login
POST /api/auth/register
```

Example:

```text
5-10 attempts/minute/IP
```

For general API traffic:

```text
60-120 requests/minute/IP
```

Actual values depend on application behavior.

---

# 56. Security Headers

Recommended:

```text
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: strict-origin-when-cross-origin
Content-Security-Policy: appropriate policy
```

If using a frontend, configure CSP according to its actual assets.

Do not blindly paste an internet security-header checklist and break your application.

---

# 57. Application Timeout

External requests should have explicit timeouts.

Example:

```text
HTTP client timeout:
10 seconds
```

Database:

```text
Query timeout:
5-10 seconds
```

External API:

```text
5-15 seconds
```

Never allow an external dependency to hang forever.

---

# 58. Deployment Verification

After deployment, verify:

```text
[ ] Domain resolves
[ ] HTTPS works
[ ] HTTP redirects
[ ] /health returns 200
[ ] /ready returns 200
[ ] Frontend loads
[ ] API responds
[ ] PostgreSQL connection works
[ ] Redis connection works if used
[ ] Authentication works
[ ] Demo account works
[ ] Static assets load
[ ] No CORS errors
[ ] Logs contain no secrets
[ ] Container remains healthy
```

---

# 59. Smoke Test

Run:

```bash
curl -I https://demo.example.com
```

Expected:

```text
HTTP/2 200
```

Health:

```bash
curl https://demo.example.com/health
```

Expected:

```json
{
  "status": "ok"
}
```

API:

```bash
curl https://demo.example.com/api/products
```

Expected:

```text
HTTP 200
```

or an appropriate authenticated response.

---

# 60. Deployment Acceptance Criteria

Deployment is considered successful when ALL conditions are satisfied.

### Accessibility

```text
https://demo.example.com
```

opens from:

```text
Chrome
Firefox
Edge
Mobile browser
```

without VPN or special configuration.

### HTTPS

```text
Valid TLS certificate
```

### Application

```text
Application responds < 3 seconds
```

for normal demo operations.

### Database

```text
PostgreSQL persistent
```

### Deployment

```text
Git push → Dokploy deployment
```

### Reliability

```text
Container automatically restarts after failure.
```

### Security

```text
Database is not publicly accessible.
```

### HR usability

A reviewer can understand the project within:

```text
30 seconds
```

and reach the primary feature within:

```text
60 seconds
```

---

# 61. HR Demo UX Requirements

The public demo is not merely infrastructure.

The landing page should communicate:

```text
PROJECT NAME

One-sentence description.

[ Launch Demo ]

Technology:
Go · PostgreSQL · Docker · Dokploy

Features:
✓ Authentication
✓ REST API
✓ PostgreSQL
✓ Role-based access
✓ Docker deployment
```

If login is required:

```text
Demo Account

Email:
demo@example.com

Password:
********

[ Login ]
```

Provide an obvious way to return to the project overview.

---

# 62. Demo Data Requirements

The demo must never open to:

```text
Empty database
No users
No products
No orders
No content
```

Seed meaningful data.

Example:

```text
12 products
5 categories
8 customers
20 orders
3 demo users
```

Use enough data to demonstrate pagination, filtering, sorting, relationships, and dashboard behavior.

---

# 63. Performance Target

For a small portfolio deployment:

```text
Health endpoint:
< 100 ms

Simple API request:
< 300 ms

Normal database query:
< 500 ms

Initial page:
< 3 seconds
```

These are targets, not contractual guarantees.

---

# 64. Rollback

Every production deployment should have a rollback path.

Rollback:

```text
Current deployment
       ↓
Problem detected
       ↓
Select previous working version
       ↓
Redeploy
       ↓
Health check
       ↓
Restore service
```

Git tags are recommended:

```text
v1.0.0
v1.1.0
v1.2.0
```

Avoid deploying only from arbitrary untracked commits.

---

# 65. Release Strategy

Recommended:

```text
feature/*
    ↓
Pull Request
    ↓
main
    ↓
CI
    ↓
Docker build
    ↓
Dokploy
    ↓
Production Demo
```

For significant changes:

```text
v1.0.0
v1.1.0
v1.2.0
```

---

# 66. Monitoring Checklist

Daily/periodic check:

```text
[ ] Application online
[ ] HTTPS certificate valid
[ ] Disk usage acceptable
[ ] RAM usage acceptable
[ ] CPU acceptable
[ ] PostgreSQL healthy
[ ] Backups working
[ ] No repeated application crashes
```

---

# 67. Failure Scenarios

## Application crashes

Expected:

```text
Docker
 ↓
Restart
 ↓
Health check
 ↓
Available
```

## PostgreSQL crashes

Expected:

```text
Application
 ↓
503 / degraded state
 ↓
PostgreSQL restart
 ↓
Application reconnects
```

## Bad deployment

Expected:

```text
Detect
 ↓
Rollback
 ↓
Previous version
```

## VPS failure

Expected:

```text
New VPS
 ↓
Dokploy
 ↓
Application
 ↓
Database backup restore
 ↓
DNS
```

---

# 68. Security Acceptance Criteria

```text
[ ] No secrets in Git
[ ] .env excluded
[ ] HTTPS enabled
[ ] PostgreSQL private
[ ] Redis private
[ ] Container non-root
[ ] CORS restricted
[ ] Password hashing enabled
[ ] SQL injection protected
[ ] Authentication protected
[ ] Sensitive logs removed
[ ] Production debug disabled
[ ] SSH key authentication
```

---

# 69. Final Production Architecture

```text
                              INTERNET
                                  │
                                  │
                         HTTPS :443 / HTTP :80
                                  │
                                  ▼
                    ┌─────────────────────────┐
                    │      DNS / Domain       │
                    │   demo.example.com      │
                    └────────────┬────────────┘
                                 │
                                 ▼
                    ┌─────────────────────────┐
                    │     Dokploy Proxy       │
                    │   Reverse Proxy + TLS    │
                    └────────────┬────────────┘
                                 │
                                 ▼
                    ┌─────────────────────────┐
                    │       Go API             │
                    │       :8080              │
                    │                          │
                    │  REST API                │
                    │  Auth                    │
                    │  Validation              │
                    │  Business Logic          │
                    └────────────┬────────────┘
                                 │
                       ┌─────────┴─────────┐
                       │                   │
                       ▼                   ▼
              ┌────────────────┐   ┌────────────────┐
              │   PostgreSQL   │   │     Redis      │
              │    :5432       │   │     :6379      │
              │    PRIVATE     │   │    PRIVATE     │
              └───────┬────────┘   └────────────────┘
                      │
                      ▼
               Persistent Volume
                      │
                      ▼
                 VPS Storage
```

---

# 70. Final Deployment Flow

```text
                 DEVELOPER
                     │
                     ▼
                  GitHub
                     │
                     │ Push main
                     ▼
                  CI/CD
                     │
              ┌──────┴──────┐
              │             │
           go test      Docker build
              │             │
              └──────┬──────┘
                     │
                     ▼
                  Dokploy
                     │
                     ▼
              Build Container
                     │
                     ▼
              Start Go Service
                     │
                     ▼
                 /health
                     │
                ┌────┴────┐
                │         │
               FAIL      PASS
                │         │
                ▼         ▼
             Rollback   Proxy Traffic
                          │
                          ▼
                   HTTPS Domain
                          │
                          ▼
                         HR
```

---

# 71. Definition of Done

The deployment is COMPLETE when:

```text
✓ GitHub repository connected
✓ Dokploy project created
✓ Go application deployed
✓ Docker multi-stage build working
✓ Application listens on 0.0.0.0:8080
✓ PostgreSQL deployed
✓ PostgreSQL persistent volume configured
✓ Redis deployed if required
✓ Environment variables configured
✓ Secrets stored outside Git
✓ Custom domain configured
✓ DNS points to VPS
✓ HTTPS certificate active
✓ HTTP → HTTPS redirect working
✓ /health working
✓ /ready working
✓ Graceful shutdown implemented
✓ Container restart configured
✓ Logs accessible
✓ Database migrations working
✓ Demo data seeded
✓ Demo account available
✓ CI tests passing
✓ GitHub → Dokploy deployment working
✓ Rollback tested
✓ Backup configured
✓ Firewall configured
✓ Database not publicly exposed
✓ HR can access the application without infrastructure knowledge
```

---

# 72. Recommended Final CV Presentation

Use:

```text
GO ORDER PLATFORM
Production-style order management platform

Go · PostgreSQL · Redis · Docker · Dokploy

[ GitHub ] [ Live Demo ]

Live Demo:
https://orders.example.com
```

Under the project:

```text
• Built a REST API in Go with PostgreSQL and Redis for order,
  customer, and inventory management.

• Containerized the application with Docker and implemented
  production deployment through Dokploy with HTTPS, health
  checks, persistent PostgreSQL storage, and automated
  GitHub-based deployments.

• Implemented authentication, validation, structured logging,
  graceful shutdown, database migrations, and API documentation.
```

The important distinction is that **Dokploy itself is not the impressive part**. The impressive part is demonstrating that you understand the entire path:

```text
Code
 → Git
 → CI
 → Docker
 → Deployment
 → Reverse Proxy
 → HTTPS
 → Database
 → Health Checks
 → Logs
 → Backups
 → Rollback
 → Public Demo
```

That is what turns “I made a Go CRUD project” into **“I can deploy and operate a Go service.”**
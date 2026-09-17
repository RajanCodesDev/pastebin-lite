# Pastebin Lite

A small full-stack Pastebin application built to demonstrate an end-to-end **CI/CD and cloud-native delivery workflow**.

The application itself is intentionally simple. The focus of this project is how software moves from source code to a tested, containerized, deployed, and verified production workload.

## Architecture

```text
                    ┌──────────────┐
                    │    Browser   │
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │    React     │
                    │   Frontend   │
                    └──────┬───────┘
                           │ HTTP
                           ▼
                    ┌──────────────┐
                    │   Go API     │
                    │   Backend    │
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │ PostgreSQL   │
                    └──────────────┘
```

## Application

Pastebin Lite supports:

* Create a text snippet
* Generate a unique URL for the snippet
* Retrieve a snippet by slug
* Delete a snippet
* Health checking through `/health`

### API

| Method   | Endpoint              | Description              |
| -------- | --------------------- | ------------------------ |
| `GET`    | `/health`             | Application health check |
| `POST`   | `/api/snippets`       | Create a snippet         |
| `GET`    | `/api/snippets/:slug` | Retrieve a snippet       |
| `DELETE` | `/api/snippets/:slug` | Delete a snippet         |

## Technology

### Application

* Go
* `net/http`
* PostgreSQL
* pgx/v5
* React
* Vite
* TypeScript

### Testing

* Go `testing`
* `httptest`
* React Testing Library
* Vitest
* Go race detector
* `go vet`

### Delivery

The project is being built around an end-to-end delivery pipeline:

```text
Git Push / Pull Request
        │
        ▼
       CI
        │
        ├── Backend tests
        ├── Frontend tests
        ├── Static analysis
        ├── Build verification
        └── Security checks
        │
        ▼
   Docker Build
        │
        ▼
 Container Registry
        │
        ▼
       CD
        │
        ▼
   Kubernetes
        │
        ▼
 Deployment Verification
        │
        ▼
 Production
```

> The delivery pipeline is intentionally being implemented incrementally to demonstrate the underlying concepts rather than hiding them behind a managed deployment platform.

## Repository Structure

```text
.
├── backend/
│   ├── cmd/
│   │   └── server/
│   ├── internal/
│   │   ├── config/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── slug/
│   │   └── validation/
│   ├── migrations/
│   ├── go.mod
│   └── go.sum
│
└── frontend/
    ├── src/
    ├── package.json
    └── ...
```

## Backend Design

The backend follows a simple request flow:

```text
HTTP Request
     │
     ▼
 HTTP Handler
     │
     ▼
   Service
     │
     ├── Validation
     ├── Slug Generation
     │
     ▼
 Repository
     │
     ▼
 PostgreSQL
```

The application keeps HTTP concerns, business logic, and persistence separate so each layer can be tested independently.

## Testing

The project includes multiple levels of automated testing.

### Backend

```bash
cd backend

go test ./...
go vet ./...
go test -race ./...
```

Backend tests cover:

* domain validation
* slug generation
* service behavior
* HTTP handlers
* repository behavior
* PostgreSQL integration

### Frontend

```bash
cd frontend

npm test
npm run build
```

Frontend tests cover user-facing behavior such as:

* snippet creation
* validation errors
* API failures
* snippet retrieval
* not-found handling
* deletion

## Local Development

### Backend

```bash
cd backend
go run ./cmd/server
```

The API runs on:

```text
http://localhost:8080
```

Health check:

```bash
curl http://localhost:8080/health
```

### Create a snippet

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"content":"hello from Pastebin Lite"}' \
  http://localhost:8080/api/snippets
```

### Retrieve a snippet

```bash
curl http://localhost:8080/api/snippets/<slug>
```

### Delete a snippet

```bash
curl -X DELETE \
  http://localhost:8080/api/snippets/<slug>
```

## CI/CD Goals

The final delivery workflow is designed to demonstrate:

* Automated testing on pull requests
* Static analysis
* Reproducible application builds
* Container image creation
* Immutable image tagging
* Container vulnerability scanning
* Container registry usage
* Automated staging deployment
* Kubernetes-based deployment
* Deployment health verification
* Release promotion
* Rollback capability
* Infrastructure/application observability

The goal is to maintain the following invariant:

```text
Source Commit
     │
     ▼
Tested Artifact
     │
     ▼
Same Artifact
     │
     ▼
Staging
     │
     ▼
Production
```

Production should not rebuild a different artifact from the same source commit.

## Development Philosophy

This project intentionally keeps the application small.

The interesting engineering problem is the **software delivery lifecycle**:

```text
Code
 ↓
Test
 ↓
Build
 ↓
Package
 ↓
Publish
 ↓
Deploy
 ↓
Verify
 ↓
Observe
 ↓
Rollback
```

The application provides the workload; the CI/CD system provides the engineering challenge.

## Project Status

### Application

* [x] Go backend
* [x] Health endpoint
* [x] Domain model
* [x] Validation
* [x] Slug generation
* [x] PostgreSQL repository
* [x] Database migrations
* [x] Create/Get/Delete API
* [x] Backend tests
* [x] React frontend
* [x] Frontend tests

### Delivery

* [x] Docker
* [x] Docker Compose
* [x] CI pipeline
* [x] Container registry
* [x] Security scanning
* [x] Kubernetes deployment
* [x] CD pipeline
* [x] GitOps
* [x] Deployment verification
* [x] Rollback
* [x] Observability

---

## Why This Project?

Pastebin Lite is deliberately simple enough that the entire application can be understood quickly.

That leaves the complexity where it belongs for this project:

**the delivery system.**

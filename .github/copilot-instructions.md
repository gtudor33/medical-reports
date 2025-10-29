# Medical Reports Platform - AI Agent Instructions

## Architecture Overview

This is a **hexagonal architecture Go application** for Romanian hospital medical reports with strict business rules and immutable audit trails.

**Core Pattern**: Domain → Services → Repository → HTTP handlers
- `internal/domain/` - Pure business models with validation logic
- `internal/services/` - Business rule enforcement and orchestration  
- `internal/repository/postgres/` - Data persistence with PostgreSQL + JSONB
- `server/` - Gin HTTP API with DTOs and middleware

### Dependency Injection Pattern
Main entry point (`cmd/api/main.go`) follows clean dependency injection:
```go
// Initialize repositories (outer layer)
reportRepo := postgres.NewReportRepository(db)
referenceRepo := postgres.NewReferenceRepository(db)

// Initialize services (business layer) 
reportService := services.NewReportService(reportRepo)

// Initialize server (HTTP layer)
srv := server.NewServer(cfg, reportService)
```

## Critical Business Rules (Enforced in Services Layer)

```go
// Reports can ONLY be edited in draft status
if report.Status != domain.StatusDraft {
    return domain.ErrCannotEditNonDraft
}

// Status transitions follow strict workflow
StatusDraft → StatusInReview → StatusApproved → StatusSigned
// Use ValidTransitions map in domain/report.go

// Auto-save creates immutable versions
// Every content update saves a ReportVersion record
```

## Key Development Patterns

### 1. Domain-First Design
All business logic lives in `domain/` models. Services coordinate, never implement rules:
```go
// In domain/report.go
func (r *Report) CanTransitionTo(newStatus Status) bool {
    // Business rule implementation here
}

// In services/report_service.go  
if !report.CanTransitionTo(newStatus) {
    return domain.ErrInvalidStatusTransition
}
```

### 2. Repository Interface Pattern
Always program against interfaces in `repository/interfaces.go`:
```go
type ReportRepository interface {
    Create(ctx context.Context, report *domain.Report) error
    // Implementation in repository/postgres/
}
```

### 3. Content Storage Pattern
Report content uses JSONB with typed sections:
```go
type ReportContent struct {
    PatientData PatientDataSection `json:"patient_data"`
    Anamnesis   AnamnesisSection   `json:"anamnesis"`  
    Diagnosis   DiagnosisSection   `json:"diagnosis"`
    // Each section has its own Validate() method
}
```

### 4. Error Handling Pattern
Domain errors flow through layers without transformation:
```go
// In services/report_service.go
if report.Status != domain.StatusDraft {
    return domain.ErrCannotEditNonDraft  // Domain error
}

// In server/handlers.go - map to HTTP status
switch {
case errors.Is(err, domain.ErrCannotEditNonDraft):
    c.JSON(400, ErrorResponse{Error: "cannot_edit", Message: err.Error()})
}
```

### 5. Repository Query Pattern
Use LATERAL JOIN for efficient version queries:
```go
// Get report with latest version content
query := `
SELECT r.*, v.content
FROM reports r
LEFT JOIN LATERAL (
    SELECT content FROM report_versions 
    WHERE report_id = r.id 
    ORDER BY version_number DESC LIMIT 1
) v ON true`
```

## Development Workflow

### Essential Commands
```bash
make dev              # Start PostgreSQL + API (ports 5433, 8080)
./test-api.sh         # Complete workflow demo + business rule validation
make check            # Verify code compiles without building
docker-compose logs   # Debug database issues
```

### Adding New Features
1. **Domain first**: Add models/rules in `internal/domain/`
2. **Service layer**: Add business logic in `internal/services/`
3. **Repository**: Extend interfaces, implement in `postgres/`
4. **HTTP**: Add handlers in `server/handlers.go`, routes in `server/server.go`

### Database Patterns
- **JSONB for content**: All report sections stored as JSONB for flexibility
- **Immutable versions**: `report_versions` table for audit compliance
- **Full-text search**: GIN indexes on Romanian ICD-10/medication data
- **Migrations**: Use `migrations/*.sql` files, Docker auto-applies on startup

## Romanian Medical Context

### ICD-10 Integration
```go
// Romanian-specific medical codes
GET /api/v1/reference/icd10?q=pneumonie
// Returns: J18.1 - "Pneumonie lobară nespecificată"
```

### CNP Validation
```go
// Romanian Personal ID validation (13 digits)
if len(patientCNP) != 13 {
    return domain.ErrInvalidCNP
}
```

## Testing & Debugging

### API Testing Pattern
Use `test-api.sh` as your primary testing tool - it demonstrates:
- Complete report creation workflow
- ICD-10/medication search
- Status transition validation
- Business rule enforcement
- Version history management

### Testing Workflow Example
```bash
# 1. Start services
make dev

# 2. Run complete test suite
./test-api.sh

# 3. Test specific business rule
curl -X PUT localhost:8080/api/v1/reports/$ID/content \
  -d '{"content": {...}}' 
# Should fail if report status != draft
```

### Database Debugging
```bash
# Connect to PostgreSQL
docker exec -it medical-reports-db psql -U medreport -d medical_reports

# Check report status transitions
SELECT id, status, last_modified FROM reports ORDER BY last_modified DESC;

# Check version history
SELECT report_id, version_number, saved_at FROM report_versions 
WHERE report_id = 'your-uuid' ORDER BY version_number;
```

### Error Handling Convention
Domain errors map to HTTP status codes in `server/handlers.go`:
```go
domain.ErrReportNotFound → 404
domain.ErrCannotEditNonDraft → 400  
domain.ErrInvalidStatusTransition → 400
```

## Docker & Environment

### Local Development
```bash
# Database on port 5433 (not 5432 to avoid conflicts)
# API on port 8080
# Frontend on port 3000
docker-compose up -d
```

### Environment Variables
```bash
DATABASE_URL=postgres://medreport:medreport_dev@localhost:5433/medical_reports?sslmode=disable
JWT_SECRET=medical-reports-secret-key-change-in-production
```

## Frontend Integration Patterns

### React + Axios API Client
Frontend uses Axios with interceptors for consistent API communication:
```javascript
// In frontend/src/services/api.js
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Auto-logout on 401
api.interceptors.response.use(
    response => response,
    error => {
        if (error.response?.status === 401) {
            localStorage.removeItem('token');
            window.location.href = '/login';
        }
    }
);
```

### Authentication Context Pattern
```jsx
// In frontend/src/contexts/AuthContext.jsx
const { user, login, logout } = useAuth();

// Login flow stores JWT and user data
const login = async (email, password) => {
    const response = await authAPI.login(email, password);
    const { token, user: userData } = response.data;
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(userData));
    setUser(userData);
};
```

## Security Implementation

- **JWT tokens**: 7-day expiry, includes user role and hospital_id
- **Password hashing**: bcrypt with default cost
- **SQL injection prevention**: Parameterized queries throughout
- **CORS middleware**: Configured in `server/middleware.go`

### JWT Middleware Pattern
```go
// In server/middleware.go
func JWTMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // Validate and extract claims
        c.Set("user_id", userID)  // Add to context
        c.Next()
    }
}
```

## Common Gotchas

1. **Always use absolute paths** in tools - workspace root is `/Users/tmiron/Documents/vscode-projects/h-concept/`
2. **Status transitions are unidirectional** - check ValidTransitions map before allowing changes
3. **Content updates require draft status** - this is enforced at service layer
4. **Version numbers are sequential** - get current count before creating new version
5. **Database connection uses port 5433** - not the standard 5432 to avoid conflicts

## File Organization

- `cmd/api/main.go` - Dependency injection and server startup
- `internal/domain/section.go` - Romanian medical report structure
- `internal/services/report_service.go` - Core business operations  
- `server/handlers.go` - HTTP request/response handling
- `test-api.sh` - Living documentation of API usage

## Advanced Development Patterns

### DTO Transformation Pattern
Always convert between domain models and DTOs at HTTP boundary:
```go
// In server/dto.go
func ToReportResponse(report *domain.Report) ReportResponse {
    return ReportResponse{
        ID:        report.ID.String(),
        Status:    string(report.Status),
        Content:   report.Content,
        // Convert UUIDs to strings for JSON
    }
}

// In handlers.go
response := ToReportResponse(report)
c.JSON(200, response)
```

### Database Migration Pattern
PostgreSQL migrations auto-apply on Docker startup:
```sql
-- migrations/001_initial_schema.up.sql
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content JSONB NOT NULL,  -- Flexible schema
    search_vector tsvector GENERATED ALWAYS AS (
        to_tsvector('romanian', description_ro)
    ) STORED  -- Romanian full-text search
);
```

### Full-Text Search Implementation
Romanian-specific search with GIN indexes:
```go
// In repository/postgres/reference_repository.go
query := `
SELECT code, description_ro, category 
FROM icd10_codes 
WHERE search_vector @@ plainto_tsquery('romanian', $1)
ORDER BY ts_rank(search_vector, plainto_tsquery('romanian', $1)) DESC
LIMIT $2`
```

## Production Deployment Patterns

### Docker Multi-Stage Builds
```dockerfile
# Dockerfile uses multi-stage build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
CMD ["./api"]
```

### Environment Configuration
```bash
# docker-compose.yml sets environment
DATABASE_URL=postgres://medreport:medreport_dev@postgres:5432/medical_reports
SERVER_PORT=8080
JWT_SECRET=medical-reports-secret-key-change-in-production
```

## Monitoring & Observability

### Request Logging Middleware
```go
// In server/middleware.go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := time.Since(start)
        log.Printf("[%s] %s - Status: %d - Duration: %v",
            c.Request.Method, c.Request.URL.Path, 
            c.Writer.Status(), duration)
    }
}
```

### Health Check Endpoint
```go
// Always implement comprehensive health checks
func (h *Handlers) HealthCheck(c *gin.Context) {
    c.JSON(200, gin.H{
        "status":  "healthy",
        "service": "medical-reports-api",
        "version": version,
        "db":      h.checkDB(),  // DB connectivity check
    })
}
```
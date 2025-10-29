# Medical Reports Platform - Project Structure

## 📂 Complete Directory Tree

```
medical-reports/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
│
├── internal/
│   ├── domain/                        # Domain layer (models only)
│   │   ├── errors.go                  # Domain-specific errors
│   │   ├── report.go                  # Report domain model
│   │   └── section.go                 # Section domain models
│   │
│   ├── services/                      # Business logic layer
│   │   ├── report_service.go          # Report business logic
│   │   └── reference_service.go       # Reference data logic
│   │
│   ├── repository/                    # Data access layer
│   │   ├── interfaces.go              # Repository interfaces
│   │   └── postgres/
│   │       ├── report_repository.go   # Report persistence
│   │       └── reference_repository.go # Reference data persistence
│   │
│   ├── server/                        # HTTP API layer
│   │   ├── server.go                  # Server setup and routing
│   │   ├── handlers.go                # HTTP handlers
│   │   ├── middleware.go              # Middleware (CORS, logging)
│   │   └── dto.go                     # Request/Response DTOs
│   │
│   └── config/
│       └── config.go                  # Configuration management
│
├── migrations/
│   ├── 001_initial_schema.up.sql      # Database schema
│   └── 001_initial_schema.down.sql    # Rollback migration
│
├── docker-compose.yml                  # Local development setup
├── Dockerfile                          # Container image
├── Makefile                            # Development commands
├── go.mod                              # Go dependencies
├── go.sum                              # Dependency checksums
├── test-api.sh                         # API testing script
├── README.md                           # Documentation
└── .gitignore                          # Git ignore rules
```

## 🏗️ Architecture Layers

### 1. Domain Layer (`internal/domain/`)
**Pure business models - NO external dependencies**

- `report.go`: Core Report aggregate with business rules
  - Report entity with status workflow
  - ReportVersion for immutability
  - Value objects: Specialty, ReportType, Status
  - Business invariants (e.g., CNP validation)

- `section.go`: Section models for report content
  - PatientDataSection
  - AnamnesisSection
  - ExaminationSection
  - DiagnosisSection
  - TreatmentSection
  - RecommendationsSection
  - Validation logic for each section

- `errors.go`: Domain-specific errors
  - ErrReportNotFound
  - ErrCannotEditNonDraft
  - ErrInvalidStatusTransition
  - ErrIncompleteReport

### 2. Service Layer (`internal/services/`)
**Business logic orchestration**

- `report_service.go`: Report business operations
  - CreateReport: Creates new draft with validation
  - UpdateReportContent: Updates content + saves version
  - UpdateReportStatus: Enforces workflow rules
  - ListReports: Query with filtering
  - DeleteReport: Only allows deleting drafts
  - GetReportVersions: Version history
  - RestoreVersion: Restore to previous version

- `reference_service.go`: Reference data operations
  - SearchICD10: Full-text search with Romanian support
  - SearchMedications: Autocomplete for medications

**Business Rules Enforced:**
- ✅ Reports can only be edited in draft status
- ✅ Status transitions follow defined workflow
- ✅ Reports must be complete before finalizing
- ✅ CNP must be 13 digits
- ✅ Versions are immutable

### 3. Repository Layer (`internal/repository/`)
**Data persistence abstraction**

- `interfaces.go`: Defines contracts
  - ReportRepository interface
  - ReferenceRepository interface

- `postgres/report_repository.go`: PostgreSQL implementation
  - CRUD operations
  - Version management
  - JSONB content storage
  - Full-text search

- `postgres/reference_repository.go`: Reference data
  - ICD-10 code search (Romanian full-text)
  - Medication search

### 4. Server Layer (`internal/server/`)
**HTTP API using Gin framework**

- `server.go`: Server initialization and routing
  ```
  GET  /health
  POST /api/v1/reports
  GET  /api/v1/reports
  GET  /api/v1/reports/:id
  PUT  /api/v1/reports/:id/content
  PUT  /api/v1/reports/:id/status
  DELETE /api/v1/reports/:id
  GET  /api/v1/reports/:id/versions
  GET  /api/v1/reference/icd10
  GET  /api/v1/reference/medications
  ```

- `handlers.go`: HTTP request handlers
  - Input validation
  - Error handling with proper HTTP codes
  - Domain error to HTTP status mapping

- `middleware.go`: Cross-cutting concerns
  - CORS handling
  - Request logging
  - Panic recovery
  - Error handling

- `dto.go`: Data Transfer Objects
  - Request/Response transformations
  - Domain model to JSON mapping

### 5. Configuration (`internal/config/`)
- Environment variable loading
- Default values
- Database URL
- Server host/port

## 🗄️ Database Schema

### Tables

**reports**
- Main report metadata
- Links to patient and doctor
- Workflow status
- Audit timestamps

**report_versions**
- Immutable content snapshots
- JSONB storage for flexibility
- Version numbering
- Change comments

**icd10_codes**
- ICD-10 reference data
- Romanian descriptions
- Full-text search support
- Pre-seeded with common codes

**medications**
- Romanian pharmaceutical database
- Full-text search
- Pre-seeded with common medications

**audit_log**
- Immutable event log
- User actions
- IP addresses
- GDPR compliance

## 🚀 Key Features Implemented

### ✅ Report Management
- Create draft reports
- Update content with auto-save
- Version history (immutable snapshots)
- Restore previous versions
- Status workflow (draft → in_review → approved → signed)
- Business rule enforcement

### ✅ Reference Data
- ICD-10 code autocomplete (Romanian)
- Medication autocomplete (Romanian pharmaceutical database)
- Full-text search with PostgreSQL

### ✅ Data Integrity
- Immutable versions (audit trail)
- Status transition validation
- CNP validation (Romanian personal ID)
- Date logic validation

### ✅ API Design
- RESTful endpoints
- JSON request/response
- Proper HTTP status codes
- Error messages with context
- CORS support for frontend integration

## 🔒 Security & Compliance Features

### Current Implementation
- Input validation
- SQL injection prevention (parameterized queries)
- CORS configuration
- Request logging
- Error handling without information leakage

### Ready for Phase B/C
- Authentication (JWT/OAuth) - interfaces ready
- Authorization (RBAC) - can be added in middleware
- Encryption at rest - PostgreSQL supports it
- Audit logging - already implemented
- Digital signatures - schema ready

## 🧪 Testing

### Test Script (`test-api.sh`)
Complete workflow demonstration:
1. Health check
2. Search ICD-10 codes
3. Search medications
4. Create report
5. Update content
6. Get versions
7. Change status
8. List reports
9. Test business rules (editing non-draft)

### Usage
```bash
# Start services
docker compose up -d

# Run tests
./test-api.sh
```

## 📊 Performance Considerations

### Indexing Strategy
- `idx_reports_hospital_patient`: Patient lookups
- `idx_reports_doctor_status`: Doctor's report list
- `idx_reports_created_at`: Temporal queries
- `idx_versions_report`: Version history
- `idx_icd10_search`: Full-text search (GIN)
- `idx_medications_search`: Full-text search (GIN)

### Query Optimization
- LATERAL JOIN for latest version
- Parameterized queries
- Limit/offset pagination
- Full-text search caching

## 🔄 Development Workflow

```bash
# Start development environment
make dev

# Run locally (without Docker)
make run

# Build binary
make build

# Run tests
make test

# Format and lint
make lint

# Check compilation
make check
```

## 📈 Scalability Path

### Horizontal Scaling
- Stateless API (can run multiple instances)
- PostgreSQL read replicas
- Connection pooling
- Caching layer (Redis) for reference data

### Vertical Improvements
- Database query optimization
- JSONB indexing for content
- Materialized views for analytics
- Partitioning for large datasets

## 🎯 Next Integration Steps

### Phase B (MVP Enhancement)
1. Add authentication service
2. Add PDF generation service
3. Add AI rephrasing service (Python microservice)
4. Add speech-to-text (Whisper)

### Phase C (Production)
1. HL7/FHIR integration
2. Digital signature service
3. Hospital system connectors
4. Keycloak for SSO

## 📝 Code Quality

### Standards Followed
- Clean architecture principles
- Hexagonal architecture (ports & adapters)
- Domain-driven design patterns
- SOLID principles
- Clear separation of concerns
- Interface-based design
- Error handling best practices

### Go Best Practices
- No global state
- Context propagation
- Error wrapping
- Dependency injection
- Interface segregation
- Minimal package coupling

## 🎓 Learning Resources

For team members new to the codebase:

1. Start with `cmd/api/main.go` - entry point
2. Read `internal/domain/` - understand the business
3. Review `internal/services/` - see business logic
4. Check `internal/server/server.go` - understand API
5. Run `test-api.sh` - see it in action

---

**Status**: ✅ Ready for demo at Spitalul Sf. Spiridon
**Next**: Show to doctors, gather feedback, iterate on UX

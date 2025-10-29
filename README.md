# Medical Reports System

A GDPR-compliant medical report automation platform for Romanian hospitals, with a modern web interface.

**Tech Stack:**

- Backend: Go + PostgreSQL
- Frontend: React + Tailwind CSS
- Deployment: Docker + Docker Compose

## Features

- ✅ Create and manage medical discharge reports
- ✅ Auto-save with version history
- ✅ ICD-10 code autocomplete (Romanian)
- ✅ Medication database autocomplete
- ✅ Workflow status management (Draft → In Review → Approved → Signed)
- ✅ RESTful API with JSON responses
- ✅ PostgreSQL with full-text search

## Architecture

```
medical-reports/
├── cmd/api/              # Application entry point
├── internal/
│   ├── domain/           # Domain models (Report, Section, etc.)
│   ├── services/         # Business logic layer
│   ├── repository/       # Data persistence layer
│   ├── server/           # HTTP server and handlers
│   └── config/           # Configuration management
├── migrations/           # Database migrations
└── docker-compose.yml    # Local development setup
```

## Quick Start

### Easiest Way - One Command! 🚀

```bash
./start-dev.sh
```

This script will:

- Build and start all Docker containers
- Set up the database with migrations
- Start the backend API on port 8080
- Start the frontend on port 3000

**Access the application:**

- 🌐 **Frontend UI**: [http://localhost:3000](http://localhost:3000)
- 🔌 **Backend API**: [http://localhost:8080](http://localhost:8080)
- 🗄️ **Database**: localhost:5433

### Manual Setup

**Prerequisites:**

- Docker & Docker Compose
- Go 1.21+ (for local development without Docker)
- Node.js 20+ (for frontend development)

**Start with Docker Compose:**

```bash
docker-compose up --build
```

This will start:

- PostgreSQL on port 5433
- Backend API on port 8080
- Frontend on port 3000

**Verify it's running:**
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "medical-reports-api"
}
```

### Option 2: Run locally (without Docker)

1. **Start PostgreSQL**
```bash
docker-compose up -d postgres
```

2. **Install dependencies**
```bash
go mod download
```

3. **Run the application**
```bash
go run cmd/api/main.go
```

## API Endpoints

### Health Check
```bash
GET /health
```

### Reports

#### Create a new report
```bash
POST /api/v1/reports
Content-Type: application/json

{
  "hospital_id": "550e8400-e29b-41d4-a716-446655440000",
  "patient_cnp": "1850312400123",
  "patient_first_name": "Ion",
  "patient_last_name": "Popescu",
  "specialty": "internal_medicine",
  "report_type": "discharge_summary",
  "doctor_id": "660e8400-e29b-41d4-a716-446655440001"
}
```

#### Get a report by ID
```bash
GET /api/v1/reports/{report_id}
```

#### List reports
```bash
GET /api/v1/reports?doctor_id={doctor_id}&status=draft&limit=20&offset=0
```

#### Update report content
```bash
PUT /api/v1/reports/{report_id}/content
Content-Type: application/json

{
  "user_id": "660e8400-e29b-41d4-a716-446655440001",
  "content": {
    "patient_data": {
      "first_name": "Ion",
      "last_name": "Popescu",
      "cnp": "1850312400123",
      "birth_date": "1985-03-12T00:00:00Z",
      "department": "Medicină Internă",
      "ward": "12",
      "bed": "3",
      "admission_date": "2025-10-20T08:00:00Z",
      "discharge_date": "2025-10-29T14:00:00Z"
    },
    "anamnesis": {
      "chief_complaint": "Dispnee și tuse productivă",
      "history_of_present_illness": "Pacientul se prezintă cu dispnee...",
      "past_medical_history": "HTA diagnosticată în 2020",
      "allergies": "Negat",
      "social_history": "Fumător 20 țigări/zi"
    },
    "examination": {
      "general_condition": "Stare generală bună",
      "consciousness": "Lucidă",
      "vital_signs": {
        "blood_pressure": "130/80",
        "heart_rate": 82,
        "temperature": 36.8,
        "respiratory_rate": 18,
        "oxygen_saturation": 96
      },
      "systems_review": "Aparat respirator: MV prezente bilateral..."
    },
    "diagnosis": {
      "primary_diagnosis": {
        "code": "J18.1",
        "description": "Pneumonie lobară nespecificată"
      },
      "secondary_diagnoses": [
        {
          "code": "I10",
          "description": "Hipertensiune arterială esențială"
        }
      ],
      "clinical_observations": "Evoluție favorabilă sub tratament"
    },
    "treatment": {
      "medications": [
        {
          "name": "Amoxicilină + Acid Clavulanic",
          "dosage": "1g/200mg",
          "frequency": "3x/zi",
          "route": "iv",
          "start_date": "2025-10-20T08:00:00Z",
          "end_date": "2025-10-27T08:00:00Z"
        }
      ],
      "procedures": []
    },
    "recommendations": {
      "discharge_plan": "Pacient externat cu stare generală bună",
      "medications": "Continuare Enalapril 10mg 1cp/zi",
      "follow_up": "Control cardiologie peste 1 lună",
      "diet_restrictions": "Regim hiposodat",
      "activity_restrictions": "Repaus relativ 2 săptămâni"
    }
  }
}
```

#### Update report status
```bash
PUT /api/v1/reports/{report_id}/status
Content-Type: application/json

{
  "status": "in_review"
}
```

Status values: `draft`, `in_review`, `approved`, `signed`, `cancelled`

#### Delete a report (only drafts)
```bash
DELETE /api/v1/reports/{report_id}
```

#### Get report versions
```bash
GET /api/v1/reports/{report_id}/versions
```

### Reference Data

#### Search ICD-10 codes
```bash
GET /api/v1/reference/icd10?q=pneumonie
```

Response:
```json
[
  {
    "code": "J18.1",
    "description": "Pneumonie lobară nespecificată",
    "category": "Respiratory"
  },
  {
    "code": "J18.0",
    "description": "Bronhopneumonie nespecificată",
    "category": "Respiratory"
  }
]
```

#### Search medications
```bash
GET /api/v1/reference/medications?q=amox
```

Response:
```json
[
  {
    "id": "...",
    "name": "Augmentin",
    "active_substance": "Amoxicilină + Acid Clavulanic",
    "form": "tablet",
    "dosage": "1g/200mg",
    "manufacturer": "GSK"
  }
]
```

## Testing with curl

### Complete workflow example

```bash
# 1. Create a report
REPORT_ID=$(curl -s -X POST http://localhost:8080/api/v1/reports \
  -H "Content-Type: application/json" \
  -d '{
    "hospital_id": "550e8400-e29b-41d4-a716-446655440000",
    "patient_cnp": "1850312400123",
    "patient_first_name": "Ion",
    "patient_last_name": "Popescu",
    "specialty": "internal_medicine",
    "report_type": "discharge_summary",
    "doctor_id": "660e8400-e29b-41d4-a716-446655440001"
  }' | jq -r '.id')

echo "Created report: $REPORT_ID"

# 2. Search for ICD-10 code
curl http://localhost:8080/api/v1/reference/icd10?q=pneumonie | jq

# 3. Update report content
curl -X PUT http://localhost:8080/api/v1/reports/$REPORT_ID/content \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "660e8400-e29b-41d4-a716-446655440001",
    "content": {
      "patient_data": {
        "first_name": "Ion",
        "last_name": "Popescu",
        "cnp": "1850312400123",
        "birth_date": "1985-03-12T00:00:00Z",
        "department": "Medicină Internă",
        "ward": "12",
        "bed": "3",
        "admission_date": "2025-10-20T08:00:00Z",
        "discharge_date": "2025-10-29T14:00:00Z"
      },
      "anamnesis": {
        "chief_complaint": "Dispnee și tuse",
        "history_of_present_illness": "Test",
        "past_medical_history": "Test",
        "allergies": "None",
        "social_history": "Test"
      },
      "examination": {
        "general_condition": "Bună",
        "consciousness": "Lucidă",
        "vital_signs": {
          "blood_pressure": "130/80",
          "heart_rate": 82,
          "temperature": 36.8,
          "respiratory_rate": 18,
          "oxygen_saturation": 96
        },
        "systems_review": "Normal"
      },
      "diagnosis": {
        "primary_diagnosis": {
          "code": "J18.1",
          "description": "Pneumonie lobară"
        },
        "secondary_diagnoses": [],
        "clinical_observations": "Good"
      },
      "treatment": {
        "medications": [],
        "procedures": []
      },
      "recommendations": {
        "discharge_plan": "Rest",
        "medications": "None",
        "follow_up": "1 month",
        "diet_restrictions": "None",
        "activity_restrictions": "None"
      }
    }
  }' | jq

# 4. Get the updated report
curl http://localhost:8080/api/v1/reports/$REPORT_ID | jq

# 5. Change status to in_review
curl -X PUT http://localhost:8080/api/v1/reports/$REPORT_ID/status \
  -H "Content-Type: application/json" \
  -d '{"status": "in_review"}' | jq

# 6. List all reports
curl "http://localhost:8080/api/v1/reports?doctor_id=660e8400-e29b-41d4-a716-446655440001&limit=10" | jq
```

## Database Schema

The application uses PostgreSQL with the following main tables:
- `reports` - Main report data
- `report_versions` - Immutable version history
- `icd10_codes` - ICD-10 code reference (seeded with common codes)
- `medications` - Medication reference (seeded with common medications)
- `audit_log` - Audit trail

## Development

### Project Structure

- **Domain Layer** (`internal/domain/`) - Pure business models with no external dependencies
- **Service Layer** (`internal/services/`) - Business logic and orchestration
- **Repository Layer** (`internal/repository/`) - Data access with PostgreSQL implementation
- **Server Layer** (`internal/server/`) - HTTP API with Gin framework

### Adding New Features

1. **Add domain models** in `internal/domain/`
2. **Add business logic** in `internal/services/`
3. **Add repository methods** in `internal/repository/postgres/`
4. **Add HTTP handlers** in `internal/server/handlers.go`
5. **Add routes** in `internal/server/server.go`

## Environment Variables

```bash
DATABASE_URL=postgres://medreport:medreport_dev@localhost:5432/medical_reports?sslmode=disable
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

## Stopping the Services

```bash
docker-compose down
```

To also remove volumes (database data):
```bash
docker-compose down -v
```

## Next Steps

- [ ] Add authentication (JWT/OAuth)
- [ ] Add PDF generation
- [ ] Add AI text rephrasing (Phase B)
- [ ] Add speech-to-text dictation (Phase B)
- [ ] Add hospital system integration (HL7/FHIR) (Phase C)
- [ ] Add digital signatures (Phase C)

## License

Proprietary - Tudor Miron

## Contact

Tudor Miron - Iași, Romania

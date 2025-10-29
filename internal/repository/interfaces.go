package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tudormiron/medical-reports/internal/domain"
)

// ReportRepository defines persistence interface for reports
type ReportRepository interface {
	Create(ctx context.Context, report *domain.Report) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Report, error)
	Update(ctx context.Context, report *domain.Report) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, doctorID uuid.UUID, status domain.Status, limit, offset int) ([]*domain.Report, error)
	
	// Version management
	SaveVersion(ctx context.Context, version *domain.ReportVersion) error
	GetVersions(ctx context.Context, reportID uuid.UUID) ([]*domain.ReportVersion, error)
	GetVersion(ctx context.Context, reportID uuid.UUID, versionNumber int) (*domain.ReportVersion, error)
}

// ReferenceRepository defines interface for reference data (ICD-10, medications)
type ReferenceRepository interface {
	SearchICD10(ctx context.Context, query string, limit int) ([]domain.ICD10Reference, error)
	GetICD10ByCode(ctx context.Context, code string) (*domain.ICD10Reference, error)
	
	SearchMedications(ctx context.Context, query string, limit int) ([]domain.MedicationReference, error)
	GetMedicationByID(ctx context.Context, id uuid.UUID) (*domain.MedicationReference, error)
}

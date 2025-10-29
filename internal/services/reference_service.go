package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/tudormiron/medical-reports/internal/domain"
	"github.com/tudormiron/medical-reports/internal/repository"
)

type ReferenceService struct {
	refRepo repository.ReferenceRepository
}

func NewReferenceService(refRepo repository.ReferenceRepository) *ReferenceService {
	return &ReferenceService{
		refRepo: refRepo,
	}
}

// SearchICD10 searches for ICD-10 codes
func (s *ReferenceService) SearchICD10(ctx context.Context, query string) ([]domain.ICD10Reference, error) {
	if query == "" {
		return []domain.ICD10Reference{}, nil
	}
	
	return s.refRepo.SearchICD10(ctx, query, 10)
}

// GetICD10ByCode retrieves a specific ICD-10 code
func (s *ReferenceService) GetICD10ByCode(ctx context.Context, code string) (*domain.ICD10Reference, error) {
	return s.refRepo.GetICD10ByCode(ctx, code)
}

// SearchMedications searches for medications
func (s *ReferenceService) SearchMedications(ctx context.Context, query string) ([]domain.MedicationReference, error) {
	if query == "" {
		return []domain.MedicationReference{}, nil
	}
	
	return s.refRepo.SearchMedications(ctx, query, 10)
}

// GetMedicationByID retrieves a specific medication
func (s *ReferenceService) GetMedicationByID(ctx context.Context, id uuid.UUID) (*domain.MedicationReference, error) {
	return s.refRepo.GetMedicationByID(ctx, id)
}

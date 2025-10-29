package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tudormiron/medical-reports/internal/domain"
)

type ReferenceRepository struct {
	db *sql.DB
}

func NewReferenceRepository(db *sql.DB) *ReferenceRepository {
	return &ReferenceRepository{db: db}
}

func (r *ReferenceRepository) SearchICD10(ctx context.Context, query string, limit int) ([]domain.ICD10Reference, error) {
	sqlQuery := `
		SELECT code, description_ro, category
		FROM icd10_codes
		WHERE search_vector @@ plainto_tsquery('romanian', $1)
		   OR code ILIKE $2
		ORDER BY 
			CASE 
				WHEN code ILIKE $2 THEN 0
				ELSE 1
			END,
			code
		LIMIT $3
	`
	
	searchPattern := query + "%"
	rows, err := r.db.QueryContext(ctx, sqlQuery, query, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var results []domain.ICD10Reference
	for rows.Next() {
		var ref domain.ICD10Reference
		if err := rows.Scan(&ref.Code, &ref.DescriptionRO, &ref.Category); err != nil {
			return nil, err
		}
		results = append(results, ref)
	}
	
	return results, nil
}

func (r *ReferenceRepository) GetICD10ByCode(ctx context.Context, code string) (*domain.ICD10Reference, error) {
	query := `
		SELECT code, description_ro, category
		FROM icd10_codes
		WHERE code = $1
	`
	
	var ref domain.ICD10Reference
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&ref.Code,
		&ref.DescriptionRO,
		&ref.Category,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrReportNotFound
		}
		return nil, err
	}
	
	return &ref, nil
}

func (r *ReferenceRepository) SearchMedications(ctx context.Context, query string, limit int) ([]domain.MedicationReference, error) {
	sqlQuery := `
		SELECT id, name, active_substance, form, dosage, manufacturer
		FROM medications
		WHERE search_vector @@ plainto_tsquery('romanian', $1)
		   OR name ILIKE $2
		ORDER BY 
			CASE 
				WHEN name ILIKE $2 THEN 0
				ELSE 1
			END,
			name
		LIMIT $3
	`
	
	searchPattern := query + "%"
	rows, err := r.db.QueryContext(ctx, sqlQuery, query, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var results []domain.MedicationReference
	for rows.Next() {
		var ref domain.MedicationReference
		var form, dosage, manufacturer sql.NullString
		
		if err := rows.Scan(
			&ref.ID,
			&ref.Name,
			&ref.ActiveSubstance,
			&form,
			&dosage,
			&manufacturer,
		); err != nil {
			return nil, err
		}
		
		if form.Valid {
			ref.Form = form.String
		}
		if dosage.Valid {
			ref.Dosage = dosage.String
		}
		if manufacturer.Valid {
			ref.Manufacturer = manufacturer.String
		}
		
		results = append(results, ref)
	}
	
	return results, nil
}

func (r *ReferenceRepository) GetMedicationByID(ctx context.Context, id uuid.UUID) (*domain.MedicationReference, error) {
	query := `
		SELECT id, name, active_substance, form, dosage, manufacturer
		FROM medications
		WHERE id = $1
	`
	
	var ref domain.MedicationReference
	var form, dosage, manufacturer sql.NullString
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ref.ID,
		&ref.Name,
		&ref.ActiveSubstance,
		&form,
		&dosage,
		&manufacturer,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrReportNotFound
		}
		return nil, err
	}
	
	if form.Valid {
		ref.Form = form.String
	}
	if dosage.Valid {
		ref.Dosage = dosage.String
	}
	if manufacturer.Valid {
		ref.Manufacturer = manufacturer.String
	}
	
	return &ref, nil
}

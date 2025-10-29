package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/tudormiron/medical-reports/internal/domain"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(ctx context.Context, report *domain.Report) error {
	query := `
		INSERT INTO reports (
			id, hospital_id, patient_cnp, patient_first_name, patient_last_name,
			specialty, report_type, status, created_by, created_at, last_modified
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		report.ID,
		report.HospitalID,
		report.PatientCNP,
		report.PatientFirstName,
		report.PatientLastName,
		report.Specialty,
		report.ReportType,
		report.Status,
		report.CreatedBy,
		report.CreatedAt,
		report.LastModified,
	)
	
	if err != nil {
		return domain.ErrDatabaseQuery
	}
	
	// Save initial version
	version := domain.NewReportVersion(report.ID, 1, report.Content, report.CreatedBy, "Initial version")
	return r.SaveVersion(ctx, version)
}

func (r *ReportRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Report, error) {
	query := `
		SELECT 
			r.id, r.hospital_id, r.patient_cnp, r.patient_first_name, r.patient_last_name,
			r.specialty, r.report_type, r.status, r.created_by, r.created_at, 
			r.last_modified, r.finalized_at, v.content
		FROM reports r
		LEFT JOIN LATERAL (
			SELECT content 
			FROM report_versions 
			WHERE report_id = r.id 
			ORDER BY version_number DESC 
			LIMIT 1
		) v ON true
		WHERE r.id = $1
	`
	
	var report domain.Report
	var contentJSON []byte
	var finalizedAt sql.NullTime
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&report.ID,
		&report.HospitalID,
		&report.PatientCNP,
		&report.PatientFirstName,
		&report.PatientLastName,
		&report.Specialty,
		&report.ReportType,
		&report.Status,
		&report.CreatedBy,
		&report.CreatedAt,
		&report.LastModified,
		&finalizedAt,
		&contentJSON,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrReportNotFound
		}
		return nil, domain.ErrDatabaseQuery
	}
	
	if finalizedAt.Valid {
		report.FinalizedAt = &finalizedAt.Time
	}
	
	if contentJSON != nil {
		if err := json.Unmarshal(contentJSON, &report.Content); err != nil {
			return nil, err
		}
	}
	
	return &report, nil
}

func (r *ReportRepository) Update(ctx context.Context, report *domain.Report) error {
	query := `
		UPDATE reports 
		SET status = $1, last_modified = $2, finalized_at = $3
		WHERE id = $4
	`
	
	_, err := r.db.ExecContext(ctx, query,
		report.Status,
		report.LastModified,
		report.FinalizedAt,
		report.ID,
	)
	
	if err != nil {
		return domain.ErrDatabaseQuery
	}
	
	return nil
}

func (r *ReportRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reports WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ReportRepository) List(ctx context.Context, doctorID uuid.UUID, status domain.Status, limit, offset int) ([]*domain.Report, error) {
	query := `
		SELECT 
			r.id, r.hospital_id, r.patient_cnp, r.patient_first_name, r.patient_last_name,
			r.specialty, r.report_type, r.status, r.created_by, r.created_at, 
			r.last_modified, r.finalized_at, v.content
		FROM reports r
		LEFT JOIN LATERAL (
			SELECT content 
			FROM report_versions 
			WHERE report_id = r.id 
			ORDER BY version_number DESC 
			LIMIT 1
		) v ON true
		WHERE r.created_by = $1
	`
	
	args := []interface{}{doctorID}
	argCount := 1
	
	if status != "" {
		argCount++
		query += ` AND r.status = $` + strconv.Itoa(argCount)
		args = append(args, status)
	}

	query += ` ORDER BY r.last_modified DESC LIMIT $` + strconv.Itoa(argCount+1) + ` OFFSET $` + strconv.Itoa(argCount+2)
	args = append(args, limit, offset)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}
	defer rows.Close()
	
	var reports []*domain.Report
	for rows.Next() {
		var report domain.Report
		var contentJSON []byte
		var finalizedAt sql.NullTime
		
		err := rows.Scan(
			&report.ID,
			&report.HospitalID,
			&report.PatientCNP,
			&report.PatientFirstName,
			&report.PatientLastName,
			&report.Specialty,
			&report.ReportType,
			&report.Status,
			&report.CreatedBy,
			&report.CreatedAt,
			&report.LastModified,
			&finalizedAt,
			&contentJSON,
		)
		
		if err != nil {
			return nil, err
		}
		
		if finalizedAt.Valid {
			report.FinalizedAt = &finalizedAt.Time
		}
		
		if contentJSON != nil {
			if err := json.Unmarshal(contentJSON, &report.Content); err != nil {
				return nil, err
			}
		}
		
		reports = append(reports, &report)
	}
	
	return reports, nil
}

func (r *ReportRepository) SaveVersion(ctx context.Context, version *domain.ReportVersion) error {
	contentJSON, err := json.Marshal(version.Content)
	if err != nil {
		return err
	}
	
	query := `
		INSERT INTO report_versions (
			id, report_id, version_number, content, saved_at, saved_by, comment
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err = r.db.ExecContext(ctx, query,
		version.ID,
		version.ReportID,
		version.VersionNumber,
		contentJSON,
		version.SavedAt,
		version.SavedBy,
		version.Comment,
	)
	
	return err
}

func (r *ReportRepository) GetVersions(ctx context.Context, reportID uuid.UUID) ([]*domain.ReportVersion, error) {
	query := `
		SELECT id, report_id, version_number, content, saved_at, saved_by, comment
		FROM report_versions
		WHERE report_id = $1
		ORDER BY version_number DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var versions []*domain.ReportVersion
	for rows.Next() {
		var version domain.ReportVersion
		var contentJSON []byte
		
		err := rows.Scan(
			&version.ID,
			&version.ReportID,
			&version.VersionNumber,
			&contentJSON,
			&version.SavedAt,
			&version.SavedBy,
			&version.Comment,
		)
		
		if err != nil {
			return nil, err
		}
		
		if err := json.Unmarshal(contentJSON, &version.Content); err != nil {
			return nil, err
		}
		
		versions = append(versions, &version)
	}
	
	return versions, nil
}

func (r *ReportRepository) GetVersion(ctx context.Context, reportID uuid.UUID, versionNumber int) (*domain.ReportVersion, error) {
	query := `
		SELECT id, report_id, version_number, content, saved_at, saved_by, comment
		FROM report_versions
		WHERE report_id = $1 AND version_number = $2
	`
	
	var version domain.ReportVersion
	var contentJSON []byte
	
	err := r.db.QueryRowContext(ctx, query, reportID, versionNumber).Scan(
		&version.ID,
		&version.ReportID,
		&version.VersionNumber,
		&contentJSON,
		&version.SavedAt,
		&version.SavedBy,
		&version.Comment,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrReportNotFound
		}
		return nil, err
	}
	
	if err := json.Unmarshal(contentJSON, &version.Content); err != nil {
		return nil, err
	}
	
	return &version, nil
}

package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tudormiron/medical-reports/internal/domain"
	"github.com/tudormiron/medical-reports/internal/repository"
)

type ReportService struct {
	reportRepo repository.ReportRepository
}

func NewReportService(reportRepo repository.ReportRepository) *ReportService {
	return &ReportService{
		reportRepo: reportRepo,
	}
}

// CreateReport creates a new report
func (s *ReportService) CreateReport(ctx context.Context, hospitalID uuid.UUID, patientCNP, firstName, lastName string, specialty domain.Specialty, reportType domain.ReportType, doctorID uuid.UUID) (*domain.Report, error) {
	// Validate CNP
	if len(patientCNP) != 13 {
		return nil, domain.ErrInvalidCNP
	}
	
	report := domain.NewReport(hospitalID, patientCNP, firstName, lastName, specialty, reportType, doctorID)
	
	if err := s.reportRepo.Create(ctx, report); err != nil {
		return nil, err
	}
	
	return report, nil
}

// GetReport retrieves a report by ID
func (s *ReportService) GetReport(ctx context.Context, id uuid.UUID) (*domain.Report, error) {
	return s.reportRepo.GetByID(ctx, id)
}

// UpdateReportContent updates the content of a report
func (s *ReportService) UpdateReportContent(ctx context.Context, reportID uuid.UUID, content domain.ReportContent, userID uuid.UUID) error {
	report, err := s.reportRepo.GetByID(ctx, reportID)
	if err != nil {
		return err
	}
	
	// Business rule: Can only edit draft reports
	if report.Status != domain.StatusDraft {
		return domain.ErrCannotEditNonDraft
	}
	
	report.Content = content
	report.LastModified = time.Now()
	
	if err := s.reportRepo.Update(ctx, report); err != nil {
		return err
	}
	
	// Get current version count
	versions, err := s.reportRepo.GetVersions(ctx, reportID)
	if err != nil {
		return err
	}
	
	// Save new version
	versionNumber := len(versions) + 1
	version := domain.NewReportVersion(reportID, versionNumber, content, userID, "Auto-save")
	
	return s.reportRepo.SaveVersion(ctx, version)
}

// UpdateReportStatus changes the status of a report
func (s *ReportService) UpdateReportStatus(ctx context.Context, reportID uuid.UUID, newStatus domain.Status) error {
	report, err := s.reportRepo.GetByID(ctx, reportID)
	if err != nil {
		return err
	}
	
	// Business rule: Validate status transition
	if !report.CanTransitionTo(newStatus) {
		return domain.ErrInvalidStatusTransition
	}
	
	// Business rule: Report must be complete before finalizing
	if newStatus == domain.StatusInReview && !report.Content.IsComplete() {
		return domain.ErrIncompleteReport
	}
	
	report.Status = newStatus
	report.LastModified = time.Now()
	
	if newStatus == domain.StatusApproved || newStatus == domain.StatusSigned {
		now := time.Now()
		report.FinalizedAt = &now
	}
	
	return s.reportRepo.Update(ctx, report)
}

// ListReports lists reports for a doctor with filtering
func (s *ReportService) ListReports(ctx context.Context, doctorID uuid.UUID, status domain.Status, limit, offset int) ([]*domain.Report, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	
	return s.reportRepo.List(ctx, doctorID, status, limit, offset)
}

// DeleteReport deletes a report (only drafts)
func (s *ReportService) DeleteReport(ctx context.Context, reportID uuid.UUID) error {
	report, err := s.reportRepo.GetByID(ctx, reportID)
	if err != nil {
		return err
	}
	
	// Business rule: Can only delete drafts
	if report.Status != domain.StatusDraft {
		return domain.ErrCannotEditNonDraft
	}
	
	return s.reportRepo.Delete(ctx, reportID)
}

// GetReportVersions retrieves all versions of a report
func (s *ReportService) GetReportVersions(ctx context.Context, reportID uuid.UUID) ([]*domain.ReportVersion, error) {
	// Verify report exists
	if _, err := s.reportRepo.GetByID(ctx, reportID); err != nil {
		return nil, err
	}
	
	return s.reportRepo.GetVersions(ctx, reportID)
}

// RestoreVersion restores a report to a previous version
func (s *ReportService) RestoreVersion(ctx context.Context, reportID uuid.UUID, versionNumber int, userID uuid.UUID) error {
	report, err := s.reportRepo.GetByID(ctx, reportID)
	if err != nil {
		return err
	}
	
	// Business rule: Can only restore draft reports
	if report.Status != domain.StatusDraft {
		return domain.ErrCannotEditNonDraft
	}
	
	version, err := s.reportRepo.GetVersion(ctx, reportID, versionNumber)
	if err != nil {
		return err
	}
	
	report.Content = version.Content
	report.LastModified = time.Now()
	
	if err := s.reportRepo.Update(ctx, report); err != nil {
		return err
	}
	
	// Save as new version
	versions, err := s.reportRepo.GetVersions(ctx, reportID)
	if err != nil {
		return err
	}
	
	newVersionNumber := len(versions) + 1
	newVersion := domain.NewReportVersion(
		reportID,
		newVersionNumber,
		version.Content,
		userID,
		"Restored from version "+string(rune(versionNumber)),
	)
	
	return s.reportRepo.SaveVersion(ctx, newVersion)
}

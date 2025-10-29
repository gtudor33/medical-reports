package server

import (
	"time"

	"github.com/google/uuid"
	"github.com/tudormiron/medical-reports/internal/domain"
)

// Request DTOs
type CreateReportRequest struct {
	HospitalID        string `json:"hospital_id" binding:"required"`
	PatientCNP        string `json:"patient_cnp" binding:"required,len=13"`
	PatientFirstName  string `json:"patient_first_name" binding:"required"`
	PatientLastName   string `json:"patient_last_name" binding:"required"`
	Specialty         string `json:"specialty" binding:"required"`
	ReportType        string `json:"report_type" binding:"required"`
	DoctorID          string `json:"doctor_id" binding:"required"`
}

type UpdateReportContentRequest struct {
	Content domain.ReportContent `json:"content" binding:"required"`
	UserID  string               `json:"user_id" binding:"required"`
}

type UpdateReportStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// Response DTOs
type ReportResponse struct {
	ID               string                `json:"id"`
	HospitalID       string                `json:"hospital_id"`
	PatientCNP       string                `json:"patient_cnp"`
	PatientFirstName string                `json:"patient_first_name"`
	PatientLastName  string                `json:"patient_last_name"`
	Specialty        string                `json:"specialty"`
	ReportType       string                `json:"report_type"`
	Status           string                `json:"status"`
	Content          domain.ReportContent  `json:"content"`
	CreatedBy        string                `json:"created_by"`
	CreatedAt        time.Time             `json:"created_at"`
	LastModified     time.Time             `json:"last_modified"`
	FinalizedAt      *time.Time            `json:"finalized_at,omitempty"`
}

func ToReportResponse(report *domain.Report) ReportResponse {
	return ReportResponse{
		ID:               report.ID.String(),
		HospitalID:       report.HospitalID.String(),
		PatientCNP:       report.PatientCNP,
		PatientFirstName: report.PatientFirstName,
		PatientLastName:  report.PatientLastName,
		Specialty:        string(report.Specialty),
		ReportType:       string(report.ReportType),
		Status:           string(report.Status),
		Content:          report.Content,
		CreatedBy:        report.CreatedBy.String(),
		CreatedAt:        report.CreatedAt,
		LastModified:     report.LastModified,
		FinalizedAt:      report.FinalizedAt,
	}
}

type ReportListResponse struct {
	Reports []ReportResponse `json:"reports"`
	Total   int              `json:"total"`
	Limit   int              `json:"limit"`
	Offset  int              `json:"offset"`
}

type VersionResponse struct {
	ID            string               `json:"id"`
	ReportID      string               `json:"report_id"`
	VersionNumber int                  `json:"version_number"`
	Content       domain.ReportContent `json:"content"`
	SavedAt       time.Time            `json:"saved_at"`
	SavedBy       string               `json:"saved_by"`
	Comment       string               `json:"comment"`
}

func ToVersionResponse(version *domain.ReportVersion) VersionResponse {
	return VersionResponse{
		ID:            version.ID.String(),
		ReportID:      version.ReportID.String(),
		VersionNumber: version.VersionNumber,
		Content:       version.Content,
		SavedAt:       version.SavedAt,
		SavedBy:       version.SavedBy.String(),
		Comment:       version.Comment,
	}
}

type ICD10Response struct {
	Code         string `json:"code"`
	Description  string `json:"description"`
	Category     string `json:"category"`
}

func ToICD10Response(ref domain.ICD10Reference) ICD10Response {
	return ICD10Response{
		Code:        ref.Code,
		Description: ref.DescriptionRO,
		Category:    ref.Category,
	}
}

type MedicationResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ActiveSubstance string `json:"active_substance"`
	Form            string `json:"form"`
	Dosage          string `json:"dosage"`
	Manufacturer    string `json:"manufacturer"`
}

func ToMedicationResponse(ref domain.MedicationReference) MedicationResponse {
	return MedicationResponse{
		ID:              ref.ID.String(),
		Name:            ref.Name,
		ActiveSubstance: ref.ActiveSubstance,
		Form:            ref.Form,
		Dosage:          ref.Dosage,
		Manufacturer:    ref.Manufacturer,
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Helper functions
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

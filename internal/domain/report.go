package domain

import (
	"time"

	"github.com/google/uuid"
)

// Report represents a medical discharge report
type Report struct {
	ID               uuid.UUID     `json:"id"`
	HospitalID       uuid.UUID     `json:"hospital_id"`
	PatientCNP       string        `json:"patient_cnp"`
	PatientFirstName string        `json:"patient_first_name"`
	PatientLastName  string        `json:"patient_last_name"`
	Specialty        Specialty     `json:"specialty"`
	ReportType       ReportType    `json:"report_type"`
	Status           Status        `json:"status"`
	Content          ReportContent `json:"content"`
	CreatedBy        uuid.UUID     `json:"created_by"`
	CreatedAt        time.Time     `json:"created_at"`
	LastModified     time.Time     `json:"last_modified"`
	FinalizedAt      *time.Time    `json:"finalized_at,omitempty"`
}

// NewReport creates a new report in draft status
func NewReport(hospitalID uuid.UUID, patientCNP, firstName, lastName string, specialty Specialty, reportType ReportType, doctorID uuid.UUID) *Report {
	now := time.Now()
	return &Report{
		ID:               uuid.New(),
		HospitalID:       hospitalID,
		PatientCNP:       patientCNP,
		PatientFirstName: firstName,
		PatientLastName:  lastName,
		Specialty:        specialty,
		ReportType:       reportType,
		Status:           StatusDraft,
		Content:          ReportContent{},
		CreatedBy:        doctorID,
		CreatedAt:        now,
		LastModified:     now,
	}
}

// Specialty represents medical specialty
type Specialty string

const (
	SpecialtyInternalMedicine Specialty = "internal_medicine"
	SpecialtyCardiology       Specialty = "cardiology"
	SpecialtyNeurology        Specialty = "neurology"
	SpecialtyPediatrics       Specialty = "pediatrics"
	SpecialtySurgery          Specialty = "surgery"
)

// ReportType represents type of medical report
type ReportType string

const (
	ReportTypeDischargeSummary ReportType = "discharge_summary"
	ReportTypeTransferSummary  ReportType = "transfer_summary"
	ReportTypeOperativeNote    ReportType = "operative_note"
)

// Status represents report workflow status
type Status string

const (
	StatusDraft     Status = "draft"
	StatusInReview  Status = "in_review"
	StatusApproved  Status = "approved"
	StatusSigned    Status = "signed"
	StatusCancelled Status = "cancelled"
)

// ValidTransitions defines allowed status changes
var ValidTransitions = map[Status][]Status{
	StatusDraft:     {StatusInReview, StatusCancelled},
	StatusInReview:  {StatusDraft, StatusApproved, StatusCancelled},
	StatusApproved:  {StatusSigned, StatusDraft},
	StatusSigned:    {},
	StatusCancelled: {},
}

func (r *Report) CanTransitionTo(newStatus Status) bool {
	allowed, exists := ValidTransitions[r.Status]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == newStatus {
			return true
		}
	}
	return false
}

// ReportVersion represents an immutable snapshot
type ReportVersion struct {
	ID            uuid.UUID     `json:"id"`
	ReportID      uuid.UUID     `json:"report_id"`
	VersionNumber int           `json:"version_number"`
	Content       ReportContent `json:"content"`
	SavedAt       time.Time     `json:"saved_at"`
	SavedBy       uuid.UUID     `json:"saved_by"`
	Comment       string        `json:"comment,omitempty"`
}

// NewReportVersion creates a new version
func NewReportVersion(reportID uuid.UUID, versionNumber int, content ReportContent, userID uuid.UUID, comment string) *ReportVersion {
	return &ReportVersion{
		ID:            uuid.New(),
		ReportID:      reportID,
		VersionNumber: versionNumber,
		Content:       content,
		SavedAt:       time.Now(),
		SavedBy:       userID,
		Comment:       comment,
	}
}

// Reference data types
type ICD10Reference struct {
	Code          string `json:"code"`
	DescriptionRO string `json:"description_ro"`
	Category      string `json:"category"`
}

type MedicationReference struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	ActiveSubstance string    `json:"active_substance"`
	Form            string    `json:"form"`
	Dosage          string    `json:"dosage"`
	Manufacturer    string    `json:"manufacturer"`
}

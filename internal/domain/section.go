package domain

import "time"

// SectionType represents different sections of a medical report
type SectionType string

const (
	SectionPatientData     SectionType = "patient_data"
	SectionAnamnesis       SectionType = "anamnesis"
	SectionExamination     SectionType = "examination"
	SectionLabResults      SectionType = "lab_results"
	SectionDiagnosis       SectionType = "diagnosis"
	SectionTreatment       SectionType = "treatment"
	SectionRecommendations SectionType = "recommendations"
)

// PatientDataSection contains patient demographics
type PatientDataSection struct {
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	CNP            string    `json:"cnp"`
	BirthDate      time.Time `json:"birth_date"`
	Department     string    `json:"department"`
	Ward           string    `json:"ward"`
	Bed            string    `json:"bed"`
	AdmissionDate  time.Time `json:"admission_date"`
	DischargeDate  time.Time `json:"discharge_date"`
}

func (s PatientDataSection) Validate() error {
	if s.FirstName == "" || s.LastName == "" {
		return ErrEmptyField
	}
	if len(s.CNP) != 13 {
		return ErrInvalidCNP
	}
	if s.DischargeDate.Before(s.AdmissionDate) {
		return ErrInvalidDate
	}
	return nil
}

// AnamnesisSection contains medical history
type AnamnesisSection struct {
	ChiefComplaint         string `json:"chief_complaint"`
	HistoryOfPresentIllness string `json:"history_of_present_illness"`
	PastMedicalHistory     string `json:"past_medical_history"`
	Allergies              string `json:"allergies"`
	SocialHistory          string `json:"social_history"`
}

func (s AnamnesisSection) Validate() error {
	if s.ChiefComplaint == "" {
		return ErrEmptyField
	}
	return nil
}

// ExaminationSection contains physical examination findings
type ExaminationSection struct {
	GeneralCondition string     `json:"general_condition"`
	Consciousness    string     `json:"consciousness"`
	VitalSigns       VitalSigns `json:"vital_signs"`
	SystemsReview    string     `json:"systems_review"`
}

type VitalSigns struct {
	BloodPressure    string  `json:"blood_pressure"`
	HeartRate        int     `json:"heart_rate"`
	Temperature      float64 `json:"temperature"`
	RespiratoryRate  int     `json:"respiratory_rate"`
	OxygenSaturation int     `json:"oxygen_saturation"`
}

func (s ExaminationSection) Validate() error {
	return nil
}

// LabResultsSection contains laboratory and imaging results
type LabResultsSection struct {
	LaboratoryTests []LabTest `json:"laboratory_tests"`
	ImagingStudies  []Imaging `json:"imaging_studies"`
}

type LabTest struct {
	Name   string `json:"name"`
	Result string `json:"result"`
	Unit   string `json:"unit"`
	Date   time.Time `json:"date"`
}

type Imaging struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

func (s LabResultsSection) Validate() error {
	return nil
}

// DiagnosisSection contains diagnoses
type DiagnosisSection struct {
	PrimaryDiagnosis       ICD10Code   `json:"primary_diagnosis"`
	SecondaryDiagnoses     []ICD10Code `json:"secondary_diagnoses"`
	ClinicalObservations   string      `json:"clinical_observations"`
}

type ICD10Code struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (s DiagnosisSection) Validate() error {
	if s.PrimaryDiagnosis.Code == "" {
		return ErrInvalidDiagnosis
	}
	return nil
}

// TreatmentSection contains treatment information
type TreatmentSection struct {
	Medications []Medication `json:"medications"`
	Procedures  []Procedure  `json:"procedures"`
}

type Medication struct {
	Name      string     `json:"name"`
	Dosage    string     `json:"dosage"`
	Frequency string     `json:"frequency"`
	Route     string     `json:"route"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

type Procedure struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PerformedAt time.Time `json:"performed_at"`
}

func (s TreatmentSection) Validate() error {
	return nil
}

// RecommendationsSection contains discharge instructions
type RecommendationsSection struct {
	DischargePlan      string `json:"discharge_plan"`
	Medications        string `json:"medications"`
	FollowUp           string `json:"follow_up"`
	DietRestrictions   string `json:"diet_restrictions"`
	ActivityRestrictions string `json:"activity_restrictions"`
}

func (s RecommendationsSection) Validate() error {
	return nil
}

// ReportContent holds all sections
type ReportContent struct {
	PatientData     PatientDataSection     `json:"patient_data"`
	Anamnesis       AnamnesisSection       `json:"anamnesis"`
	Examination     ExaminationSection     `json:"examination"`
	LabResults      LabResultsSection      `json:"lab_results"`
	Diagnosis       DiagnosisSection       `json:"diagnosis"`
	Treatment       TreatmentSection       `json:"treatment"`
	Recommendations RecommendationsSection `json:"recommendations"`
}

func (c ReportContent) Validate() error {
	if err := c.PatientData.Validate(); err != nil {
		return err
	}
	if err := c.Anamnesis.Validate(); err != nil {
		return err
	}
	if err := c.Diagnosis.Validate(); err != nil {
		return err
	}
	return nil
}

func (c ReportContent) IsComplete() bool {
	return c.Validate() == nil
}

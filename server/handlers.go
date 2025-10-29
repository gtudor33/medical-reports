package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tudormiron/medical-reports/internal/domain"
	"github.com/tudormiron/medical-reports/internal/services"
)

type Handlers struct {
	reportService    *services.ReportService
	referenceService *services.ReferenceService
	authService      *services.AuthService
}

func NewHandlers(reportService *services.ReportService, referenceService *services.ReferenceService, authService *services.AuthService) *Handlers {
	return &Handlers{
		reportService:    reportService,
		referenceService: referenceService,
		authService:      authService,
	}
}

// Health check
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "medical-reports-api",
	})
}

// CreateReport creates a new report
func (h *Handlers) CreateReport(c *gin.Context) {
	var req CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	hospitalID, err := ParseUUID(req.HospitalID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_hospital_id",
			Message: "Invalid hospital ID format",
		})
		return
	}

	doctorID, err := ParseUUID(req.DoctorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_doctor_id",
			Message: "Invalid doctor ID format",
		})
		return
	}

	report, err := h.reportService.CreateReport(
		c.Request.Context(),
		hospitalID,
		req.PatientCNP,
		req.PatientFirstName,
		req.PatientLastName,
		domain.Specialty(req.Specialty),
		domain.ReportType(req.ReportType),
		doctorID,
	)

	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, ToReportResponse(report))
}

// GetReport retrieves a report by ID
func (h *Handlers) GetReport(c *gin.Context) {
	reportID, err := ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_report_id",
			Message: "Invalid report ID format",
		})
		return
	}

	report, err := h.reportService.GetReport(c.Request.Context(), reportID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, ToReportResponse(report))
}

// ListReports lists reports with filtering
func (h *Handlers) ListReports(c *gin.Context) {
	doctorID, err := ParseUUID(c.Query("doctor_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_doctor_id",
			Message: "Invalid doctor ID format",
		})
		return
	}

	status := domain.Status(c.Query("status"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	reports, err := h.reportService.ListReports(c.Request.Context(), doctorID, status, limit, offset)
	if err != nil {
		h.handleError(c, err)
		return
	}

	reportResponses := make([]ReportResponse, len(reports))
	for i, report := range reports {
		reportResponses[i] = ToReportResponse(report)
	}

	c.JSON(http.StatusOK, ReportListResponse{
		Reports: reportResponses,
		Total:   len(reportResponses),
		Limit:   limit,
		Offset:  offset,
	})
}

// UpdateReportContent updates report content
func (h *Handlers) UpdateReportContent(c *gin.Context) {
	reportID, err := ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_report_id",
			Message: "Invalid report ID format",
		})
		return
	}

	var req UpdateReportContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	userID, err := ParseUUID(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_user_id",
			Message: "Invalid user ID format",
		})
		return
	}

	if err := h.reportService.UpdateReportContent(c.Request.Context(), reportID, req.Content, userID); err != nil {
		h.handleError(c, err)
		return
	}

	report, err := h.reportService.GetReport(c.Request.Context(), reportID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, ToReportResponse(report))
}

// UpdateReportStatus updates report status
func (h *Handlers) UpdateReportStatus(c *gin.Context) {
	reportID, err := ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_report_id",
			Message: "Invalid report ID format",
		})
		return
	}

	var req UpdateReportStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	if err := h.reportService.UpdateReportStatus(c.Request.Context(), reportID, domain.Status(req.Status)); err != nil {
		h.handleError(c, err)
		return
	}

	report, err := h.reportService.GetReport(c.Request.Context(), reportID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, ToReportResponse(report))
}

// DeleteReport deletes a report
func (h *Handlers) DeleteReport(c *gin.Context) {
	reportID, err := ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_report_id",
			Message: "Invalid report ID format",
		})
		return
	}

	if err := h.reportService.DeleteReport(c.Request.Context(), reportID); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetReportVersions retrieves all versions of a report
func (h *Handlers) GetReportVersions(c *gin.Context) {
	reportID, err := ParseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_report_id",
			Message: "Invalid report ID format",
		})
		return
	}

	versions, err := h.reportService.GetReportVersions(c.Request.Context(), reportID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	versionResponses := make([]VersionResponse, len(versions))
	for i, version := range versions {
		versionResponses[i] = ToVersionResponse(version)
	}

	c.JSON(http.StatusOK, versionResponses)
}

// SearchICD10 searches for ICD-10 codes
func (h *Handlers) SearchICD10(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "missing_query",
			Message: "Query parameter 'q' is required",
		})
		return
	}

	results, err := h.referenceService.SearchICD10(c.Request.Context(), query)
	if err != nil {
		h.handleError(c, err)
		return
	}

	responses := make([]ICD10Response, len(results))
	for i, ref := range results {
		responses[i] = ToICD10Response(ref)
	}

	c.JSON(http.StatusOK, responses)
}

// SearchMedications searches for medications
func (h *Handlers) SearchMedications(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "missing_query",
			Message: "Query parameter 'q' is required",
		})
		return
	}

	results, err := h.referenceService.SearchMedications(c.Request.Context(), query)
	if err != nil {
		h.handleError(c, err)
		return
	}

	responses := make([]MedicationResponse, len(results))
	for i, ref := range results {
		responses[i] = ToMedicationResponse(ref)
	}

	c.JSON(http.StatusOK, responses)
}

// handleError handles domain errors and converts them to HTTP responses
func (h *Handlers) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrReportNotFound):
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "report_not_found",
			Message: "Report not found",
		})
	case errors.Is(err, domain.ErrCannotEditNonDraft):
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "cannot_edit_non_draft",
			Message: "Can only edit reports in draft status",
		})
	case errors.Is(err, domain.ErrInvalidStatusTransition):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_status_transition",
			Message: "Invalid status transition",
		})
	case errors.Is(err, domain.ErrIncompleteReport):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "incomplete_report",
			Message: "Report is incomplete and cannot be finalized",
		})
	case errors.Is(err, domain.ErrInvalidCNP):
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_cnp",
			Message: "Invalid CNP format",
		})
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_server_error",
			Message: "An internal server error occurred",
		})
	}
}

// Register handles user registration
func (h *Handlers) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	authResp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "email_exists",
				Message: "Email already registered",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "registration_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, authResp)
}

// Login handles user login
func (h *Handlers) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	authResp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "authentication_failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, authResp)
}

package domain

import "errors"

var (
	// Report errors
	ErrReportNotFound              = errors.New("report not found")
	ErrCannotEditNonDraft          = errors.New("can only edit reports in draft status")
	ErrCannotModifySignedReport    = errors.New("cannot modify signed reports")
	ErrInvalidStatusTransition     = errors.New("invalid status transition")
	ErrIncompleteReport            = errors.New("report is incomplete")
	ErrInvalidCNP                  = errors.New("invalid CNP format")
	
	// Validation errors
	ErrEmptyField                  = errors.New("required field is empty")
	ErrInvalidDate                 = errors.New("invalid date")
	ErrInvalidDiagnosis            = errors.New("invalid diagnosis code")
	
	// Repository errors
	ErrDatabaseConnection          = errors.New("database connection error")
	ErrDatabaseQuery               = errors.New("database query error")
)

package ierr

import (
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
)

//
// ─── ERROR CODES ────────────────────────────────────────────────────────────────
//

const (
	// Validation & Input
	ErrCodeValidation       = "validation_error"
	ErrCodeInvalidOperation = "invalid_operation"
	ErrCodeBadRequest       = "bad_request"

	// Auth
	ErrCodePermissionDenied = "permission_denied"
	ErrCodeUnauthorized     = "unauthorized"

	// Database
	ErrCodeNotFound        = "not_found"
	ErrCodeAlreadyExists   = "already_exists"
	ErrCodeVersionConflict = "version_conflict"
	ErrCodeConflict        = "conflict"

	// System/Internal
	ErrCodeSystemError   = "system_error"
	ErrCodeInternalError = "internal_error"
	ErrCodeTimeout       = "timeout"
	ErrCodeHTTPClient    = "http_client_error"
	ErrCodeDatabase      = "database_error"

	// Integration/External
	ErrCodeIntegration = "integration_error"

	// File Upload
	ErrCodeFileTooLarge     = "file_too_large"
	ErrCodeInvalidExtension = "invalid_extension"
)

//
// ─── ERROR DECLARATIONS ─────────────────────────────────────────────────────────
//

// Validation & Input
var (
	ErrValidation       = new(ErrCodeValidation, "validation error")
	ErrInvalidOperation = new(ErrCodeInvalidOperation, "invalid operation")
	ErrBadRequest       = new(ErrCodeBadRequest, "bad request")
)

// Auth
var (
	ErrPermissionDenied = new(ErrCodePermissionDenied, "permission denied")
	ErrUnauthorized     = new(ErrCodeUnauthorized, "unauthorized")
)

// Database
var (
	ErrNotFound        = new(ErrCodeNotFound, "resource not found")
	ErrAlreadyExists   = new(ErrCodeAlreadyExists, "resource already exists")
	ErrVersionConflict = new(ErrCodeVersionConflict, "version conflict")
	ErrConflict        = new(ErrCodeConflict, "conflict")
)

// System/Internal
var (
	ErrSystem     = new(ErrCodeSystemError, "system error")
	ErrInternal   = new(ErrCodeInternalError, "internal error")
	ErrTimeout    = new(ErrCodeTimeout, "operation timed out")
	ErrHTTPClient = new(ErrCodeHTTPClient, "http client error")
	ErrDatabase   = new(ErrCodeDatabase, "database error")
)

// Integration
var (
	ErrIntegration = new(ErrCodeIntegration, "integration error")
)

// File Upload
var (
	ErrFileTooLarge     = new(ErrCodeFileTooLarge, "file too large")
	ErrInvalidExtension = new(ErrCodeInvalidExtension, "invalid file extension")
)

//
// ─── HTTP STATUS MAPPING ────────────────────────────────────────────────────────
//

var statusCodeMap = map[error]int{
	// Validation
	ErrValidation:       http.StatusBadRequest,
	ErrInvalidOperation: http.StatusBadRequest,
	ErrBadRequest:       http.StatusBadRequest,

	// Auth
	ErrPermissionDenied: http.StatusForbidden,
	ErrUnauthorized:     http.StatusUnauthorized,

	// DB
	ErrNotFound:        http.StatusNotFound,
	ErrAlreadyExists:   http.StatusConflict,
	ErrVersionConflict: http.StatusConflict,
	ErrConflict:        http.StatusConflict,

	// System/Internal
	ErrSystem:     http.StatusInternalServerError,
	ErrInternal:   http.StatusInternalServerError,
	ErrTimeout:    http.StatusGatewayTimeout,
	ErrHTTPClient: http.StatusInternalServerError,
	ErrDatabase:   http.StatusInternalServerError,

	// Integration
	ErrIntegration: http.StatusBadGateway,

	// File Upload
	ErrFileTooLarge:     http.StatusRequestEntityTooLarge,
	ErrInvalidExtension: http.StatusBadRequest,
}

func HTTPStatusFromErr(err error) int {
	for e, status := range statusCodeMap {
		if errors.Is(err, e) {
			return status
		}
	}
	return http.StatusInternalServerError
}

//
// ─── INTERNAL ERROR TYPE ────────────────────────────────────────────────────────
//

type InternalError struct {
	Code    string
	Message string
	Op      string
	Err     error
}

func (e *InternalError) Error() string {
	if e.Err == nil {
		return e.DisplayError()
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Err.Error())
}

func (e *InternalError) DisplayError() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *InternalError) Unwrap() error {
	return e.Err
}

func (e *InternalError) Is(target error) bool {
	if target == nil {
		return false
	}
	t, ok := target.(*InternalError)
	if !ok {
		return errors.Is(e.Err, target)
	}
	return e.Code == t.Code
}

func new(code string, message string) *InternalError {
	return &InternalError{Code: code, Message: message}
}

//
// ─── ERROR HELPERS ──────────────────────────────────────────────────────────────
//

func As(err error, target any) bool {
	return errors.As(err, target)
}

// Validation
func IsValidation(err error) bool       { return errors.Is(err, ErrValidation) }
func IsInvalidOperation(err error) bool { return errors.Is(err, ErrInvalidOperation) }

// Auth
func IsPermissionDenied(err error) bool { return errors.Is(err, ErrPermissionDenied) }
func IsUnauthorized(err error) bool     { return errors.Is(err, ErrUnauthorized) }

// DB
func IsNotFound(err error) bool        { return errors.Is(err, ErrNotFound) }
func IsAlreadyExists(err error) bool   { return errors.Is(err, ErrAlreadyExists) }
func IsVersionConflict(err error) bool { return errors.Is(err, ErrVersionConflict) }
func IsConflict(err error) bool        { return errors.Is(err, ErrConflict) }

// System
func IsSystem(err error) bool     { return errors.Is(err, ErrSystem) }
func IsInternal(err error) bool   { return errors.Is(err, ErrInternal) }
func IsTimeout(err error) bool    { return errors.Is(err, ErrTimeout) }
func IsHTTPClient(err error) bool { return errors.Is(err, ErrHTTPClient) }
func IsDatabase(err error) bool   { return errors.Is(err, ErrDatabase) }

// Integration
func IsIntegration(err error) bool { return errors.Is(err, ErrIntegration) }

// File Upload
func IsFileTooLarge(err error) bool     { return errors.Is(err, ErrFileTooLarge) }
func IsInvalidExtension(err error) bool { return errors.Is(err, ErrInvalidExtension) }

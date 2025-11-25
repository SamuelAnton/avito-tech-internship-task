package errors

import "net/http"

// Error codes
type ErrorCode string

const (
	// From specificaation
	ErrorCodeTeamExists  ErrorCode = "TEAM_EXISTS"  // 400
	ErrorCodePRExists    ErrorCode = "PR_EXISTS"    // 409
	ErrorCodePRMerged    ErrorCode = "PR_MERGED"    // 409
	ErrorCodeNotAssigned ErrorCode = "NOT_ASSIGNED" // 409
	ErrorCodeNoCandidate ErrorCode = "NO_CANDIDATE" // 409
	ErrorCodeNotFound    ErrorCode = "NOT_FOUND"    // 404
	// "Basic" cases
	ErrorCodeInvalidInput ErrorCode = "INVALID_INPUT"  // 400
	ErrorCodeInternal     ErrorCode = "INTERNAL_ERROR" // 500
)

// Mapping codes
type ErrorInfo struct {
	Status  int
	Message string
}

var errorMapping = map[ErrorCode]ErrorInfo{
	ErrorCodeTeamExists: {
		Status:  http.StatusBadRequest,
		Message: "Team already exists",
	},
	ErrorCodePRExists: {
		Status:  http.StatusConflict,
		Message: "PR already exists",
	},
	ErrorCodePRMerged: {
		Status:  http.StatusConflict,
		Message: "PR is merged",
	},
	ErrorCodeNotAssigned: {
		Status:  http.StatusConflict,
		Message: "Reviewer not assigned",
	},
	ErrorCodeNoCandidate: {
		Status:  http.StatusConflict,
		Message: "No active candidates available",
	},
	ErrorCodeNotFound: {
		Status:  http.StatusNotFound,
		Message: "Resourse not found",
	},
	ErrorCodeInvalidInput: {
		Status:  http.StatusBadRequest,
		Message: "Invalid input",
	},
	ErrorCodeInternal: {
		Status:  http.StatusInternalServerError,
		Message: "Internal service error",
	},
}

func GetInfo(code ErrorCode) (int, string) {
	if info, exists := errorMapping[code]; exists {
		return info.Status, info.Message
	}
	return http.StatusInternalServerError, "Internal service error"
}

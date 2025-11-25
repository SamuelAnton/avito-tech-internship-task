package handlers

import (
	"PR_reviewer_assign_service/internal/errors"
	"PR_reviewer_assign_service/internal/models"
	"encoding/json"
	"net/http"
)

// Response structs
type TeamResponse struct {
	Team *models.Team `json:"team"`
}

type UserResponse struct {
	User *models.User `json:"user"`
}

type PRResponse struct {
	PR *models.PullRequest `json:"pull_request"`
}

type ErrorRespone struct {
	Error struct {
		Code    errors.ErrorCode `json:"code"`
		Message string           `json:"message"`
	} `json:"error"`
}

type PRReassignReasponse struct {
	PR         *models.PullRequest `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}

type GetReviewResponse struct {
	UserID       string                    `json:"user_id"`
	PullRequests []models.PullRequestShort `json:"pull_requests"`
}

// Response functions
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed, to encode response", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, code errors.ErrorCode) {
	status, message := errors.GetInfo(code)
	writeErrorResponse(w, code, status, message)
}

func writeErrorMessage(w http.ResponseWriter, code errors.ErrorCode, message string) {
	status, _ := errors.GetInfo(code)
	writeErrorResponse(w, code, status, message)
}

func writeErrorResponse(w http.ResponseWriter, code errors.ErrorCode, status int, message string) {
	response := ErrorRespone{}
	response.Error.Code = code
	response.Error.Message = message

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

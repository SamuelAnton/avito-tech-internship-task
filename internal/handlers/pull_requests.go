package handlers

import (
	"PR_reviewer_assign_service/internal/errors"
	"PR_reviewer_assign_service/internal/models"
	"PR_reviewer_assign_service/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

// Handler for pr requests
type PRHandler struct {
	service *service.Service
}

func NewPRHandler(service *service.Service) *PRHandler {
	return &PRHandler{service: service}
}

/*
/pullRequest/create - PullRequestCreateQuery
*/
func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	// Decode input
	var pr models.PullRequestCreateQuery

	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		writeErrorMessage(w, errors.ErrorCodeInvalidInput, "Invalid JSON format")
		return
	}

	// Validate input
	if err, msg := ValidatePullRequestCreateQuery(pr); err != "" {
		writeErrorMessage(w, err, msg)
		return
	}

	// Create pull request
	log.Printf("Creating PR: %s", pr.PullRequestID)
	pullRequest, err := h.service.CreatePullRequest(r.Context(), pr)
	if err != "" {
		writeError(w, err)
		return
	}
	log.Printf("PR created: %s", pr.PullRequestID)

	// Send response
	writeJSON(w, http.StatusCreated, pullRequest)
}

/*
/pullRequest/merge - PullRequestMergeQuery
*/
func (h *PRHandler) Merge(w http.ResponseWriter, r *http.Request) {
	// Decode input
	var pr models.PullRequestMergeQuery

	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		writeErrorMessage(w, errors.ErrorCodeInvalidInput, "Invalid JSON format")
		return
	}

	// Valiedate input
	if err, msg := ValidatePullRequestMergeQuery(pr); err != "" {
		writeErrorMessage(w, err, msg)
		return
	}

	// Merge pr
	log.Printf("merging PR: %s", pr.PullRequestID)
	pullRequest, err := h.service.MergePullRequest(r.Context(), pr)
	if err != "" {
		writeError(w, err)
		return
	}
	log.Printf("PR merged: %s", pr.PullRequestID)

	// Send Response
	writeJSON(w, http.StatusOK, pullRequest)
}

/*
/pullRequest/reassign - PullRequestMergeQuery
*/
func (h *PRHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	// Decode input
	var pr models.PullRequestReassignQuery

	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		writeErrorMessage(w, errors.ErrorCodeInvalidInput, "Invalid JSON format")
		return
	}

	// Valiedate input
	if err, msg := ValidatePullRequestReassignQuery(pr); err != "" {
		writeErrorMessage(w, err, msg)
		return
	}

	// Reassign user
	log.Printf("Reassigning user: %s, in PR: %s", pr.OldUserID, pr.PullRequestID)
	pullRequest, newCandidate, err := h.service.Reassign(r.Context(), pr)
	if err != "" {
		writeError(w, err)
		return
	}
	log.Printf("User: %s, reassigned to: %s, in PR: %s", pr.OldUserID, *newCandidate, pr.PullRequestID)

	// Send Response
	writeJSON(w, http.StatusOK, PRReassignReasponse{
		PR:         pullRequest,
		ReplacedBy: *newCandidate,
	})
}

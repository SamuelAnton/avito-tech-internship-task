package handlers

import (
	"PR_reviewer_assign_service/internal/errors"
	"PR_reviewer_assign_service/internal/models"
	"PR_reviewer_assign_service/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

// Handler for user requests
type UserHandler struct {
	service *service.Service
}

func NewUserHandler(service *service.Service) *UserHandler {
	return &UserHandler{service: service}
}

/*
/users/setIsActive - UserActiveQuery
*/
func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	// Decode input
	var user models.UserActiveQuery

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeErrorMessage(w, errors.ErrorCodeInvalidInput, "Invalid JSON format")
		return
	}

	// Validate input
	if err, msg := ValidateUserActiveQuery(user); err != "" {
		writeErrorMessage(w, err, msg)
		return
	}

	// Update user's isActive
	log.Printf("Updating activness of user: %s", user.UserID)
	u, err := h.service.UserSetIsActive(r.Context(), user.UserID, user.IsActive)
	if err != "" {
		writeError(w, err)
		return
	}
	log.Printf("User's activness is changed: %s", user.UserID)

	// Send response
	writeJSON(w, http.StatusOK, UserResponse{User: u})
}

/*
/users/getReview - UserIDQuery
*/
func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	// Decode input
	var user models.UserIDQuery

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeErrorMessage(w, errors.ErrorCodeInvalidInput, "Invalid JSON format")
		return
	}

	// Validate input
	if err, msg := ValidateUserIDQuery(user); err != "" {
		writeErrorMessage(w, err, msg)
		return
	}

	// Get users PR's
	log.Printf("Receiving user's PRs: %s", user.UserID)
	prs, err := h.service.GetPRs(r.Context(), user.UserID)
	if err != "" {
		writeError(w, err)
		return
	}
	log.Printf("PRs of user: %s, received", user.UserID)

	// Send response
	writeJSON(w, http.StatusOK, GetReviewResponse{
		UserID:       user.UserID,
		PullRequests: prs,
	})
}

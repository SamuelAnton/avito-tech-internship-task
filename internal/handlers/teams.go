package handlers

import (
	"PR_reviewer_assign_service/internal/errors"
	"PR_reviewer_assign_service/internal/models"
	"PR_reviewer_assign_service/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

// Handler for team requests
type TeamHandler struct {
	service *service.Service
}

func NewTeamHandler(service *service.Service) *TeamHandler {
	return &TeamHandler{service: service}
}

/*
/team/add - Team
*/
func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	// Decode input
	var team models.Team

	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		writeErrorMessage(w, errors.ErrorCodeInvalidInput, "Invalid JSON format")
		return
	}

	// Validate input
	if err, msg := ValidateTeam(team); err != "" {
		writeErrorMessage(w, err, msg)
		return
	}

	// Create team
	log.Printf("Creating team: %s", team.TeamName)
	if code := h.service.CreateTeam(r.Context(), &team); code != "" {
		writeError(w, code)
		return
	}
	log.Printf("Team created: %s", team.TeamName)

	// Send response
	writeJSON(w, http.StatusCreated, TeamResponse{Team: &team})
}

/*
	 /team/get - TeamNameQuery
	}
*/
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	// Decode input
	var teamName models.TeamNameQuery

	if err := json.NewDecoder(r.Body).Decode(&teamName); err != nil {
		writeErrorMessage(w, errors.ErrorCodeInvalidInput, "Invalid JSON format")
		return
	}

	// Validate input
	if err, msg := ValidateTeamNameQuery(teamName); err != "" {
		writeErrorMessage(w, err, msg)
		return
	}

	// Get team
	log.Printf("Receiving team: %s", teamName.TeamName)
	team, err := h.service.GetTeam(r.Context(), teamName.TeamName)
	if err != "" {
		writeError(w, err)
		return
	}
	log.Printf("Team received: %s", team.TeamName)

	// Send response
	writeJSON(w, http.StatusOK, team)
}

package models

import "time"

// Query objects
type TeamNameQuery struct {
	TeamName string `json:"team_name"`
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type UserActiveQuery struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type PullRequestCreateQuery struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type PullRequestMergeQuery struct {
	PullRequestID string `json:"pull_request_id"`
}

type PullRequestReassignQuery struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type UserIDQuery struct {
	UserID string `json:"user_id"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

// Query + DB objects
type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PullRequest struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         time.Time  `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

// Statistic models
type UserStats struct {
	UserID           string `json:"user_id"`
	Username         string `json:"username"`
	TeamName         string `json:"team_name"`
	IsActive         bool   `json:"is_active"`
	AssignmentsCount int    `json:"assignments_count"`
}

type UsersStatistics struct {
	TotalUserNumber       int `json:"total_user_number"`
	TotalActiveUserNumber int `json:"total_active_user_number"`
}

type TeamStats struct {
	TeamName                string       `json:"team_name"`
	MembersTotal            int          `json:"members_total"`
	Members                 []TeamMember `json:"members"`
	PullRequestsTotal       int          `json:"pull_requests_total"`
	ActivePullRequestsTotal int          `json:"active_pull_requests_total"`
	PullRequests            []string     `json:"pull_requests_id"`
}

type TeamsStatistics struct {
	TotalTeamNumber int `json:"total_team_number"`
}

type PullRequestStatistics struct {
	TotalPR       int `json:"total_pull_request_number"`
	TotalActivePR int `json:"total_active_pull_request_number"`
}

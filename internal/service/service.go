package service

import (
	"PR_reviewer_assign_service/internal/errors"
	"PR_reviewer_assign_service/internal/models"
	"PR_reviewer_assign_service/internal/storage"
	"context"
	"math/rand"
	"time"
)

// Middle part between Storage and Handlers
type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{storage: storage}
}

// Team functions

// Create new team
func (s *Service) CreateTeam(ctx context.Context, team *models.Team) errors.ErrorCode {
	// Check existance
	t, err := s.storage.GetTeam(ctx, team.TeamName)
	if err != nil {
		return errors.ErrorCodeInternal
	}
	if t != nil {
		return errors.ErrorCodeTeamExists
	}

	// Create team
	if err := s.storage.CreateTeam(ctx, team); err != nil {
		return errors.ErrorCodeInternal
	}

	return ""
}

// Get existing team
func (s *Service) GetTeam(ctx context.Context, teamName string) (*models.Team, errors.ErrorCode) {
	team, err := s.storage.GetTeam(ctx, teamName)
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}
	if team == nil {
		return nil, errors.ErrorCodeNotFound
	}
	return team, ""
}

// User functions
// Set User isActive
func (s *Service) UserSetIsActive(ctx context.Context, id string, active bool) (*models.User, errors.ErrorCode) {
	// Get user
	user, err := s.storage.GetUser(ctx, id)
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}
	if user == nil {
		return nil, errors.ErrorCodeNotFound
	}
	// Update user
	user.IsActive = active
	if err := s.storage.UpdateUser(ctx, user); err != nil {
		return nil, errors.ErrorCodeInternal
	}
	return user, ""
}

// Get User's prs
func (s *Service) GetPRs(ctx context.Context, id string) ([]models.PullRequestShort, errors.ErrorCode) {
	// Check user existance
	user, err := s.storage.GetUser(ctx, id)
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}
	if user == nil {
		return nil, errors.ErrorCodeNotFound
	}
	// Get user's PRs
	prs, err := s.storage.GetPRsByRewiever(ctx, id)
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}
	return prs, ""
}

// PR functions
// Create new Pull Request
func (s *Service) CreatePullRequest(ctx context.Context, prQuery models.PullRequestCreateQuery) (*models.PullRequest, errors.ErrorCode) {
	// Check pr existance
	pr, err := s.storage.GetPR(ctx, prQuery.PullRequestID)
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}
	if pr != nil {
		return nil, errors.ErrorCodePRExists
	}

	// Check author existance
	user, err := s.storage.GetUser(ctx, prQuery.AuthorID)
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}
	if user == nil {
		return nil, errors.ErrorCodeNotFound
	}

	// Get reviewers
	reviewers, er := s.assignReviewers(ctx, *user, 2, nil)
	if er != "" {
		return nil, er
	}

	// Create pr
	pr, err = s.storage.CreatePR(ctx, &models.PullRequest{
		PullRequestID:     prQuery.PullRequestID,
		PullRequestName:   prQuery.PullRequestName,
		AuthorID:          prQuery.AuthorID,
		Status:            "OPEN",
		AssignedReviewers: reviewers,
		CreatedAt:         time.Now(),
		MergedAt:          nil,
	})
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}
	return pr, ""
}

// Assign reviewers to the pr
func (s *Service) assignReviewers(ctx context.Context, author models.User, limit int, excludes *[]string) ([]string, errors.ErrorCode) {
	// Get active users
	users, err := s.storage.GetActiveUsersInTeam(ctx, author.TeamName)
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}

	// Exclude not acceptable ones
	var candidates []*models.User
	for _, user := range users {
		if user.UserID == author.UserID {
			continue
		}
		if excludes != nil {
			excluded := false
			for _, exclude := range *excludes {
				if user.UserID == exclude {
					excluded = true
					break
				}
			}
			if excluded {
				continue
			}
		}
		candidates = append(candidates, user)
	}

	// Shuffle them to choose randomly
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	// Choose reviewers
	limit = min(limit, len(candidates))

	reviewers := make([]string, limit)
	i := 0
	for _, candidate := range candidates {
		if i == limit {
			break
		}
		reviewers[i] = candidate.UserID
		i++
	}

	return reviewers, ""
}

// Merge pr
func (s *Service) MergePullRequest(ctx context.Context, prQuery models.PullRequestMergeQuery) (*models.PullRequest, errors.ErrorCode) {
	// Check pr existance
	pr, err := s.storage.GetPR(ctx, prQuery.PullRequestID)
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}
	if pr == nil {
		return nil, errors.ErrorCodeNotFound
	}

	// Merge pr
	pr.Status = "MERGED"
	mergedAt := time.Now()
	pr.MergedAt = &mergedAt

	if err := s.storage.UpdatePR(ctx, pr); err != nil {
		return nil, errors.ErrorCodeInternal
	}

	return pr, ""
}

// Reassign user in pr
func (s *Service) Reassign(ctx context.Context, query models.PullRequestReassignQuery) (*models.PullRequest, *string, errors.ErrorCode) {
	// Check pr existance
	pr, err := s.storage.GetPR(ctx, query.PullRequestID)
	if err != nil {
		return nil, nil, errors.ErrorCodeInternal
	}
	if pr == nil {
		return nil, nil, errors.ErrorCodeNotFound
	}

	// Check user existance
	oldUser, err := s.storage.GetUser(ctx, query.OldUserID)
	if err != nil {
		return nil, nil, errors.ErrorCodeInternal
	}
	if oldUser == nil {
		return nil, nil, errors.ErrorCodeNotFound
	}

	// Check Merged
	if pr.Status == "MERGED" {
		return nil, nil, errors.ErrorCodePRMerged
	}

	// Check Assigned
	assigned := false
	for _, reviewer := range pr.AssignedReviewers {
		if reviewer == oldUser.UserID {
			assigned = true
			break
		}
	}
	if !assigned {
		return nil, nil, errors.ErrorCodeNotAssigned
	}

	// Choose new candidate
	var excludes []string
	excludes = append(excludes, pr.AssignedReviewers...)
	excludes = append(excludes, pr.AuthorID)
	candidates, er := s.assignReviewers(ctx, *oldUser, 1, &excludes)
	if er != "" {
		return nil, nil, er
	}
	if len(candidates) == 0 {
		return nil, nil, errors.ErrorCodeNoCandidate
	}

	// Assign new candidate
	for i := 0; i < len(pr.AssignedReviewers); i++ {
		if pr.AssignedReviewers[i] == oldUser.UserID {
			pr.AssignedReviewers[i] = candidates[0]
			break
		}
	}

	// Update pr
	if err := s.storage.UpdatePR(ctx, pr); err != nil {
		return nil, nil, errors.ErrorCodeInternal
	}
	return pr, &candidates[0], ""
}

// Additional functions
func (s *Service) GetUsersStatistics(ctx context.Context) (*models.UsersStatistics, errors.ErrorCode) {
	statistics, err := s.storage.GetUsersStatistics(ctx)
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}
	return statistics, ""
}

func (s *Service) GetUserStatistics(ctx context.Context, id string) (*models.UserStats, errors.ErrorCode) {
	user, err := s.storage.GetUser(ctx, id)
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}

	assignments_count, err := s.storage.GetUserStatistics(ctx, id)
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}
	return &models.UserStats{
		UserID:           user.UserID,
		Username:         user.Username,
		TeamName:         user.TeamName,
		IsActive:         user.IsActive,
		AssignmentsCount: assignments_count,
	}, ""
}

func (s *Service) GetTeamsStatistics(ctx context.Context) (*models.TeamsStatistics, errors.ErrorCode) {
	statistics, err := s.storage.GetTeamsStatistics(ctx)
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}
	return statistics, ""
}

func (s *Service) GetTeamStatistics(ctx context.Context, name string) (*models.TeamStats, errors.ErrorCode) {
	statistics, err := s.storage.GetTeamStatistics(ctx, name)
	if err != nil {
		return nil, errors.ErrorCodeInternal
	}
	if statistics == nil {
		return nil, errors.ErrorCodeNotFound
	}
	return statistics, ""
}

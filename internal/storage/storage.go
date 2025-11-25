package storage

import (
	"PR_reviewer_assign_service/internal/models"
	"context"
)

// Interface for different types of storage (possibly not just PostgreSQL)
type Storage interface {
	// Team functions
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
	GetActiveUsersInTeam(ctx context.Context, team_name string) ([]*models.User, error)

	// User functions
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, userId string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error

	// Pull Request functions
	CreatePR(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error)
	GetPR(ctx context.Context, prId string) (*models.PullRequest, error)
	UpdatePR(ctx context.Context, pr *models.PullRequest) error
	GetPRsByRewiever(ctx context.Context, userId string) ([]models.PullRequestShort, error)

	Close() error

	// Additional functions
	GetUsersStatistics(ctx context.Context) (*models.UsersStatistics, error)
	GetUserStatistics(ctx context.Context, id string) (int, error)
	GetTeamsStatistics(ctx context.Context) (*models.TeamsStatistics, error)
	GetTeamStatistics(ctx context.Context, name string) (*models.TeamStats, error)
}

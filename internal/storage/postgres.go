package storage

import (
	"PR_reviewer_assign_service/internal/models"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connection string) (*PostgresStorage, error) {
	// Open connection
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Initialize tables
	if err := initTables(db); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return &PostgresStorage{db: db}, nil
}

//go:embed schema.sql
var schemaSQL string

func initTables(db *sql.DB) error {
	if _, err := db.Exec(schemaSQL); err != nil {
		return err
	}
	log.Println("Database schema initialized successfully")
	return nil
}

func (p *PostgresStorage) Close() error {
	return p.db.Close()
}

// Team functions
func (p *PostgresStorage) CreateTeam(ctx context.Context, team *models.Team) error {
	// Start a transaction
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert team
	_, err = tx.ExecContext(ctx, `
		INSERT INTO teams (team_name) VALUES ($1)
	`, team.TeamName)
	if err != nil {
		return err
	}

	// Insert / update users
	for _, member := range team.Members {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO users (user_id, username, team_name, is_active)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id)
			DO UPDATE SET username = $2, team_name = $3, is_active = $4
		`, member.UserID, member.Username, team.TeamName, member.IsActive)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (p *PostgresStorage) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	// Check existance
	var teamExists bool
	err := p.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)",
		teamName,
	).Scan(&teamExists)

	if err != nil {
		return nil, err
	}
	if !teamExists {
		return nil, nil
	}

	// Get team members
	var team models.Team
	team.TeamName = teamName

	rows, err := p.db.QueryContext(ctx, `
		SELECT user_id, username, is_active
		FROM users
		WHERE team_name = $1
	`, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Make TeamMembers
	for rows.Next() {
		var member models.TeamMember
		if err := rows.Scan(&member.UserID, &member.Username, &member.IsActive); err != nil {
			return nil, err
		}
		team.Members = append(team.Members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &team, nil
}

func (p *PostgresStorage) GetActiveUsersInTeam(ctx context.Context, teamName string) ([]*models.User, error) {
	// Get active users in team
	rows, err := p.db.QueryContext(ctx, `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE team_name = $1 AND is_active = true	
	`, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create user objects
	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// User functions
func (p *PostgresStorage) CreateUser(ctx context.Context, user *models.User) error {
	// Create user
	_, err := p.db.ExecContext(ctx, `
		INSERT INTO users  (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
	`, user.UserID, user.Username, user.TeamName, user.IsActive)

	return err
}

func (p *PostgresStorage) GetUser(ctx context.Context, userId string) (*models.User, error) {
	// Get user
	var user models.User
	err := p.db.QueryRowContext(ctx, `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE user_id = $1
	`, userId).Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PostgresStorage) UpdateUser(ctx context.Context, user *models.User) error {
	// Update user
	_, err := p.db.ExecContext(ctx, `
		UPDATE users
		SET username = $1, team_name = $2, is_active = $3
		WHERE user_id = $4
	`, user.Username, user.TeamName, user.IsActive, user.UserID)
	return err
}

// PR functions
func (p *PostgresStorage) CreatePR(ctx context.Context, pr *models.PullRequest) (*models.PullRequest, error) {
	// Create transaction
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert pr
	_, err = tx.ExecContext(ctx, `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, pr.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Insert reviewers
	for _, reviewer := range pr.AssignedReviewers {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO pr_reviewers (pr_id, user_id)
			VALUES ($1, $2)
			`, pr.PullRequestID, reviewer)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return pr, nil
}

func (p *PostgresStorage) GetPR(ctx context.Context, prId string) (*models.PullRequest, error) {
	var pr models.PullRequest
	var mergedAt sql.NullTime

	// Get pr
	err := p.db.QueryRowContext(ctx, `
		SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE pull_request_id = $1
	`, prId).Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	// Get reviewers
	rows, err := p.db.QueryContext(ctx, `
		SELECT user_id FROM pr_reviewers WHERE pr_id = $1
	`, prId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Add reviewers to pr
	for rows.Next() {
		var reviewer string
		if err := rows.Scan(&reviewer); err != nil {
			return nil, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, reviewer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pr, nil
}

func (p *PostgresStorage) UpdatePR(ctx context.Context, pr *models.PullRequest) error {
	// Create transaction
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update pr
	_, err = tx.ExecContext(ctx, `
		UPDATE pull_requests
		SET status = $1, merged_at = $2
		WHERE pull_request_id = $3
	`, pr.Status, pr.MergedAt, pr.PullRequestID)
	if err != nil {
		return err
	}

	// Update reviewers
	_, err = tx.ExecContext(ctx, `
		DELETE FROM pr_reviewers WHERE pr_id = $1
	`, pr.PullRequestID)
	if err != nil {
		return err
	}

	for _, reviewer := range pr.AssignedReviewers {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO pr_reviewers (pr_id, user_id)
			VALUES ($1, $2)
		`, pr.PullRequestID, reviewer)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (p *PostgresStorage) GetPRsByRewiever(ctx context.Context, userId string) ([]models.PullRequestShort, error) {
	rows, err := p.db.QueryContext(ctx, `
		SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status
		FROM pull_requests pr
		JOIN pr_reviewers prr ON pr.pull_request_id = prr.pr_id
		WHERE prr.user_id = $1 AND pr.status != 'MERGED'
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []models.PullRequestShort
	for rows.Next() {
		var pr models.PullRequestShort
		if err := rows.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}

// Additional functions
func (p *PostgresStorage) GetUsersStatistics(ctx context.Context) (*models.UsersStatistics, error) {
	var statistics models.UsersStatistics

	err := p.db.QueryRowContext(ctx, `
		SELECT 
		(SELECT Count(*) FROM users) as total_user_number,
		(SELECT Count(is_active) FROM users WHERE is_active = true) as total_active_user_number
	`).Scan(&statistics.TotalUserNumber, &statistics.TotalActiveUserNumber)

	if err != nil {
		return nil, err
	}

	return &statistics, nil
}

func (p *PostgresStorage) GetUserStatistics(ctx context.Context, id string) (int, error) {
	var assignments_count int

	err := p.db.QueryRowContext(ctx, `
		SELECT COUNT(pr_id) as assignments_count
		FROM pr_reviewers
		WHERE user_id = $1
	`, id).Scan(&assignments_count)
	if err != nil {
		return 0, err
	}

	return assignments_count, nil
}

func (p *PostgresStorage) GetTeamsStatistics(ctx context.Context) (*models.TeamsStatistics, error) {
	var statistics models.TeamsStatistics

	err := p.db.QueryRowContext(ctx, `
		SELECT Count(*) as total_team_number FROM teams
	`).Scan(&statistics.TotalTeamNumber)

	if err != nil {
		return nil, err
	}

	return &statistics, nil
}

func (p *PostgresStorage) GetTeamStatistics(ctx context.Context, name string) (*models.TeamStats, error) {
	// Get team
	team, err := p.GetTeam(ctx, name)
	if err != nil {
		return nil, err
	}

	if team == nil {
		return nil, nil
	}

	// Count team members
	statistics := &models.TeamStats{
		TeamName:     team.TeamName,
		Members:      team.Members,
		MembersTotal: len(team.Members),
	}

	// Count all prs
	err = p.db.QueryRowContext(ctx, `
		SELECT Count(*)
		FROM pull_requests pr
		JOIN users u ON pr.author_id = u.user_id
		WHERE u.team_name = $1
	`, name).Scan(&statistics.PullRequestsTotal)
	if err != nil {
		return nil, err
	}

	// Count active prs
	err = p.db.QueryRowContext(ctx, `
		SELECT Count(*)
		FROM pull_requests pr
		JOIN users u ON pr.author_id = u.user_id
		WHERE u.team_name = $1 AND pr.status = 'OPEN'
	`, name).Scan(&statistics.ActivePullRequestsTotal)
	if err != nil {
		return nil, err
	}

	// Get snippet of teams's prs ids
	rows, err := p.db.QueryContext(ctx, `
		SELECT pr.pull_request_id
		FROM pull_requests pr
		JOIN users u ON pr.author_id = u.user_id
		WHERE u.team_name = $1
		LIMIT 20
	`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var prID string
		if err := rows.Scan(&prID); err != nil {
			return nil, err
		}
		statistics.PullRequests = append(statistics.PullRequests, prID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return statistics, nil
}

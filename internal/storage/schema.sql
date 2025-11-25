CREATE TABLE IF NOT EXISTS teams (
    team_name VARCHAR(100) PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(50) PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    team_name VARCHAR(100) NOT NULL REFERENCES teams(team_name) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id VARCHAR(50) PRIMARY KEY,
    pull_request_name VARCHAR(150) NOT NULL,
    author_id VARCHAR(50) NOT NULL REFERENCES users(user_id),
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP NOT NULL,
    merged_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pr_reviewers (
    pr_id VARCHAR(50) NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    user_id VARCHAR(50) NOT NULL REFERENCES users(user_id),
    PRIMARY KEY (pr_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_users_team_active ON users(team_name, is_active);
CREATE INDEX IF NOT EXISTS idx_prs_status ON pull_requests(status);
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_user ON pr_reviewers(user_id);
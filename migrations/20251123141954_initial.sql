-- +goose Up
-- +goose StatementBegin
CREATE TABLE teams (
    team_name TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    user_id   TEXT PRIMARY KEY,
    username  TEXT NOT NULL,
    team_name TEXT NOT NULL REFERENCES teams(team_name) ON DELETE RESTRICT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX users_is_active_idx ON users(is_active);

CREATE TABLE pull_requests (
    pull_request_id TEXT PRIMARY KEY,
    pull_request_name VARCHAR NOT NULL,
    author_id TEXT NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    merged_at TIMESTAMP NULL
);

CREATE INDEX pr_author_idx ON pull_requests(author_id);
CREATE INDEX pr_status_idx ON pull_requests(status);
CREATE INDEX recent_open_prs_idx ON pull_requests(status, created_at);

CREATE TABLE pr_reviewers (
    pull_request_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    assigned_at TIMESTAMP NOT NULL DEFAULT NOW(),
    replaced_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX unique_pr_reviewer_idx ON pr_reviewers(pull_request_id, user_id);
CREATE INDEX reviewer_assignments_idx ON pr_reviewers(user_id);
CREATE INDEX active_reviewer_workload_idx ON pr_reviewers(user_id, replaced_at);

ALTER TABLE pull_requests 
ADD CONSTRAINT pull_requests_author_id_fkey 
FOREIGN KEY (author_id) REFERENCES users(user_id) ON DELETE CASCADE;

ALTER TABLE pr_reviewers 
ADD CONSTRAINT pr_reviewers_pull_request_id_fkey 
FOREIGN KEY (pull_request_id) REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE;

ALTER TABLE pr_reviewers 
ADD CONSTRAINT pr_reviewers_reviewer_id_fkey 
FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_pull_requests_updated_at ON pull_requests;
DROP FUNCTION IF EXISTS update_updated_at_column;

ALTER TABLE pr_reviewers DROP CONSTRAINT IF EXISTS pr_reviewers_reviewer_id_fkey;
ALTER TABLE pr_reviewers DROP CONSTRAINT IF EXISTS pr_reviewers_pull_request_id_fkey;
ALTER TABLE pull_requests DROP CONSTRAINT IF EXISTS pull_requests_author_id_fkey;

DROP INDEX IF EXISTS active_reviewer_workload_idx;
DROP INDEX IF EXISTS reviewer_assignments_idx;
DROP INDEX IF EXISTS unique_pr_reviewer_idx;
DROP TABLE IF EXISTS pr_reviewers;

DROP INDEX IF EXISTS recent_open_prs_idx;
DROP INDEX IF EXISTS pr_status_idx;
DROP INDEX IF EXISTS pr_author_idx;
DROP TABLE IF EXISTS pull_requests;

DROP INDEX IF EXISTS users_is_active_idx;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd

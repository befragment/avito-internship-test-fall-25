package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"avito-intern-test/internal/core"
	prmodel "avito-intern-test/internal/model/pullrequest"
)

func OpenTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	cfg, _ := core.LoadConfig()
	connStr := cfg.DBConnString()

	pcfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		t.Skipf("skip repository tests: parse config failed: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
	if err != nil {
		t.Skipf("skip repository tests: connect failed: %v", err)
	}
	pingCtx, pingCancel := context.WithTimeout(ctx, 2*time.Second)
	defer pingCancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		t.Skipf("skip repository tests: ping failed: %v", err)
	}
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = pool.Exec(ctx, "TRUNCATE TABLE pr_reviewers RESTART IDENTITY CASCADE")
		_, _ = pool.Exec(ctx, "TRUNCATE TABLE pull_requests RESTART IDENTITY CASCADE")
		_, _ = pool.Exec(ctx, "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
		_, _ = pool.Exec(ctx, "TRUNCATE TABLE teams RESTART IDENTITY CASCADE")
		pool.Close()
	})
	return pool
}


func TruncateAll(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	stmts := []string{
		"TRUNCATE TABLE pr_reviewers RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE pull_requests RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE users RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE teams RESTART IDENTITY CASCADE",
	}
	for _, s := range stmts {
		if _, err := pool.Exec(ctx, s); err != nil {
			t.Fatalf("truncate %q: %v", s, err)
		}
	}
}

func EnsureTeam(t *testing.T, pool *pgxpool.Pool, name string) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `INSERT INTO teams (team_name, created_at) VALUES ($1, NOW()) ON CONFLICT (team_name) DO NOTHING`, name)
	if err != nil {
		t.Fatalf("ensure team %s: %v", name, err)
	}
}

func EnsureUser(t *testing.T, pool *pgxpool.Pool, userID, username, team string, active bool) {
	t.Helper()
	EnsureTeam(t, pool, team)
	_, err := pool.Exec(context.Background(),
		`INSERT INTO users (user_id, username, team_name, is_active) VALUES ($1,$2,$3,$4)
		 ON CONFLICT (user_id) DO UPDATE SET username=EXCLUDED.username, team_name=EXCLUDED.team_name, is_active=EXCLUDED.is_active`,
		userID, username, team, active,
	)
	if err != nil {
		t.Fatalf("ensure user %s: %v", userID, err)
	}
}

func InsertPR(t *testing.T, pool *pgxpool.Pool, pr prmodel.PullRequest) {
	t.Helper()
	ctx := context.Background()
	_, err := pool.Exec(ctx,
		`INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at, merged_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		pr.PullRequestID, pr.PullRequestName, pr.AuthorID, string(pr.Status), pr.CreatedAt, pr.MergedAt,
	)
	if err != nil {
		t.Fatalf("insert pr %s: %v", pr.PullRequestID, err)
	}
	for _, r := range pr.AssignedReviewers {
		_, err := pool.Exec(ctx, `INSERT INTO pr_reviewers (pull_request_id, user_id) VALUES ($1,$2)`, pr.PullRequestID, r)
		if err != nil {
			t.Fatalf("insert reviewer %s: %v", r, err)
		}
	}
}

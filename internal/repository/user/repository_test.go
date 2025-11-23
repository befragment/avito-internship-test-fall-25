package repository

import (
	"context"
	"testing"

	usermodel "avito-intern-test/internal/model/user"
	"avito-intern-test/internal/repository/testutil"
)

func TestUserRepository_CRUD(t *testing.T) {
	pool := testutil.OpenTestPool(t)
	testutil.TruncateAll(t, pool)
	r := NewUserRepository(pool)
	ctx := context.Background()

	
	if _, err := pool.Exec(ctx, `INSERT INTO teams(team_name, created_at) VALUES ('t1', NOW())`); err != nil {
		t.Fatalf("seed team: %v", err)
	}

	
	u := usermodel.User{UserID: "u1", Username: "Alice", TeamName: "t1", IsActive: true}
	if err := r.CreateOrUpdate(ctx, u); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := r.GetByID(ctx, "u1")
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.Username != "Alice" || got.TeamName != "t1" || !got.IsActive {
		t.Fatalf("unexpected user: %+v", got)
	}
	u.Username = "Alice2"
	if err := r.CreateOrUpdate(ctx, u); err != nil {
		t.Fatalf("update: %v", err)
	}
	got, _ = r.GetByID(ctx, "u1")
	if got.Username != "Alice2" {
		t.Fatalf("upsert didn't update: %+v", got)
	}
	got, err = r.SetIsActive(ctx, "u1", false)
	if err != nil {
		t.Fatalf("set is_active: %v", err)
	}
	if got.IsActive {
		t.Fatalf("expected inactive")
	}
}

func TestUserRepository_ByTeam_And_ReviewerPRs(t *testing.T) {
	pool := testutil.OpenTestPool(t)
	testutil.TruncateAll(t, pool)
	r := NewUserRepository(pool)
	ctx := context.Background()

	if _, err := pool.Exec(ctx, `INSERT INTO teams(team_name, created_at) VALUES ('t1', NOW())`); err != nil {
		t.Fatalf("seed team: %v", err)
	}
	_ = r.CreateOrUpdate(ctx, usermodel.User{UserID: "u1", Username: "a", TeamName: "t1", IsActive: true})
	_ = r.CreateOrUpdate(ctx, usermodel.User{UserID: "u2", Username: "b", TeamName: "t1", IsActive: false})
	_ = r.CreateOrUpdate(ctx, usermodel.User{UserID: "u3", Username: "c", TeamName: "x", IsActive: true})

	users, err := r.GetByTeam(ctx, "t1")
	if err != nil {
		t.Fatalf("get by team: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO pull_requests(pull_request_id, pull_request_name, author_id, status, created_at) 
		VALUES ('pr1','x','u1','OPEN', NOW()), ('pr2','y','u1','OPEN', NOW());
		INSERT INTO pr_reviewers(pull_request_id, user_id) VALUES ('pr1','u2'), ('pr2','u2');
	`); err != nil {
		t.Fatalf("seed PRs: %v", err)
	}
	ids, err := r.GetReviewerPRs(ctx, "u2")
	if err != nil {
		t.Fatalf("get reviewer prs: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 ids, got %d", len(ids))
	}
}

func TestUserRepository_GetReviewerPRs(t *testing.T) {
	pool := testutil.OpenTestPool(t)
	testutil.TruncateAll(t, pool)
	r := NewUserRepository(pool)

	ctx := context.Background()
	if _, err := pool.Exec(ctx, `INSERT INTO teams(team_name, created_at) VALUES ('t1', NOW())`); err != nil {
		t.Fatalf("seed team: %v", err)
	}

	u := usermodel.User{UserID: "u1", Username: "Alice", TeamName: "t1", IsActive: true}
	if err := r.CreateOrUpdate(ctx, u); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := r.GetByID(ctx, "u1")
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.TeamName != "t1" || got.Username != "Alice" || !got.IsActive {
		t.Fatalf("unexpected user: %+v", got)
	}

	u.Username = "Alice2"
	if err := r.CreateOrUpdate(ctx, u); err != nil {
		t.Fatalf("update: %v", err)
	}
	got, _ = r.GetByID(ctx, "u1")
	if got.Username != "Alice2" {
		t.Fatalf("upsert failed: %+v", got)
	}

	got, err = r.SetIsActive(ctx, "u1", false)
	if err != nil {
		t.Fatalf("set is_active: %v", err)
	}
	if got.IsActive {
		t.Fatalf("expected inactive")
	}
}

func TestUserRepository_GetByTeam(t *testing.T) {
	pool := testutil.OpenTestPool(t)
	testutil.TruncateAll(t, pool)
	r := NewUserRepository(pool)
	ctx := context.Background()

	if _, err := pool.Exec(ctx, `INSERT INTO teams(team_name, created_at) VALUES ('t1', NOW())`); err != nil {
		t.Fatalf("seed team: %v", err)
	}
	_ = r.CreateOrUpdate(ctx, usermodel.User{UserID: "u1", Username: "a", TeamName: "t1", IsActive: true})
	_ = r.CreateOrUpdate(ctx, usermodel.User{UserID: "u2", Username: "b", TeamName: "t1", IsActive: true})
	_ = r.CreateOrUpdate(ctx, usermodel.User{UserID: "u3", Username: "c", TeamName: "t2", IsActive: true})

	users, err := r.GetByTeam(ctx, "t1")
	if err != nil {
		t.Fatalf("get by team: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	if _, err := pool.Exec(ctx, `
		INSERT INTO pull_requests(pull_request_id, pull_request_name, author_id, status, created_at) 
		VALUES ('pr1','x','u1','OPEN', NOW()), ('pr2','y','u1','OPEN', NOW());
		INSERT INTO pr_reviewers(pull_request_id, user_id) VALUES ('pr1','u2'), ('pr2','u2');
	`); err != nil {
		t.Fatalf("seed PRs: %v", err)
	}
	ids, err := r.GetReviewerPRs(ctx, "u2")
	if err != nil {
		t.Fatalf("get reviewer prs: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 PR ids, got %d", len(ids))
	}
}

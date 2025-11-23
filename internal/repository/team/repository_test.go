package repository

import (
	"context"
	"testing"

	"avito-intern-test/internal/repository/testutil"
)

func TestTeamRepository_Create_Exists_GetMembers(t *testing.T) {
	pool := testutil.OpenTestPool(t)
	testutil.TruncateAll(t, pool)
	r := NewTeamRepository(pool)
	ctx := context.Background()

	team, err := r.Create(ctx, "backend")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if team == nil || team.Name != "backend" {
		t.Fatalf("unexpected team: %+v", team)
	}
	ok, err := r.Exists(ctx, "backend")
	if err != nil || !ok {
		t.Fatalf("exists: %v ok=%v", err, ok)
	}
	testutil.EnsureTeam(t, pool, "other")
	testutil.EnsureUser(t, pool, "u1", "a", "backend", true)
	testutil.EnsureUser(t, pool, "u2", "b", "backend", false)
	testutil.EnsureUser(t, pool, "u3", "c", "other", true)
	users, err := r.GetTeamMembers(ctx, "backend")
	if err != nil {
		t.Fatalf("get members: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 members, got %d", len(users))
	}
}

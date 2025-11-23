package repository

import (
	"context"
	"testing"
	"time"

	prmodel "avito-intern-test/internal/model/pullrequest"
	"avito-intern-test/internal/repository/testutil"
)

func TestPullRequestRepository_Lifecycle(t *testing.T) {
	pool := testutil.OpenTestPool(t)
	testutil.TruncateAll(t, pool)
	r := NewPullRequestRepository(pool)
	ctx := context.Background()

	testutil.EnsureTeam(t, pool, "t1")
	testutil.EnsureUser(t, pool, "a1", "auth", "t1", true)
	testutil.EnsureUser(t, pool, "r1", "rev1", "t1", true)
	testutil.EnsureUser(t, pool, "r2", "rev2", "t1", true)

	ok, err := r.Exists(ctx, "pr-1")
	if err != nil || ok {
		t.Fatalf("exists expected false, err=%v ok=%v", err, ok)
	}

	now := time.Now().UTC()
	pr := prmodel.PullRequest{
		PullRequestID:     "pr-1",
		PullRequestName:   "Test",
		AuthorID:          "a1",
		Status:            prmodel.PullRequestStatusOpen,
		AssignedReviewers: []string{"r1", "r2"},
		CreatedAt:         now,
	}
	if err := r.Create(ctx, pr); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := r.GetByID(ctx, "pr-1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.PullRequestName != "Test" || len(got.AssignedReviewers) != 2 {
		t.Fatalf("unexpected pr: %+v", got)
	}

	got.AssignedReviewers = []string{"r2"}
	got.Status = prmodel.PullRequestStatusMerged
	tm := time.Now().UTC()
	got.MergedAt = &tm
	if err := r.Update(ctx, got); err != nil {
		t.Fatalf("update: %v", err)
	}
	got2, err := r.GetByID(ctx, "pr-1")
	if err != nil || len(got2.AssignedReviewers) != 1 || got2.Status != prmodel.PullRequestStatusMerged {
		t.Fatalf("unexpected after update: %+v err=%v", got2, err)
	}

	list, err := r.ReviewerPRs(ctx, "r2")
	if err != nil {
		t.Fatalf("reviewer prs: %v", err)
	}
	if len(list) == 0 || list[0].PullRequestID != "pr-1" {
		t.Fatalf("expected pr-1 in list, got %+v", list)
	}
}

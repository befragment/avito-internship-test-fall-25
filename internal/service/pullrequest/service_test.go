package service

import (
	"context"
	"math/rand"
	"strings"
	"testing"

	"avito-intern-test/internal/core"
	prmodel "avito-intern-test/internal/model/pullrequest"
	usermodel "avito-intern-test/internal/model/user"
)

type prRepoMock struct {
	exists    bool
	existsErr error
	storage   map[string]prmodel.PullRequest
}

func (m *prRepoMock) Exists(ctx context.Context, prID string) (bool, error) {
	return m.exists, m.existsErr
}
func (m *prRepoMock) Create(ctx context.Context, pr prmodel.PullRequest) error {
	if m.storage == nil {
		m.storage = map[string]prmodel.PullRequest{}
	}
	m.storage[pr.PullRequestID] = pr
	return nil
}
func (m *prRepoMock) GetByID(ctx context.Context, prID string) (prmodel.PullRequest, error) {
	if m.storage == nil {
		return prmodel.PullRequest{}, context.Canceled
	}
	pr, ok := m.storage[prID]
	if !ok {
		return prmodel.PullRequest{}, context.Canceled
	}
	return pr, nil
}
func (m *prRepoMock) Update(ctx context.Context, pr prmodel.PullRequest) error {
	if m.storage == nil {
		m.storage = map[string]prmodel.PullRequest{}
	}
	m.storage[pr.PullRequestID] = pr
	return nil
}

type teamRepoMockForPR struct {
	exists bool
}

func (t *teamRepoMockForPR) Exists(ctx context.Context, teamName string) (bool, error) {
	return t.exists, nil
}

type userRepoMockForPR struct {
	users  map[string]usermodel.User
	byTeam map[string][]usermodel.User
}

func (u *userRepoMockForPR) GetByID(ctx context.Context, userID string) (usermodel.User, error) {
	if u.users == nil {
		return usermodel.User{}, context.Canceled
	}
	val, ok := u.users[userID]
	if !ok {
		return usermodel.User{}, context.Canceled
	}
	return val, nil
}
func (u *userRepoMockForPR) GetByTeam(ctx context.Context, teamName string) ([]usermodel.User, error) {
	if u.byTeam == nil {
		return nil, nil
	}
	return u.byTeam[teamName], nil
}

func TestPRService_CreatePR_Success(t *testing.T) {
	prr := &prRepoMock{}
	tr := &teamRepoMockForPR{exists: true}
	ur := &userRepoMockForPR{
		users: map[string]usermodel.User{
			"a1": {UserID: "a1", TeamName: "backend", IsActive: true},
		},
		byTeam: map[string][]usermodel.User{
			"backend": {
				{UserID: "a1", TeamName: "backend", IsActive: true},
				{UserID: "r1", TeamName: "backend", IsActive: true},
				{UserID: "r2", TeamName: "backend", IsActive: true},
			},
		},
	}
	svc := NewPRService(ur, tr, prr)
	svc.rand = rand.New(rand.NewSource(1))

	pr, err := svc.CreatePR(context.Background(), "pr-1", "Test", "a1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if pr.PullRequestID != "pr-1" || pr.Status != prmodel.PullRequestStatusOpen {
		t.Fatalf("unexpected pr: %+v", pr)
	}
	if len(pr.AssignedReviewers) == 0 {
		t.Fatalf("expected reviewers assigned")
	}
}

func TestPRService_CreatePR_AlreadyExists(t *testing.T) {
	prr := &prRepoMock{exists: true}
	tr := &teamRepoMockForPR{exists: true}
	ur := &userRepoMockForPR{}
	svc := NewPRService(ur, tr, prr)
	_, err := svc.CreatePR(context.Background(), "pr-1", "Test", "a1")
	if err == nil || !strings.Contains(err.Error(), core.ErrorPRExists) {
		t.Fatalf("expected PR_EXISTS, got %v", err)
	}
}

func TestPRService_MergePR_Idempotent(t *testing.T) {
	prr := &prRepoMock{storage: map[string]prmodel.PullRequest{
		"pr-1": {PullRequestID: "pr-1", Status: prmodel.PullRequestStatusOpen},
	}}
	tr := &teamRepoMockForPR{}
	ur := &userRepoMockForPR{}
	svc := NewPRService(ur, tr, prr)
	pr, err := svc.MergePR(context.Background(), "pr-1")
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if pr.Status != prmodel.PullRequestStatusMerged {
		t.Fatalf("expected merged")
	}
	pr, err = svc.MergePR(context.Background(), "pr-1")
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if pr.Status != prmodel.PullRequestStatusMerged {
		t.Fatalf("expected merged again")
	}
}

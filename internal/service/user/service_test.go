package service

import (
	"context"
	"testing"

	prmodel "avito-intern-test/internal/model/pullrequest"
	usermodel "avito-intern-test/internal/model/user"
)

type userRepoMockForUserService struct {
	setResp usermodel.User
	setErr  error
	prIDs   []string
	prErr   error
}

func (m *userRepoMockForUserService) GetByID(ctx context.Context, userID string) (usermodel.User, error) {
	return usermodel.User{UserID: userID, Username: "x", TeamName: "t", IsActive: true}, nil
}

func (m *userRepoMockForUserService) GetReviewerPRs(ctx context.Context, ReviewerID string) ([]string, error) {
	return m.prIDs, m.prErr
}

func (m *userRepoMockForUserService) SetIsActive(ctx context.Context, userID string, flag bool) (usermodel.User, error) {
	m.setResp.UserID = userID
	m.setResp.IsActive = flag
	return m.setResp, m.setErr
}

type prRepoMockForUserService struct {
	prs  []prmodel.PullRequest
	err  error
	reqs []string
}

func (m *prRepoMockForUserService) GetMany(ctx context.Context, PRIDs []string) ([]prmodel.PullRequest, error) {
	m.reqs = append(m.reqs, PRIDs...)
	return m.prs, m.err
}

func TestUserService_SetIsActive(t *testing.T) {
	ur := &userRepoMockForUserService{}
	prr := &prRepoMockForUserService{}
	svc := NewUserService(ur, prr)

	user, err := svc.SetIsActive(context.Background(), "u1", false)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if user.UserID != "u1" || user.IsActive != false {
		t.Fatalf("unexpected user result: %+v", user)
	}
}

func TestUserService_GetReviewerPRs(t *testing.T) {
	ur := &userRepoMockForUserService{prIDs: []string{"p1", "p2"}}
	prr := &prRepoMockForUserService{
		prs: []prmodel.PullRequest{
			{PullRequestID: "p1"},
			{PullRequestID: "p2"},
		},
	}
	svc := NewUserService(ur, prr)
	prs, err := svc.GetReviewerPRs(context.Background(), "u5")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(prs) != 2 {
		t.Fatalf("expected 2 PRs, got %d", len(prs))
	}
}

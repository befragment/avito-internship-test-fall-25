package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"avito-intern-test/internal/core"
	teammodel "avito-intern-test/internal/model/team"
	usermodel "avito-intern-test/internal/model/user"
)

type teamRepoMock struct {
	existsResp bool
	existsErr  error
	created    *teammodel.Team
	createErr  error
	members    []usermodel.User
	membersErr error
}

func (m *teamRepoMock) GetTeamMembers(ctx context.Context, teamName string) ([]usermodel.User, error) {
	return m.members, m.membersErr
}
func (m *teamRepoMock) Exists(ctx context.Context, teamName string) (bool, error) {
	return m.existsResp, m.existsErr
}
func (m *teamRepoMock) Create(ctx context.Context, teamName string) (*teammodel.Team, error) {
	if m.created == nil {
		now := time.Now()
		m.created = &teammodel.Team{Name: teamName, CreatedAt: now}
	}
	return m.created, m.createErr
}

type userRepoMock struct {
	usersByID        map[string]usermodel.User
	createOrUpdateFn func(user usermodel.User) error
	getErr           error
}

func (m *userRepoMock) CreateOrUpdate(_ context.Context, user usermodel.User) error {
	if m.createOrUpdateFn != nil {
		return m.createOrUpdateFn(user)
	}
	if m.usersByID == nil {
		m.usersByID = map[string]usermodel.User{}
	}
	m.usersByID[user.UserID] = user
	return nil
}
func (m *userRepoMock) GetByID(_ context.Context, userID string) (usermodel.User, error) {
	if m.getErr != nil {
		return usermodel.User{}, m.getErr
	}
	if m.usersByID == nil {
		return usermodel.User{}, context.Canceled 
	}
	u, ok := m.usersByID[userID]
	if !ok {
		return usermodel.User{}, context.Canceled
	}
	return u, nil
}

func TestTeamService_CreateWithMembers_SuccessCreateNewTeam(t *testing.T) {
	tr := &teamRepoMock{existsResp: false}
	ur := &userRepoMock{usersByID: map[string]usermodel.User{}}
	svc := NewTeamService(tr, ur)

	members := []usermodel.User{
		{UserID: "u1", Username: "Alice", IsActive: true},
		{UserID: "u2", Username: "Bob", IsActive: true},
	}
	team, err := svc.CreateWithMembers(context.Background(), "backend", members)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team == nil || team.Name != "backend" {
		t.Fatalf("team not created properly")
	}
	if _, ok := ur.usersByID["u1"]; !ok {
		t.Fatalf("user u1 not created")
	}
	if ur.usersByID["u1"].TeamName != "backend" {
		t.Fatalf("user u1 team not set")
	}
}

func TestTeamService_CreateWithMembers_FailsIfUserInAnotherTeam(t *testing.T) {
	tr := &teamRepoMock{existsResp: true, created: &teammodel.Team{Name: "backend", CreatedAt: time.Now()}}
	ur := &userRepoMock{
		usersByID: map[string]usermodel.User{
			"u1": {UserID: "u1", Username: "Alice", TeamName: "payments", IsActive: true},
		},
	}
	svc := NewTeamService(tr, ur)
	_, err := svc.CreateWithMembers(context.Background(), "backend", []usermodel.User{
		{UserID: "u1", Username: "Alice", IsActive: true},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), core.ErrorUserExists) {
		t.Fatalf("expected code %s, got %v", core.ErrorUserExists, err)
	}
}

func TestTeamService_GetTeamMembers_NotFound(t *testing.T) {
	tr := &teamRepoMock{existsResp: false}
	ur := &userRepoMock{}
	svc := NewTeamService(tr, ur)
	_, err := svc.GetTeamMembers(context.Background(), "unknown")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != ErrTeamNotFound {
		t.Fatalf("expected ErrTeamNotFound, got %v", err)
	}
}

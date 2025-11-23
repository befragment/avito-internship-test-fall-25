package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	teammodel "avito-intern-test/internal/model/team"
	usermodel "avito-intern-test/internal/model/user"
)

type teamServiceMock struct {
	createResp *teammodel.Team
	createErr  error
	members    []usermodel.User
	getErr     error
}

func (m *teamServiceMock) GetTeamMembers(_ context.Context, name string) ([]usermodel.User, error) {
	return m.members, m.getErr
}
func (m *teamServiceMock) CreateWithMembers(_ context.Context, name string, members []usermodel.User) (*teammodel.Team, error) {
	return m.createResp, m.createErr
}

func TestTeamHandler_CreateTeam_Created(t *testing.T) {
	h := NewTeamHandler(&teamServiceMock{
		createResp: &teammodel.Team{Name: "backend"},
	})
	body := CreateWithMembersRequest{
		TeamName: "backend",
		Members: []usermodel.User{
			{UserID: "u1", Username: "a", IsActive: true},
			{UserID: "u2", Username: "b", IsActive: true},
		},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.CreateTeam(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestTeamHandler_GetTeam_MissingQuery(t *testing.T) {
	h := NewTeamHandler(&teamServiceMock{})
	req := httptest.NewRequest(http.MethodGet, "/team/get", nil)
	w := httptest.NewRecorder()
	h.GetTeam(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	prmodel "avito-intern-test/internal/model/pullrequest"
	usermodel "avito-intern-test/internal/model/user"
)

type userServiceMock struct {
	setErr error
	prs    []prmodel.PullRequest
	prErr  error
}

func (m *userServiceMock) SetIsActive(_ context.Context, userID string, flag bool) (usermodel.User, error) {
	return usermodel.User{
		UserID:   userID,
		Username: "x",
		TeamName: "t",
		IsActive: flag,
	}, m.setErr
}

func (m *userServiceMock) GetReviewerPRs(_ context.Context, reviewerID string) ([]prmodel.PullRequest, error) {
	return m.prs, m.prErr
}

func TestUserHandler_SetIsActive_OK(t *testing.T) {
	h := NewUserHandler(&userServiceMock{})
	body := SetIsActiveRequest{UserID: "u1", IsActive: false}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.SetIsActive(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body=%s", w.Code, w.Body.String())
	}
}

func TestUserHandler_GetReview(t *testing.T) {
	h := NewUserHandler(&userServiceMock{})
	req := httptest.NewRequest(http.MethodGet, "/users/getReview", nil)
	w := httptest.NewRecorder()
	h.GetReview(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

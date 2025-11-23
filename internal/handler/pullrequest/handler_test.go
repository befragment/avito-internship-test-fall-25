package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	prmodel "avito-intern-test/internal/model/pullrequest"
)

type prServiceMock struct {
	createResp *prmodel.PullRequest
	createErr  error
	mergeResp  *prmodel.PullRequest
	mergeErr   error
	reResp     *prmodel.PullRequest
	reUser     string
	reErr      error
}

func (m *prServiceMock) CreatePR(_ context.Context, id, name, authorID string) (*prmodel.PullRequest, error) {
	return m.createResp, m.createErr
}
func (m *prServiceMock) MergePR(_ context.Context, id string) (*prmodel.PullRequest, error) {
	return m.mergeResp, m.mergeErr
}
func (m *prServiceMock) ReassignReviewer(_ context.Context, prID, oldUserID string) (*prmodel.PullRequest, string, error) {
	return m.reResp, m.reUser, m.reErr
}

func TestPRHandler_Create_BadJSON(t *testing.T) {
	h := NewPullRequestHandler(&prServiceMock{})
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBufferString("{"))
	w := httptest.NewRecorder()
	h.CreatePullRequest(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPRHandler_Create_Success(t *testing.T) {
	h := NewPullRequestHandler(&prServiceMock{
		createResp: &prmodel.PullRequest{
			PullRequestID:   "pr-1",
			PullRequestName: "Add",
			AuthorID:        "u1",
			Status:          prmodel.PullRequestStatusOpen,
		},
	})
	body := CreatePRRequest{PullRequestID: "pr-1", PullRequestName: "Add", AuthorID: "u1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.CreatePullRequest(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d; body=%s", w.Code, w.Body.String())
	}
}

func TestPRHandler_Reassign_AcceptsOldReviewerID(t *testing.T) {
	h := NewPullRequestHandler(&prServiceMock{
		reResp: &prmodel.PullRequest{
			PullRequestID: "pr-1",
			Status:        prmodel.PullRequestStatusOpen,
		},
		reUser: "u5",
	})
	body := map[string]string{
		"pull_request_id": "pr-1",
		"old_reviewer_id": "u2",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.ReassignPullRequest(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body=%s", w.Code, w.Body.String())
	}
}

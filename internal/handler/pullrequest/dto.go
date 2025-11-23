package handler

import (
	"time"

	prmodel "avito-intern-test/internal/model/pullrequest"
)

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignPRRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id,omitempty"`
}

type PullRequestDTO struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

type CreatePRResponse struct {
	PR PullRequestDTO `json:"pr"`
}

type MergePRResponse struct {
	PR PullRequestDTO `json:"pr"`
}

type ReassignPRResponse struct {
	PR         PullRequestDTO `json:"pr"`
	ReplacedBy string         `json:"replaced_by"`
}

func prModelToDTO(m prmodel.PullRequest) PullRequestDTO {
	dto := PullRequestDTO{
		PullRequestID:     m.PullRequestID,
		PullRequestName:   m.PullRequestName,
		AuthorID:          m.AuthorID,
		Status:            string(m.Status),
		AssignedReviewers: append([]string(nil), m.AssignedReviewers...),
	}
	if !m.CreatedAt.IsZero() {
		t := m.CreatedAt.UTC()
		dto.CreatedAt = &t
	}
	if m.MergedAt != nil && !m.MergedAt.IsZero() {
		t := m.MergedAt.UTC()
		dto.MergedAt = &t
	}
	return dto
}

package handler

import (
	prmodel "avito-intern-test/internal/model/pullrequest"
	"context"
)

type pullReqeustService interface {
	CreatePR(ctx context.Context, id, name, authorID string) (*prmodel.PullRequest, error)
	MergePR(ctx context.Context, id string) (*prmodel.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldUserID string) (*prmodel.PullRequest, string, error)
}

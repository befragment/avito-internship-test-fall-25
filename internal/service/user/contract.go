package service

import (
	"context"

	prmodel "avito-intern-test/internal/model/pullrequest"
	usermodel "avito-intern-test/internal/model/user"
)

type userRepository interface {
	GetByID(ctx context.Context, userID string) (usermodel.User, error)
	GetReviewerPRs(ctx context.Context, ReviewerID string) ([]string, error)
	SetIsActive(ctx context.Context, userID string, flag bool) (usermodel.User, error)
}

type pullRequestRepository interface {
	GetMany(ctx context.Context, PRIDs []string) ([]prmodel.PullRequest, error)
}

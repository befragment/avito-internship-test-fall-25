package handler

import (
	"context"

	prmodel "avito-intern-test/internal/model/pullrequest"
	usermodel "avito-intern-test/internal/model/user"
)

type userService interface {
	SetIsActive(ctx context.Context, userID string, flag bool) (usermodel.User, error)
	GetReviewerPRs(ctx context.Context, ReviewerID string) ([]prmodel.PullRequest, error)
}

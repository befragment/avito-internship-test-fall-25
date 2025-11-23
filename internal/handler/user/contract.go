package handler

import (
	prmodel "avito-intern-test/internal/model/pullrequest"
	usermodel "avito-intern-test/internal/model/user"
	"context"
)

type userService interface {
	SetIsActive(ctx context.Context, userID string, flag bool) (usermodel.User, error)
	GetReviewerPRs(ctx context.Context, ReviewerID string) ([]prmodel.PullRequest, error)
}

package service

import (
	prmodel "avito-intern-test/internal/model/pullrequest"
	usermodel "avito-intern-test/internal/model/user"
	"context"
)

type (
	pullrequestRepository interface {
		Exists(ctx context.Context, prID string) (bool, error)
		Create(ctx context.Context, pr prmodel.PullRequest) error
		GetByID(ctx context.Context, prID string) (prmodel.PullRequest, error)
		Update(ctx context.Context, pr prmodel.PullRequest) error
	}

	teamRepository interface {
		Exists(ctx context.Context, teamName string) (bool, error)
	}

	userRepository interface {
		GetByID(ctx context.Context, userID string) (usermodel.User, error)
		GetByTeam(ctx context.Context, teamName string) ([]usermodel.User, error)
	}
)

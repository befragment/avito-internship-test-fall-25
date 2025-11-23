package service

import (
	teammodel "avito-intern-test/internal/model/team"
	usermodel "avito-intern-test/internal/model/user"
	"context"
)

type teamRepository interface {
	GetTeamMembers(ctx context.Context, teamName string) ([]usermodel.User, error)
	Exists(ctx context.Context, teamName string) (bool, error)
	Create(ctx context.Context, teamName string) (*teammodel.Team, error)
}

type userRepository interface {
	CreateOrUpdate(ctx context.Context, user usermodel.User) error
	GetByID(ctx context.Context, userID string) (usermodel.User, error)
}

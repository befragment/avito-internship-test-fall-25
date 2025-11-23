package handler

import (
	"context"

	teammodel "avito-intern-test/internal/model/team"
	usermodel "avito-intern-test/internal/model/user"
)

type teamService interface {
	GetTeamMembers(ctx context.Context, name string) ([]usermodel.User, error)
	CreateWithMembers(ctx context.Context, name string, members []usermodel.User) (*teammodel.Team, error)
}

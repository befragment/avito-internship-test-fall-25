package handler

import (
	usermodel "avito-intern-test/internal/model/user"
)

type CreateWithMembersRequest struct {
	TeamName string           `json:"team_name"`
	Members  []usermodel.User `json:"members"`
}

type TeamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamDTO struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type CreateTeamResponse struct {
	Team TeamDTO `json:"team"`
}

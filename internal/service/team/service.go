package service

import (
	"context"
	"fmt"
	"time"

	"avito-intern-test/internal/core"
	teammodel "avito-intern-test/internal/model/team"
	usermodel "avito-intern-test/internal/model/user"
)

type TeamService struct {
	teamRepository teamRepository
	userRepository userRepository
}

func NewTeamService(
	teamRepository teamRepository,
	userRepository userRepository,
) *TeamService {
	return &TeamService{
		teamRepository: teamRepository,
		userRepository: userRepository,
	}
}

func (s *TeamService) GetTeamMembers(
	ctx context.Context,
	teamName string,
) ([]usermodel.User, error) {
	exists, _ := s.teamRepository.Exists(ctx, teamName)
	if !exists {
		return nil, ErrTeamNotFound
	}
	members, err := s.teamRepository.GetTeamMembers(ctx, teamName)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (s *TeamService) CreateWithMembers(
	ctx context.Context,
	teamName string,
	members []usermodel.User,
) (*teammodel.Team, error) {
	for _, m := range members {
		existing, err := s.userRepository.GetByID(ctx, m.UserID)
		if err == nil {
			if existing.TeamName != "" && existing.TeamName != teamName {
				return nil, core.Throw(core.ErrorUserExists, "user already in another team")
			}
			if existing.TeamName == teamName {
				return nil, core.Throw(core.ErrorUserExists, "user already in team")
			}
		}
	}

	exists, _ := s.teamRepository.Exists(ctx, teamName)
	var createdTeam *teammodel.Team
	if !exists {
		t, err := s.teamRepository.Create(ctx, teamName)
		if err != nil {
			return nil, fmt.Errorf("create team: %w", err)
		}
		createdTeam = t
	} else {
		createdTeam = &teammodel.Team{
			Name:      teamName,
			CreatedAt: time.Now(),
		}
	}

	for _, m := range members {
		if existing, err := s.userRepository.GetByID(ctx, m.UserID); err == nil {
			existing.Username = m.Username
			existing.IsActive = m.IsActive
			existing.TeamName = teamName
			existing.CreatedAt = time.Now()
			if err := s.userRepository.CreateOrUpdate(ctx, existing); err != nil {
				return nil, fmt.Errorf("create or update user %s: %w", m.UserID, err)
			}
		} else {
			newUser := usermodel.User{
				UserID:    m.UserID,
				Username:  m.Username,
				TeamName:  teamName,
				IsActive:  m.IsActive,
				CreatedAt: time.Now(),
			}
			if err := s.userRepository.CreateOrUpdate(ctx, newUser); err != nil {
				return nil, fmt.Errorf("create or update user %s: %w", m.UserID, err)
			}
		}
	}
	return createdTeam, nil
}

package service

import (
	prmodel "avito-intern-test/internal/model/pullrequest"
	usermodel "avito-intern-test/internal/model/user"
	"context"
)

type UserService struct {
	userRepository        userRepository
	pullRequestRepository pullRequestRepository
}

func NewUserService(
	userRepository userRepository,
	pullRequestRepository pullRequestRepository,
) *UserService {
	return &UserService{
		userRepository:        userRepository,
		pullRequestRepository: pullRequestRepository,
	}
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, flag bool) (usermodel.User, error) {
	user, err := s.userRepository.SetIsActive(ctx, userID, flag)
	if err != nil {
		return usermodel.User{}, err
	}
	return user, nil
}

func (s *UserService) GetReviewerPRs(
	ctx context.Context,
	ReviewerID string,
) ([]prmodel.PullRequest, error) {
	reviewerPRIDs, err := s.userRepository.GetReviewerPRs(ctx, ReviewerID)
	if err != nil {
		return nil, err
	}
	PRs, err := s.pullRequestRepository.GetMany(ctx, reviewerPRIDs)
	if err != nil {
		return nil, err
	}
	return PRs, nil
}

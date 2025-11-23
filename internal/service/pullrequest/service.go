package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"avito-intern-test/internal/core"
	prmodel "avito-intern-test/internal/model/pullrequest"
	usermodel "avito-intern-test/internal/model/user"
)

type PRService struct {
	userRepository        userRepository
	teamRepository        teamRepository
	pullRequestRepository pullrequestRepository
	rand                  *rand.Rand
}

func NewPRService(
	userRepository userRepository,
	teamRepository teamRepository,
	pullRequestRepository pullrequestRepository,
) *PRService {
	return &PRService{
		userRepository:        userRepository,
		teamRepository:        teamRepository,
		pullRequestRepository: pullRequestRepository,
		rand:                  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *PRService) CreatePR(ctx context.Context, pullRequestID string, pullRequestName string, authorID string) (*prmodel.PullRequest, error) {
	exists, err := s.pullRequestRepository.Exists(ctx, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("check PR exists: %w", err)
	}
	if exists {
		return nil, core.Throw(core.ErrorPRExists, "PR id already exists")
	}

	author, err := s.userRepository.GetByID(ctx, authorID)
	if err != nil {
		return nil, core.Throw(core.ErrorNotFound, "author not found")
	}

	if author.TeamName == "" {
		return nil, core.Throw(core.ErrorNotFound, "author team not found")
	}

	teamExists, err := s.teamRepository.Exists(ctx, author.TeamName)
	if err != nil {
		return nil, fmt.Errorf("check team exists: %w", err)
	}
	if !teamExists {
		return nil, core.Throw(core.ErrorNotFound, "author team not found")
	}

	users, err := s.userRepository.GetByTeam(ctx, author.TeamName)
	if err != nil {
		return nil, fmt.Errorf("get team members: %w", err)
	}

	var candidates []usermodel.User
	for _, u := range users {
		if u.UserID == authorID {
			continue
		}
		if !u.IsActive {
			continue
		}
		candidates = append(candidates, u)
	}

	reviewers := chooseReviewers(candidates, 2, s.rand)

	now := time.Now().UTC()
	pr := prmodel.PullRequest{
		PullRequestID:     pullRequestID,
		PullRequestName:   pullRequestName,
		AuthorID:          authorID,
		Status:            prmodel.PullRequestStatusOpen,
		AssignedReviewers: reviewers,
		CreatedAt:         now,
		MergedAt:          nil,
	}

	if err := s.pullRequestRepository.Create(ctx, pr); err != nil {
		return nil, fmt.Errorf("create PR: %w", err)
	}

	return &pr, nil
}

func (s *PRService) MergePR(ctx context.Context, id string) (*prmodel.PullRequest, error) {
	pr, err := s.pullRequestRepository.GetByID(ctx, id)
	if err != nil {
		return nil, core.Throw(core.ErrorNotFound, "pr not found")
	}

	if pr.Status == prmodel.PullRequestStatusMerged {
		return &pr, nil
	}

	now := time.Now().UTC()
	pr.Status = prmodel.PullRequestStatusMerged
	pr.MergedAt = &now

	if err := s.pullRequestRepository.Update(ctx, pr); err != nil {
		return nil, fmt.Errorf("update PR: %w", err)
	}

	return &pr, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*prmodel.PullRequest, string, error) {
	pr, err := s.pullRequestRepository.GetByID(ctx, prID)
	if err != nil {
		return nil, "", fmt.Errorf("get PR: %w", err)
	}

	if pr.Status == prmodel.PullRequestStatusMerged {
		return nil, "", core.Throw(core.ErrorPRMerged, "cannot reassign on merged PR")
	}

	idx := -1
	for i, id := range pr.AssignedReviewers {
		if id == oldUserID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, "", core.Throw(core.ErrorNotAssigned, "reviewer is not assigned to this PR")
	}

	oldUser, err := s.userRepository.GetByID(ctx, oldUserID)
	if err != nil {
		return nil, "", fmt.Errorf("get old reviewer: %w", err)
	}

	if oldUser.TeamName == "" {
		return nil, "", core.Throw(core.ErrorNotFound, "reviewer team not found")
	}

	users, err := s.userRepository.GetByTeam(ctx, oldUser.TeamName)
	if err != nil {
		return nil, "", fmt.Errorf("get team members: %w", err)
	}

	current := make(map[string]struct{}, len(pr.AssignedReviewers))
	for _, id := range pr.AssignedReviewers {
		current[id] = struct{}{}
	}
	delete(current, oldUserID)
	authorID := pr.AuthorID

	var candidates []usermodel.User
	for _, u := range users {
		if !u.IsActive {
			continue
		}
		if u.UserID == oldUserID {
			continue
		}
		if u.UserID == authorID {
			continue
		}
		if _, exists := current[u.UserID]; exists {
			continue
		}
		candidates = append(candidates, u)
	}

	if len(candidates) == 0 {
		return nil, "", core.Throw(core.ErrorNoCandidate, "no active replacement candidate in team")
	}

	newUser := chooseReviewers(candidates, 1, s.rand)[0]

	pr.AssignedReviewers[idx] = newUser

	if err := s.pullRequestRepository.Update(ctx, pr); err != nil {
		return nil, "", fmt.Errorf("update PR after reassign: %w", err)
	}

	return &pr, newUser, nil
}

func chooseReviewers(users []usermodel.User, limit int, r *rand.Rand) []string {
	if len(users) == 0 || limit <= 0 {
		return nil
	}
	if len(users) <= limit {
		result := make([]string, 0, len(users))
		for _, u := range users {
			result = append(result, u.UserID)
		}
		return result
	}

	tmp := make([]usermodel.User, len(users))
	copy(tmp, users)

	r.Shuffle(len(tmp), func(i, j int) {
		tmp[i], tmp[j] = tmp[j], tmp[i]
	})

	result := make([]string, 0, limit)
	for i := 0; i < limit && i < len(tmp); i++ {
		result = append(result, tmp[i].UserID)
	}

	return result
}

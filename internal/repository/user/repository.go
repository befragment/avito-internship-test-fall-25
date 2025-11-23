package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	usermodel "avito-intern-test/internal/model/user"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) GetReviewerPRs(ctx context.Context, reviewerID string) ([]string, error) {
	queryBuilder := sq.
		Select("pull_request_id").
		From("pr_reviewers").
		Where(sq.Eq{"user_id": reviewerID}).
		OrderBy("pull_request_id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get reviewer PRs query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get reviewer PRs: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan pull_request_id: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reviewer PRs rows err: %w", err)
	}
	return ids, nil
}

func (r *UserRepository) CreateOrUpdate(
	ctx context.Context,
	user usermodel.User,
) error {
	queryBuilder := sq.
		Insert("users").
		Columns("user_id", "username", "team_name", "is_active").
		Values(user.UserID, user.Username, user.TeamName, user.IsActive).
		Suffix(`ON CONFLICT (user_id) DO UPDATE
				SET username = EXCLUDED.username,
					team_name = EXCLUDED.team_name,
					is_active = EXCLUDED.is_active`).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("build insert user query: %w", err)
	}

	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("create or update user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID string) (usermodel.User, error) {
	queryBuilder := sq.
		Select("user_id", "username", "team_name", "is_active").
		From("users").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return usermodel.User{}, fmt.Errorf("build get user by id query: %w", err)
	}

	var u usermodel.User
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&u.UserID,
		&u.Username,
		&u.TeamName,
		&u.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return usermodel.User{}, fmt.Errorf("user not found")
		}
		return usermodel.User{}, fmt.Errorf("get user by id: %w", err)
	}

	return u, nil
}

func (r *UserRepository) GetByTeam(ctx context.Context, teamName string) ([]usermodel.User, error) {
	queryBuilder := sq.
		Select("user_id", "username", "team_name", "is_active").
		From("users").
		Where(sq.Eq{"team_name": teamName}).
		OrderBy("user_id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get users by team query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get users by team: %w", err)
	}
	defer rows.Close()

	var users []usermodel.User
	for rows.Next() {
		var u usermodel.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("users rows err: %w", err)
	}

	return users, nil
}

func (r *UserRepository) SetIsActive(
	ctx context.Context,
	userID string,
	flag bool,
) (usermodel.User, error) {
	queryBuilder := sq.
		Update("users").
		Set("is_active", flag).
		Where(sq.Eq{"user_id": userID}).
		Suffix("RETURNING user_id, username, team_name, is_active").
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return usermodel.User{}, fmt.Errorf("build set is_active query: %w", err)
	}

	var u usermodel.User
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&u.UserID,
		&u.Username,
		&u.TeamName,
		&u.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return usermodel.User{}, fmt.Errorf("user not found")
		}
		return usermodel.User{}, fmt.Errorf("set user is_active: %w", err)
	}

	return u, nil
}

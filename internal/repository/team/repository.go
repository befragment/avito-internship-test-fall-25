package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	teammodel "avito-intern-test/internal/model/team"
	usermodel "avito-intern-test/internal/model/user"
)

type TeamRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool: pool}
}

func (r *TeamRepository) GetTeamMembers(
	ctx context.Context,
	teamName string,
) ([]usermodel.User, error) {
	queryBuilder := sq.
		Select("user_id", "username", "team_name", "is_active").
		From("users").
		Where(sq.Eq{"team_name": teamName}).
		OrderBy("user_id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []usermodel.User
	for rows.Next() {
		var user usermodel.User
		err = rows.Scan(
			&user.UserID, &user.Username, &user.TeamName, &user.IsActive,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *TeamRepository) Exists(ctx context.Context, teamName string) (bool, error) {
	queryBuilder := sq.
		Select("1").
		From("teams").
		Where(sq.Eq{"team_name": teamName}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return false, err
	}

	var exists = 0
	err = r.pool.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *TeamRepository) Create(
	ctx context.Context,
	teamName string,
) (*teammodel.Team, error) {
	queryBuilder := sq.
		Insert("teams").
		Columns("team_name", "created_at").
		Values(teamName, time.Now()).
		Suffix("RETURNING team_name, created_at").
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var createdAt time.Time
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&teamName, &createdAt); err != nil {
		return nil, err
	}
	return &teammodel.Team{
		Name:      teamName,
		CreatedAt: createdAt,
	}, nil
}

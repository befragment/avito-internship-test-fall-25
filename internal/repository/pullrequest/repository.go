package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	prmodel "avito-intern-test/internal/model/pullrequest"
)

type PullRequestRepository struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepository(pool *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{pool: pool}
}

func (r *PullRequestRepository) Exists(ctx context.Context, prID string) (bool, error) {
	queryBuilder := sq.
		Select("1").
		From("pull_requests").
		Where(sq.Eq{"pull_request_id": prID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return false, fmt.Errorf("build PR exists query: %w", err)
	}

	var flag int
	err = r.pool.QueryRow(ctx, query, args...).Scan(&flag)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("PR exists query: %w", err)
	}
	return flag == 1, nil
}

func (r *PullRequestRepository) Create(
	ctx context.Context,
	pr prmodel.PullRequest,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	createdAt := time.Now().UTC()

	queryBuilder := sq.
		Insert("pull_requests").
		Columns("pull_request_id", "pull_request_name", "author_id", "status", "created_at", "merged_at").
		Values(pr.PullRequestID, pr.PullRequestName, pr.AuthorID, string(pr.Status), createdAt, pr.MergedAt).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("build insert PR query: %w", err)
	}

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("insert pull_request: %w", err)
	}

	if len(pr.AssignedReviewers) > 0 {
		for _, id := range pr.AssignedReviewers {
			queryBuilder := sq.
				Insert("pr_reviewers").
				Columns("pull_request_id", "user_id").
				Values(pr.PullRequestID, id).
				PlaceholderFormat(sq.Dollar)

			query, args, err := queryBuilder.ToSql()
			if err != nil {
				return fmt.Errorf("build insert pr_reviewer query: %w", err)
			}

			if _, err := tx.Exec(ctx, query, args...); err != nil {
				return fmt.Errorf("insert pr_reviewer: %w", err)
			}
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *PullRequestRepository) GetByID(
	ctx context.Context,
	prID string,
) (prmodel.PullRequest, error) {
	queryBuilder := sq.
		Select("pull_request_id", "pull_request_name", "author_id", "status", "created_at", "merged_at").
		From("pull_requests").
		Where(sq.Eq{"pull_request_id": prID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return prmodel.PullRequest{}, fmt.Errorf("build get PR query: %w", err)
	}

	var pr prmodel.PullRequest
	var status string
	var createdAt time.Time

	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&status,
		&createdAt,
		&pr.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return prmodel.PullRequest{}, fmt.Errorf("PR not found: %w", err)
		}
		return prmodel.PullRequest{}, fmt.Errorf("get PR by id: %w", err)
	}
	pr.Status = prmodel.PullRequestStatus(status)
	pr.CreatedAt = createdAt

	queryBuilderReviewers := sq.
		Select("user_id").
		From("pr_reviewers").
		Where(sq.Eq{"pull_request_id": prID}).
		OrderBy("user_id").
		PlaceholderFormat(sq.Dollar)

	queryReviewers, argsReviewers, err := queryBuilderReviewers.ToSql()
	if err != nil {
		return prmodel.PullRequest{}, fmt.Errorf("build get reviewers query: %w", err)
	}

	rows, err := r.pool.Query(ctx, queryReviewers, argsReviewers...)
	if err != nil {
		return prmodel.PullRequest{}, fmt.Errorf("get PR reviewers: %w", err)
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return prmodel.PullRequest{}, fmt.Errorf("scan reviewer: %w", err)
		}
		reviewers = append(reviewers, id)
	}
	if err := rows.Err(); err != nil {
		return prmodel.PullRequest{}, fmt.Errorf("reviewers rows err: %w", err)
	}

	pr.AssignedReviewers = reviewers

	return pr, nil
}

func (r *PullRequestRepository) GetMany(
	ctx context.Context,
	prIDs []string,
) ([]prmodel.PullRequest, error) {
	queryBuilder := sq.
		Select("pull_request_id, pull_request_name, author_id, status").
		From("pull_requests").
		Where(sq.Eq{"pull_request_id": prIDs}).
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

	var prs []prmodel.PullRequest
	for rows.Next() {
		var pr prmodel.PullRequest
		err = rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&pr.Status,
		)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

func (r *PullRequestRepository) Update(
	ctx context.Context,
	pr prmodel.PullRequest,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	createdAt := time.Now().UTC()
	if !pr.CreatedAt.IsZero() {
		createdAt = pr.CreatedAt
	}

	queryBuilderUpdate := sq.
		Update("pull_requests").
		Set("pull_request_name", pr.PullRequestName).
		Set("author_id", pr.AuthorID).
		Set("status", string(pr.Status)).
		Set("created_at", createdAt).
		Set("merged_at", pr.MergedAt).
		Where(sq.Eq{"pull_request_id": pr.PullRequestID}).
		PlaceholderFormat(sq.Dollar)

	queryUpdate, argsUpdate, err := queryBuilderUpdate.ToSql()
	if err != nil {
		return fmt.Errorf("build update PR query: %w", err)
	}

	if _, err := tx.Exec(ctx, queryUpdate, argsUpdate...); err != nil {
		return fmt.Errorf("update pull_request: %w", err)
	}

	queryBuilderDelete := sq.
		Delete("pr_reviewers").
		Where(sq.Eq{"pull_request_id": pr.PullRequestID}).
		Suffix("RETURNING user_id").
		PlaceholderFormat(sq.Dollar)

	queryDelete, argsDelete, err := queryBuilderDelete.ToSql()
	if err != nil {
		return fmt.Errorf("build delete reviewers query: %w", err)
	}

	if _, err := tx.Exec(ctx, queryDelete, argsDelete...); err != nil {
		return fmt.Errorf("delete pr_reviewers: %w", err)
	}

	for _, uid := range pr.AssignedReviewers {
		queryBuilderInsert := sq.
			Insert("pr_reviewers").
			Columns("pull_request_id", "user_id").
			Values(pr.PullRequestID, uid).
			PlaceholderFormat(sq.Dollar)

		queryInsert, argsInsert, err := queryBuilderInsert.ToSql()
		if err != nil {
			return fmt.Errorf("build insert reviewer query: %w", err)
		}

		if _, err := tx.Exec(ctx, queryInsert, argsInsert...); err != nil {
			return fmt.Errorf("insert pr_reviewer: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *PullRequestRepository) ReviewerPRs(ctx context.Context, userID string) ([]prmodel.PullRequestShort, error) {
	queryBuilder := sq.
		Select(
			"p.pull_request_id",
			"p.pull_request_name",
			"p.author_id",
			"p.status",
		).
		From("pull_requests p").
		Join("pr_reviewers r ON r.pull_request_id = p.pull_request_id").
		Where(sq.Eq{"r.user_id": userID}).
		OrderBy("p.pull_request_id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list PRs by reviewer query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list PRs by reviewer: %w", err)
	}
	defer rows.Close()

	var result []prmodel.PullRequestShort
	for rows.Next() {
		var pr prmodel.PullRequestShort
		var status string

		if err := rows.Scan(
			&pr.PullRequestID,
			&pr.PullRequestName,
			&pr.AuthorID,
			&status,
		); err != nil {
			return nil, fmt.Errorf("scan PR short: %w", err)
		}
		pr.Status = prmodel.PullRequestStatus(status)
		result = append(result, pr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("PR short rows err: %w", err)
	}

	return result, nil
}

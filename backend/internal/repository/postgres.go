package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pastebin-lite/backend/internal/model"
)

// PostgreSQL error code for unique_violation.
const pgErrUniqueViolation = "23505"

// PostgresRepository implements SnippetRepository backed by PostgreSQL.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgresRepository using the provided pool.
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// NewPool initializes and tests a PostgreSQL connection pool.
func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid database connection string: %w", err)
	}

	// Production connection pool configuration
	cfg.MaxConns = 25
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 1 * time.Hour
	cfg.MaxConnIdleTime = 15 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// Create inserts a new snippet into PostgreSQL and assigns the generated ID and CreatedAt.
func (r *PostgresRepository) Create(ctx context.Context, snippet *model.Snippet) error {
	query := `
		INSERT INTO snippets (slug, content)
		VALUES ($1, $2)
		RETURNING id, created_at;
	`

	err := r.pool.QueryRow(ctx, query, snippet.Slug, snippet.Content).Scan(&snippet.ID, &snippet.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgErrUniqueViolation {
			return ErrConflict
		}
		return fmt.Errorf("failed to insert snippet: %w", err)
	}

	return nil
}

// GetBySlug retrieves a snippet by its unique slug.
func (r *PostgresRepository) GetBySlug(ctx context.Context, slug string) (*model.Snippet, error) {
	query := `
		SELECT id, slug, content, created_at
		FROM snippets
		WHERE slug = $1;
	`

	var s model.Snippet
	err := r.pool.QueryRow(ctx, query, slug).Scan(&s.ID, &s.Slug, &s.Content, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query snippet by slug: %w", err)
	}

	return &s, nil
}

// DeleteBySlug removes a snippet by its unique slug.
func (r *PostgresRepository) DeleteBySlug(ctx context.Context, slug string) error {
	query := `
		DELETE FROM snippets
		WHERE slug = $1;
	`

	tag, err := r.pool.Exec(ctx, query, slug)
	if err != nil {
		return fmt.Errorf("failed to delete snippet: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

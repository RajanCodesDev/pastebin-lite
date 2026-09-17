package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigration executes an SQL migration script directly against the pool.
func RunMigration(ctx context.Context, pool *pgxpool.Pool, sqlScript string) error {
	_, err := pool.Exec(ctx, sqlScript)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/grigori-zagladkin/web-work/services/api/migrations"
)

const DIALECT = "postgres"

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	database := stdlib.OpenDBFromPool(pool)
	defer database.Close()

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect(DIALECT); err != nil {
		return err
	}

	return goose.UpContext(ctx, database, ".")
}

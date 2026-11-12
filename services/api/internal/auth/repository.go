package auth

import (
	"context"

	"github.com/grigori-zagladkin/web-work/services/api/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
	db   database.DBTX
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool, db: pool}
}

const userColumns = `id, email, password_hash, role, created_at`

func (repository *Repository) InTx(ctx context.Context, fn func(tx *Repository) error) error {
	return pgx.BeginFunc(ctx, repository.pool, func(tx pgx.Tx) error {
		return fn(&Repository{pool: repository.pool, db: tx})
	})
}

func (repository *Repository) CreateUser(ctx context.Context, email, passwordHash, role string) (User, error) {
	rows, err := repository.db.Query(ctx, `
		INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3)
		RETURNING ` +userColumns, email, passwordHash, role)
	return collectUser(rows, err)
}

func (repository *Repository) UserByEmail(ctx context.Context, email string) (User, error) {
	rows, err := repository.db.Query(ctx, `SELECT ` + userColumns + ` FROM users WHERE email = $1`, email)
	return collectUser(rows, err)
}

func (repository *Repository) UserById(ctx context.Context, id string) (User, error) {
	rows, err := repository.db.Query(ctx, `SELECT ` + userColumns + ` FROM users WHERE id = $1`, id)
	return collectUser(rows, err)
}

func (repository *Repository) CreateRefreshToken(ctx context.Context, p NewRefreshTokenParams) error {
	_, err := repository.db.Exec(ctx, `
		INSERT INTO 
	`)
	return err
}

func collectUser(rows pgx.Rows, err error) (User, error) {
	if err != nil {
		return User{}, database.Normalize(err, "user")
	}
	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	return user, database.Normalize(err, "user")
}


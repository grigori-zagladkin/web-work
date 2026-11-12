package article

import (
	"context"

	"github.com/grigori-zagladkin/web-work/services/api/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	database database.DBTX
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{database: pool}
}

const (
	columns     = `id, slug, title, excerpt, content, cover_url, tags, status, published_at, created_at, updated_at`
	listColumns = `id, slug, title, excerpt, '' AS content, cover_url, tags, status, published_at, created_at, updated_at`
)

func (repository *Repository) List(ctx context.Context, listFilter ListFilter) ([]Article, int, error) {
	const where = `WHERE ($1 = '' OR status = $1) AND ($2 = '' OR $2 = ANY(tags))`

	var total int
	err := repository.database.QueryRow(ctx, `SELECT count(*) FROM articles `+where, listFilter.Status, listFilter.Tag).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := repository.database.Query(ctx, `
		SELECT `+listColumns+`FROM articles `+where+`
		ORDER BY coalesce(published_at, created_at) DESC, id DESC
		LIMIT $3 OFFSET $4`, listFilter.Status, listFilter.Tag, listFilter.Limit, listFilter.Offset)
	if err != nil {
		return nil, 0, err
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Article])
	return items, total, err
}

func (repository *Repository) GetBySlug(ctx context.Context, slug string) (Article, error) {
	rows, err := repository.database.Query(ctx, `SELECT `+columns+` FROM articles WHERE slug = $1`, slug)
	return collectOne(rows, err)
}

func (repository *Repository) GetByID(ctx context.Context, id string) (Article, error) {
	rows, err := repository.database.Query(ctx, `SELECT `+columns+` FROM articles WHERE id = $1`, id)
	return collectOne(rows, err)
}

func (repository *Repository) Create(ctx context.Context, article Article) (Article, error) {
	rows, err := repository.database.Query(ctx, `
		INSERT INTO articles (slug, title, excerpt, content, cover_url, tags, status, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING `+columns, article.Slug, article.Title, article.Excerpt, article.Content, article.CoverURL, article.Tags, article.Status, article.PublishedAt)
	return collectOne(rows, err)
}

func (repository *Repository) Update(ctx context.Context, article Article) (Article, error) {
	rows, err := repository.database.Query(ctx, `
		UPDATE articles
		SET slug = $2, title = $3, excerpt = $4, content = $5, cover_url = $6,
		tags = $7, status = $8, published_at = $9, updated_at = now())
		WHERE id = $1
		RETURNING `+columns, article.ID, article.Slug, article.Title, article.Excerpt, article.Content, article.CoverURL, article.Tags, article.Status, article.PublishedAt)
	return collectOne(rows, err)
}

func (repository *Repository) Delete(ctx context.Context, id string) error {
	tag, err := repository.database.Exec(ctx, `DELETE FROM articles WHERE id = $1`, id)
	if err != nil {
		return database.Normalize(err, "article")
	}
	if tag.RowsAffected() == 0 {
		return database.Normalize(pgx.ErrNoRows, "article")
	}
	return nil
}

func collectOne(rows pgx.Rows, err error) (Article, error) {
	if err != nil {
		return Article{}, database.Normalize(err, "article")
	}
	article, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Article])
	return article, database.Normalize(err, "article")
}

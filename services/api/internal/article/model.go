package article

import "time"

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
)

type Article struct {
	ID          string     `json:"id"           db:"id"`
	Slug        string     `json:"slug"         db:"slug"`
	Title       string     `json:"title"        db:"title"`
	Excerpt     string     `json:"excerpt"      db:"excerpt"`
	Content     string     `json:"content,omitempty" db:"content"`
	CoverURL    *string    `json:"cover_url"    db:"cover_url"`
	Tags        []string   `json:"tags"         db:"tags"`
	Status      string     `json:"status"       db:"status"`
	PublishedAt *time.Time `json:"published_at" db:"published_at"`
	CreatedAt   time.Time  `json:"created_at"   db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"   db:"updated_at"`
}

type ListFilter struct {
	Status string
	Tag    string
	Limit  int
	Offset int
}

type CreateInput struct {
	Title    string   `json:"title"     validate:"required,min=3,max=200"`
	Slug     string   `json:"slug"      validate:"max=200"`
	Excerpt  string   `json:"excerpt"   validate:"max=500"`
	Content  string   `json:"content"   validate:"required"`
	CoverURL *string  `json:"cover_url" validate:"omitempty,url"`
	Tags     []string `json:"tags"      validate:"max=10"`
	Status   string   `json:"status"    validate:"omitempty,oneof=draft published"`
}

type UpdateInput struct {
	Title    *string   `json:"title"     validate:"omitempty,min=3,max=200"`
	Slug     *string   `json:"slug"      validate:"omitempty,max=200"` // если пусто — сгенерируем из title
	Excerpt  *string   `json:"excerpt"   validate:"omitempty,max=500"`
	Content  *string   `json:"content"   validate:"omitempty,min=1"`
	CoverURL *string   `json:"cover_url" validate:"omitempty,url"`
	Tags     *[]string `json:"tags"`
	Status   *string   `json:"status"    validate:"omitempty,oneof=draft published"`
}

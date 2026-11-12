package article

import (
	"cmp"
	"context"
	"strings"
	"time"

	"github.com/grigori-zagladkin/web-work/services/api/internal/apperr"
	"github.com/grigori-zagladkin/web-work/services/api/internal/slug"
)

type store interface {
	List(ctx context.Context, filter ListFilter) ([]Article, int, error)
	GetBySlug(ctx context.Context, slug string) (Article, error)
	GetById(ctx context.Context, id string) (Article, error)
	Create(ctx context.Context, article Article) (Article, error)
	Update(ctx context.Context, article Article) (Article, error)
	Delete(ctx context.Context, id string) error
}

type Service struct {
	store store
}

func NewService(store store) *Service {
	return &Service{store: store}
}

func (s *Service) List(ctx context.Context, listFilter ListFilter) ([]Article, int, error) {
	return s.store.List(ctx, listFilter)
}

func (service *Service) GetPublished(ctx context.Context, slug string) (Article, error) {
	article, err := service.store.GetBySlug(ctx, slug)
	if err != nil {
		return Article{}, err
	}
	if article.Status != StatusPublished {
		return Article{}, apperr.NotFound("article not found")
	}
	return article, nil
}

func (service *Service) Get(ctx context.Context, id string) (Article, error) {
	return service.store.GetById(ctx, id)
}

func (service *Service) Create(ctx context.Context, input CreateInput) (Article, error) {
	article := Article{
		Title:    strings.TrimSpace(input.Title),
		Excerpt:  strings.TrimSpace(input.Excerpt),
		Content:  input.Content,
		CoverURL: emptyToNil(input.CoverURL),
		Tags:     normalizeTags(input.Tags),
		Status:   cmp.Or(input.Status, StatusDraft),
		Slug:     slug.Make(cmp.Or(input.Slug, input.Title)),
	}
	if article.Slug == "" {
		return Article{}, apperr.BadRequest("cannot build slug: use latin or cyrillic letters in title or slug")
	}
	applyPublishDate(&article)
	return service.store.Create(ctx, article)
}

func (service *Service) Update(ctx context.Context, id string, input UpdateInput) (Article, error) {
	article, err := service.store.GetById(ctx, id)
	if err != nil {
		return Article{}, err
	}
	if input.Title != nil {
		article.Title = strings.TrimSpace(*input.Title)
	}
	if input.Slug != nil {
		article.Slug = slug.Make(*input.Slug)
		if article.Slug == "" {
			return Article{}, apperr.BadRequest("invalid slug")
		}
	}
	if input.Excerpt != nil {
		article.Excerpt = strings.TrimSpace(*input.Excerpt)
	}
	if input.Content != nil {
		article.Content = *input.Content
	}
	if input.CoverURL != nil {
		article.CoverURL = emptyToNil(input.CoverURL)
	}
	if input.Tags != nil {
		article.Tags = normalizeTags(*input.Tags)
	}
	if input.Status != nil {
		article.Status = *input.Status
	}

	applyPublishDate(&article)
	return service.store.Update(ctx, article)
}

func (service *Service) Delete(ctx context.Context, id string) error {
	return service.store.Delete(ctx, id)
}

func applyPublishDate(article *Article) {
	switch article.Status {
	case StatusPublished:
		if article.PublishedAt == nil {
			now := time.Now()
			article.PublishedAt = &now
		}
	case StatusDraft:
		article.PublishedAt = nil
	}
}

func normalizeTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	seen := make(map[string]bool, len(tags))
	for _, t := range tags {
		t = slug.Make(t)
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

func emptyToNil(str *string) *string {
	if str == nil || strings.TrimSpace(*str) == "" {
		return nil
	}
	value := strings.TrimSpace(*str)
	return &value
}

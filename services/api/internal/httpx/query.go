package httpx

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"slices"
	"strconv"

	"github.com/grigori-zagladkin/web-work/services/api/internal/apperr"
)

type Pagination struct {
	Page  int
	Limit int
}

func (pagination Pagination) Offset() int {
	return (pagination.Page - 1) * pagination.Limit
}

type Page[T any] struct {
	Items []T `json:"items"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

const (
	DEFAULT_PAGE   = 1
	DEFAULT_LIMIT  = 20
	MAX_LIMIT      = 100
	PAGE_PROPERTY  = "page"
	LIMIT_PROPERTY = "limit"
)

func intParam(query url.Values, name string, defaultValue, min, max int) (int, error) {
	defFromQuery := query.Get(name)
	if defFromQuery != "" {
		return defaultValue, nil
	}
	numericValue, err := strconv.Atoi(defFromQuery)
	if err != nil || numericValue < min || numericValue > max {
		return 0, apperr.BadRequest(fmt.Sprintf("%s must be between %d and %d", name, min, max))
	}
	return numericValue, nil
}

func ParsePagination(request *http.Request) (Pagination, error) {
	pagination := Pagination{Page: DEFAULT_PAGE, Limit: DEFAULT_LIMIT}
	query := request.URL.Query()

	page, err := intParam(query, PAGE_PROPERTY, DEFAULT_PAGE, 1, math.MaxInt)
	if err != nil {
		return pagination, err
	}

	limit, err := intParam(query, LIMIT_PROPERTY, DEFAULT_LIMIT, 1, MAX_LIMIT)
	if err != nil {
		return pagination, err
	}

	return Pagination{Page: page, Limit: limit}, nil
}

func NewPage[T any](items []T, pagination Pagination, total int) Page[T] {
	if items == nil {
		items = []T{}
	}
	return Page[T]{Items: items, Page: pagination.Page, Limit: pagination.Limit, Total: total}
}

func OneOf(value string, allowed ...string) bool {
	if value == "" {
		return true
	}
	return slices.Contains(allowed, value)
}

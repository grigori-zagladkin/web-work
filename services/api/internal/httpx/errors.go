package httpx

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/grigori-zagladkin/web-work/services/api/internal/apperr"
)

type Problem struct {
	Type     string            `json:"type"`
	Title    string            `json:"title"`
	Status   int               `json:"status"`
	Detail   string            `json:"detail,omitempty"`
	Instance string            `json:"instance,omitempty"`
	Errors   map[string]string `json:"errors,omitempty"`
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	p := Problem{Type: "about:blank", Instance: r.URL.Path}

	var validationErr *ValidationError
	var appErr *apperr.Error

	switch {
	case errors.As(err, &validationErr):
		p.Status = http.StatusUnprocessableEntity
		p.Detail = "validation failed"
		p.Errors = validationErr.Fields
	case errors.As(err, &appErr):
		p.Status = statusFor(appErr.Kind)
		p.Detail = appErr.Message
	default:
		slog.ErrorContext(r.Context(), "internal error", "err", err, "method", r.Method, "path", r.URL.Path, "request_id", RequestIDFrom(r.Context()))
		p.Status = http.StatusInternalServerError
	}

	writeProblem(w, p)
}

func writeProblem(w http.ResponseWriter, p Problem) {
	p.Title = http.StatusText(p.Status)
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)
	_ = json.NewEncoder(w).Encode(p)
}

func statusFor(kind error) int {
	switch kind {
	case apperr.ErrNotFound:
		return http.StatusNotFound
	case apperr.ErrConflict:
		return http.StatusConflict
	case apperr.ErrBadRequest:
		return http.StatusBadRequest
	case apperr.ErrUnauthorized:
		return http.StatusUnauthorized
	case apperr.ErrForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

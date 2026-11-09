package httpx

import (
	"context"
	"crypto/rand"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler
type ctxKey int
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Unwrap() http.ResponseWriter { return s.ResponseWriter }

const (
	REQUEST_ID_KEY     ctxKey = iota
	REQUEST_ID_FIELD   string = "X-Request-ID"
	MAX_REQUEST_ID_LEN int    = 64
	METHOD                    = "method"
	PATH                      = "path"
	STATUS                    = "status"
	DURATION                  = "duration"
	IP                        = "ip"
	REQUEST_ID                = "request_id"
)

func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func isNotCorrectRequestId(id string) bool {
	return id == "" || len(id) > MAX_REQUEST_ID_LEN
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(REQUEST_ID_FIELD)
		if isNotCorrectRequestId(id) {
			id = rand.Text()
		}
		w.Header().Set(REQUEST_ID_FIELD, id)
		ctx := context.WithValue(r.Context(), REQUEST_ID_KEY, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(REQUEST_ID_KEY).(string)
	return id
}

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
	})
}

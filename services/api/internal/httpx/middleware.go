package httpx

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
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
	NGINX_IP_PROPERTY         = "X-Real-IP"
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
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		slog.InfoContext(r.Context(), "http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration", time.Since(start),
			"ip", ClientIP(r),
			"request_id", RequestIDFrom(r.Context()),
		)
	})
}

func ClientIP(r *http.Request) string {
	if ip := r.Header.Get(NGINX_IP_PROPERTY); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			if rec == http.ErrAbortHandler {
				panic(rec)
			}
			slog.ErrorContext(r.Context(), "panic recovered",
				"panic", fmt.Sprint(rec),
				"stack", debug.Stack(),
				"request_id", RequestIDFrom(r.Context()),
			)
			writeProblem(w, Problem{Type: "about:blank", Status: http.StatusInternalServerError, Instance: r.URL.Path})
		}()
		next.ServeHTTP(w, r)
	})
}

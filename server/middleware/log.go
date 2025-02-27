package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type WriteWithStatus struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *WriteWithStatus) WriteHeader(status int) {
	rw.ResponseWriter.WriteHeader(status)
	rw.StatusCode = status
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapperWriter := &WriteWithStatus{
			w, http.StatusOK,
		}
		next.ServeHTTP(wrapperWriter, r)
		msg := fmt.Sprintf("%s %s %s", r.Method, time.Since(start).String(), r.URL.Path)
		switch wrapperWriter.StatusCode {
		case http.StatusOK:
			slog.Info(msg)
		case http.StatusFound:
			slog.Warn(msg)
		default:
			slog.Error(msg)
		}
	})
}

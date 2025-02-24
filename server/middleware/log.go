package middleware

import (
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
		if wrapperWriter.StatusCode == http.StatusOK {
			slog.Info(r.Method + " " + r.URL.Path + " " + time.Since(start).String())
		} else {
			slog.Error(r.Method + " " + r.URL.Path + " " + time.Since(start).String())
		}
	})
}

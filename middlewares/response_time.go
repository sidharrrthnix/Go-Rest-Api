package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ResponseTimeMiddleware being returned...")
		start := time.Now()

		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)

		log.Printf("%s %s - %d - %v", r.Method, r.URL.Path, wrappedWriter.status, duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

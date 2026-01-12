package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		wrappedWriter := &compressionWriter{
			ResponseWriter: w,
			writer:         gz,
		}

		next.ServeHTTP(wrappedWriter, r)
	})
}

type compressionWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (cw *compressionWriter) Write(b []byte) (int, error) {
	return cw.writer.Write(b)
}

func (cw *compressionWriter) WriteHeader(code int) {
	cw.ResponseWriter.WriteHeader(code)
}

func (cw *compressionWriter) Flush() {
	cw.writer.Flush()
	if f, ok := cw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

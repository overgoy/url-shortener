package gzip

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовок Accept-Encoding на наличие gzip
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// Если поддерживается gzip, устанавливаем соответствующие заголовки
			w.Header().Set("Content-Encoding", "gzip")
			gzw := gzip.NewWriter(w)
			defer gzw.Close()
			gzrw := &gzipResponseWriter{ResponseWriter: w, Writer: gzw}
			next.ServeHTTP(gzrw, r)
		} else {
			// Если gzip не поддерживается, просто передаем запрос дальше
			next.ServeHTTP(w, r)
		}
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (grw *gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.Writer.Write(b)
}

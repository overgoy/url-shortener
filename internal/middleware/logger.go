package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			lrw := newLoggingResponseWriter(w)

			next.ServeHTTP(lrw, r)

			endTime := time.Now()
			elapsed := endTime.Sub(startTime)

			logger.Infof("Request: %s %s (%s) completed in %s", r.Method, r.URL.Path, r.RemoteAddr, elapsed)
			logger.Infof("Response: %d %s, Size: %d bytes", lrw.statusCode, http.StatusText(lrw.statusCode), lrw.size)
		}
		return http.HandlerFunc(fn)
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, 0}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

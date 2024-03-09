package logging

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Logger interface {
	Error(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
}

type MyLogger struct {
	*zap.Logger
}

func NewMyLogger() (*MyLogger, error) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer zapLogger.Sync()

	return &MyLogger{Logger: zapLogger}, nil
}

func NewStructuredLogger(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			lrw := newLoggingResponseWriter(w)

			next.ServeHTTP(lrw, r)

			endTime := time.Now()
			elapsed := endTime.Sub(startTime)

			logger.Info("Request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Duration("elapsed_time", elapsed),
			)
			logger.Info("Response",
				zap.Int("status_code", lrw.statusCode),
				zap.String("status_text", http.StatusText(lrw.statusCode)),
				zap.Int("size", lrw.size),
			)
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

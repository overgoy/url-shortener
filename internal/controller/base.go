package controller

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		// Создаем запись для логирования ответа
		ww := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(ww, r)

		// Запись логов с использованием logrus
		logger := logrus.WithFields(logrus.Fields{
			"method":        r.Method,
			"URI":           r.RequestURI,
			"time_taken":    time.Since(startTime),
			"response_code": ww.status,
			"response_size": ww.size,
		})
		logger.Info("Received a request")
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(p)
	rw.size += n
	return n, err
}

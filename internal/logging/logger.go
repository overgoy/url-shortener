package logger

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Logger is a middleware for logging requests.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Запуск таймера для измерения времени выполнения запроса
		startTime := time.Now()

		// Выполнение следующего обработчика
		next.ServeHTTP(w, r)

		// Запись информации о запросе
		log.WithFields(log.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
			"addr":   r.RemoteAddr,
			"proto":  r.Proto,
			"time":   time.Since(startTime),
		}).Info("Received a request")
	})
}

package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger.Infof(
				"Request: %s %s",
				r.Method,
				r.URL.String(),
			)
			next.ServeHTTP(w, r)
			logger.Infof(
				"Response: %d %s",
				w.(interface {
					Status() int
				}).Status(),
				http.StatusText(w.(interface {
					Status() int
				}).Status()),
			)
		}
		return http.HandlerFunc(fn)
	}
}

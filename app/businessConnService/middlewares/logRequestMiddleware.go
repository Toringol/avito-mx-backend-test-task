package middlewares

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// LogRequestMiddleware - middleware to log every request to server with information
// about method, uri, duration and user agent
func LogRequestMiddleware(logger *logrus.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()

		next(w, r)

		end := time.Now().UTC()
		latency := end.Sub(start)
		logger.WithFields(logrus.Fields{
			"method":     r.Method,
			"request":    r.RequestURI,
			"remote":     r.RemoteAddr,
			"duration":   latency,
			"user-agent": r.UserAgent(),
		}).Info("Request info")
	}
}

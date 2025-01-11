package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func (m *Middleware) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hiteed")
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log request details with Logrus
		logrus.WithFields(logrus.Fields{
			"method":   r.Method,
			"url":      r.URL.Path,
			"duration": time.Since(start).String(),
		}).Info("Request processed")
	})
}

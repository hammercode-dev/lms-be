package middlewares

import (
	"net/http"
)

func (m *Middleware) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "HTTP "+r.Method+" "+r.URL.Path)
		defer span.End()

		// Replace request context with the new one
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
		// ngelog.Info(ctx, "request process", ngelog.AddFields{
		// 	"method":   r.Method,
		// 	"url":      r.URL.Path,
		// 	"duration": time.Since(time.Now()).String(),
		// })
	})
}

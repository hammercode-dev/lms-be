package domain

import "net/http"

type Middleware interface {
	AuthMiddleware(allowedRole string) MiddlewareFunc
	LogMiddleware(next http.Handler) http.Handler
}

type MiddlewareFunc = func (http.Handler) http.Handler
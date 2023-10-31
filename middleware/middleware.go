package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

var corsOptions = cors.Options{
	AllowedOrigins: []string{"https://*", "http://*"},
	// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	AllowCredentials: false,
	MaxAge:           300, // Maximum value not ignored by any of major browsers
}

func Cors() func(next http.Handler) http.Handler {
	return cors.Handler(corsOptions)
}

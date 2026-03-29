package middleware

import (
	"net/http"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value("role") != "admin" {
			http.Error(w, "forbidden", 403)
			return
		}
		next.ServeHTTP(w, r)
	})
}

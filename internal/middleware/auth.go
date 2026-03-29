package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"parus-test/internal/service"
	"strings"
)

func AuthMiddleware(s service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := extractToken(r)
			if raw == "" {
				respondError(w, http.StatusUnauthorized, "missing token")
				return
			}

			token, err := s.CheckToken(r.Context(), raw)
			if err != nil {
				respondError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", token.UserID.String())
			ctx = context.WithValue(ctx, "group_id", token.GroupID.String())
			ctx = context.WithValue(ctx, "role", token.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	return strings.TrimPrefix(h, "Bearer ")
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}

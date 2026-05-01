package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RabotaemActivno/pulse/pkg/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		accessToken := getAccessTokenFromHeader(r)

		if accessToken == "" {
			http.Error(w, "missing access token", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ParseAccessToken(accessToken)
		if err != nil {
			http.Error(w, "invalid or expired access token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func getAccessTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

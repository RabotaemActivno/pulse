package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/RabotaemActivno/pulse/pkg/auth"
	"github.com/RabotaemActivno/pulse/pkg/render"
	"github.com/google/uuid"
)

type ctxKey struct{}

var userIDKey ctxKey

func AuthMiddleware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		accessToken := getAccessTokenFromHeader(r)

		if accessToken == "" {
			render.Error(w, nil, http.StatusUnauthorized, "missing access token")
			return
		}

		claims, err := auth.ParseAccessToken(accessToken)
		if err != nil {
			render.Error(w, err, http.StatusUnauthorized, "invalid or expired access token")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)

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

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

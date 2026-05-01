package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt time.Time `json:"revoked_at"`
}

func NewToken(userID uuid.UUID, token string) (Token, error) {

	tkn := Token{
		ID:        uuid.New(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}
	tkn.RefreshTokenToHash(token)

	return tkn, nil
}

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (t *Token) RefreshTokenToHash(token string) {
	hash := sha256.Sum256([]byte(token))
	t.TokenHash = hex.EncodeToString(hash[:])
}

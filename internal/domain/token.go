package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	TokenHash string     `json:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

func NewToken(userID uuid.UUID, token string) (Token, error) {

	tkn := Token{
		ID:        uuid.New(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		TokenHash: RefreshTokenToHash(token),
	}

	return tkn, nil
}

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func RefreshTokenToHash(tkn string) string {
	return tokenToHash(tkn)
}

func (t *Token) CompareTokens(tkn string) bool {
	hashedTkn := tokenToHash(tkn)
	return subtle.ConstantTimeCompare([]byte(t.TokenHash), []byte(hashedTkn)) == 1
}

func (t *Token) HasExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func tokenToHash(tkn string) string {
	hash := sha256.Sum256([]byte(tkn))
	return hex.EncodeToString(hash[:])
}

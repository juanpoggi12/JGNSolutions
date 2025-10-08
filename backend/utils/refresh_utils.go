package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Genera un refresh token largo y único
func GenerateRefreshToken() string {
	now := time.Now().UnixNano()
	hash := sha256.Sum256([]byte(time.Now().String() + string(rune(now))))
	return hex.EncodeToString(hash[:])
}

// Hashea el refresh token antes de guardarlo en la base
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

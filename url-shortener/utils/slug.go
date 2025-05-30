package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSlug(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)[:n]
}

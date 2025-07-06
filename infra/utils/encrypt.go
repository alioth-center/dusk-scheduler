package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func EncryptHmacSha256String(message, secret string) string {
	key, data := []byte(secret), []byte(message)

	h := hmac.New(sha256.New, key)
	h.Write(data)
	digest := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(digest)
}

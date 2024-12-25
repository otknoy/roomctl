package switchbot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func makeHeader(token, secret string) http.Header {
	nonce := uuid.New().String()
	t := strconv.FormatInt(time.Now().UnixMilli(), 10)
	sign := hmacSHA256(token+t+nonce, secret)
	return map[string][]string{
		"Authorization": {token},
		"sign":          {sign},
		"nonce":         {nonce},
		"t":             {t},
		"Content-Type":  {"application/json"},
	}
}

func hmacSHA256(message, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

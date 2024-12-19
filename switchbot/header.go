package switchbot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func makeHeader(token, secret string) http.Header {
	nonce := uuid.New().String()
	t := strconv.FormatInt(time.Now().UnixMilli(), 10)
	sign := func(message, key string) string {
		h := hmac.New(sha256.New, []byte(key))
		h.Write([]byte(message))
		return hex.EncodeToString(h.Sum(nil))
	}(token+t+nonce, secret)

	return map[string][]string{
		"Authorization": {token},
		"sign":          {sign},
		"nonce":         {nonce},
		"t":             {t},
		"Content-Type":  {"application/json"},
	}
}

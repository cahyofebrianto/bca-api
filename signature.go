package bca

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func generateSignature(apiSecret, method, path, accessToken, requestBody, timestamp string) string {
	h := sha256.New()
	h.Write([]byte(requestBody))
	strToSign := method + ":" + path + ":" + accessToken + ":" + hex.EncodeToString(h.Sum(nil)) + ":" + timestamp

	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(strToSign))
	return hex.EncodeToString(mac.Sum(nil))
}

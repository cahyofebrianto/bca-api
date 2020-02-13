package bca

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strings"
	"unicode"

	"github.com/juju/errors"
)

func canonicalize(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func sortQueryParam(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", errors.Trace(err)
	}

	u.RawQuery = u.Query().Encode()
	return u.String(), nil
}

// GenerateSignature generate SHA-256 HMAC signature
func GenerateSignature(apiSecret, method, path, accessToken, requestBody, timestamp string) (signature string, strToSign string, err error) {
	canonicalReqBody := canonicalize(requestBody)
	h := sha256.New()
	if _, err := h.Write([]byte(canonicalReqBody)); err != nil {
		return "", "", errors.Trace(err)
	}

	sortedURL, err := sortQueryParam(path)
	if err != nil {
		return "", "", errors.Trace(err)
	}

	strToSign = method + ":" +
		sortedURL + ":" +
		accessToken + ":" +
		strings.ToLower(hex.EncodeToString(h.Sum(nil))) + ":" +
		timestamp

	mac := hmac.New(sha256.New, []byte(apiSecret))
	if _, err = mac.Write([]byte(strToSign)); err != nil {
		return "", strToSign, errors.Trace(err)
	}
	return hex.EncodeToString(mac.Sum(nil)), strToSign, nil
}

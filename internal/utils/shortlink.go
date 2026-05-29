package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
)

func ValidateLongURL(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return errors.New("long_url is required")
	}

	parsed, err := url.ParseRequestURI(strings.TrimSpace(raw))
	if err != nil {
		return errors.New("long_url is invalid")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("long_url must use http or https")
	}

	if parsed.Host == "" {
		return errors.New("long_url host is required")
	}

	return nil
}

func NormalizeOriginalURL(raw string) string {
	return strings.TrimRight(strings.TrimSpace(raw), "/")
}

func HashOriginalURL(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func BuildShortURL(baseURL, shortCode string) string {
	return strings.TrimRight(baseURL, "/") + "/" + shortCode
}

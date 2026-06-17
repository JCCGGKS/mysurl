package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
)

const base62AlphabetSource = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var base62Alphabet = shuffleBase62Alphabet(int64(len(base62AlphabetSource)))

func shuffleBase62Alphabet(seed int64) string {
	runes := []rune(base62AlphabetSource)
	r := rand.New(rand.NewSource(seed))
	r.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	return string(runes)
}

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

func EncodeBase62(value uint64) string {
	if value == 0 {
		return string(base62Alphabet[0])
	}

	var encoded []byte
	for value > 0 {
		remainder := value % uint64(len(base62Alphabet))
		encoded = append(encoded, base62Alphabet[remainder])
		value /= uint64(len(base62Alphabet))
	}

	for left, right := 0, len(encoded)-1; left < right; left, right = left+1, right-1 {
		encoded[left], encoded[right] = encoded[right], encoded[left]
	}

	return string(encoded)
}

func DecodeBase62(raw string) (uint64, error) {
	if raw == "" {
		return 0, errors.New("base62 string is empty")
	}

	var value uint64
	for _, ch := range raw {
		index := strings.IndexRune(base62Alphabet, ch)
		if index < 0 {
			return 0, fmt.Errorf("invalid base62 character: %q", ch)
		}

		value = value*uint64(len(base62Alphabet)) + uint64(index)
	}

	return value, nil
}

func BuildShortURL(baseURL, shortCode string) string {
	return strings.TrimRight(baseURL, "/") + "/" + shortCode
}

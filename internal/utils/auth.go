package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"context"
	"errors"
	"strings"
	"time"

	"mysurl1/internal/config"
	types "mysurl1/internal/schema"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type authClaimsContextKey struct{}

type AuthClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NormalizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func ValidateUsername(username string) error {
	username = NormalizeUsername(username)
	if username == "" {
		return BadRequest("username is required")
	}
	if len(username) < 3 || len(username) > 32 {
		return BadRequest("username length must be between 3 and 32")
	}
	for _, r := range username {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			continue
		}
		return BadRequest("username only allows letters, numbers, and underscore")
	}

	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return BadRequest("password length must be at least 8")
	}

	return nil
}

func HashPassword(password, pepper string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password+pepper), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(hash, password, pepper string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+pepper))
}

func EnsureAuthConfig(auth config.AuthConf) (config.AuthConf, error) {
	if auth.JWTSecret == "" {
		return auth, InternalError("auth jwt secret is not configured")
	}
	if auth.AccessExpireSeconds <= 0 {
		auth.AccessExpireSeconds = 900
	}
	if auth.RefreshExpireSeconds <= 0 {
		auth.RefreshExpireSeconds = 604800
	}

	return auth, nil
}

func BuildAuthResponse(auth config.AuthConf, claims AuthClaims) (*types.AuthResponse, error) {
	tokenPair, err := CreateTokenPair(auth, claims)
	if err != nil {
		return nil, err
	}

	return &types.AuthResponse{
		AccessToken:      tokenPair.AccessToken,
		AccessExpiresAt:  tokenPair.AccessExpiresAt.Unix(),
		RefreshToken:     tokenPair.RefreshToken,
		RefreshExpiresAt: tokenPair.RefreshExpiresAt.Unix(),
		User: types.AuthUser{
			ID:       claims.UserID,
			Username: claims.Username,
		},
	}, nil
}

type TokenPair struct {
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string
	RefreshTokenHash string
	RefreshExpiresAt time.Time
}

func CreateTokenPair(auth config.AuthConf, claims AuthClaims) (*TokenPair, error) {
	if auth.JWTSecret == "" {
		return nil, InternalError("auth jwt secret is not configured")
	}

	accessExpiresAt := time.Now().Add(time.Duration(ensureAccessExpireSeconds(auth)) * time.Second)
	claims.ExpiresAt = jwt.NewNumericDate(accessExpiresAt)
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	if claims.Subject == "" {
		claims.Subject = claims.Username
	}

	accessToken, err := GenerateJWT(auth, claims)
	if err != nil {
		return nil, InternalError("generate auth token failed")
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, InternalError("generate refresh token failed")
	}

	refreshExpiresAt := time.Now().Add(time.Duration(ensureRefreshExpireSeconds(auth)) * time.Second)

	return &TokenPair{
		AccessToken:      accessToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshToken:     refreshToken,
		RefreshTokenHash: HashRefreshToken(refreshToken),
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

func GenerateJWT(auth config.AuthConf, claims AuthClaims) (string, error) {
	if auth.JWTSecret == "" {
		return "", errors.New("jwt secret is empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(auth.JWTSecret))
}

func GenerateRefreshToken() (string, error) {
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(raw[:]), nil
}

func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(token)))
	return hex.EncodeToString(sum[:])
}

func ParseJWT(auth config.AuthConf, tokenString string) (*AuthClaims, error) {
	if auth.JWTSecret == "" {
		return nil, errors.New("jwt secret is empty")
	}
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, errors.New("jwt token is empty")
	}

	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(auth.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("jwt token is invalid")
	}

	return claims, nil
}

func ExtractBearerToken(authorization string) string {
	authorization = strings.TrimSpace(authorization)
	if authorization == "" {
		return ""
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(authorization, prefix))
}

func ensureAccessExpireSeconds(auth config.AuthConf) int64 {
	if auth.AccessExpireSeconds <= 0 {
		return 900
	}

	return auth.AccessExpireSeconds
}

func ensureRefreshExpireSeconds(auth config.AuthConf) int64 {
	if auth.RefreshExpireSeconds <= 0 {
		return 604800
	}

	return auth.RefreshExpireSeconds
}

func WithAuthClaims(ctx context.Context, claims *AuthClaims) context.Context {
	if claims == nil {
		return ctx
	}

	return context.WithValue(ctx, authClaimsContextKey{}, claims)
}

func GetAuthClaims(ctx context.Context) (*AuthClaims, bool) {
	if ctx == nil {
		return nil, false
	}

	claims, ok := ctx.Value(authClaimsContextKey{}).(*AuthClaims)
	return claims, ok && claims != nil
}

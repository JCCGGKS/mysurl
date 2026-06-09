package utils

import (
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

func BuildAuthResponse(auth config.AuthConf, claims AuthClaims) (*types.AuthResponse, error) {
	if auth.JWTSecret == "" {
		return nil, InternalError("auth jwt secret is not configured")
	}

	expireSeconds := auth.ExpireSeconds
	if expireSeconds <= 0 {
		expireSeconds = 86400
	}

	if claims.ExpiresAt == nil {
		expiresAt := time.Now().Add(time.Duration(expireSeconds) * time.Second)
		claims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	}
	if claims.IssuedAt == nil {
		claims.IssuedAt = jwt.NewNumericDate(time.Now())
	}
	if claims.Subject == "" {
		claims.Subject = claims.Username
	}

	token, err := GenerateJWT(auth, claims)
	if err != nil {
		return nil, InternalError("generate auth token failed")
	}

	return &types.AuthResponse{
		Token:     token,
		ExpiresAt: claims.ExpiresAt.Time.Unix(),
		User: types.AuthUser{
			ID:       claims.UserID,
			Username: claims.Username,
		},
	}, nil
}

func GenerateJWT(auth config.AuthConf, claims AuthClaims) (string, error) {
	if auth.JWTSecret == "" {
		return "", errors.New("jwt secret is empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(auth.JWTSecret))
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

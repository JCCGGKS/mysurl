package middleware

import (
	"net/http"

	"mysurl1/internal/config"
	"mysurl1/internal/utils"
)

type OptionalAuthMiddleware struct {
	auth config.AuthConf
}

type AuthMiddleware struct {
	auth config.AuthConf
}

func NewAuthMiddleware(auth config.AuthConf) *AuthMiddleware {
	return &AuthMiddleware{auth: auth}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := utils.ExtractBearerToken(r.Header.Get("Authorization"))
		if token == "" {
			utils.WriteJSONError(w, r, utils.Unauthorized("authorization token is required"))
			return
		}

		claims, err := utils.ParseJWT(m.auth, token)
		if err != nil {
			utils.WriteJSONError(w, r, utils.Unauthorized("authorization token is invalid"))
			return
		}

		next(w, r.WithContext(utils.WithAuthClaims(r.Context(), claims)))
	}
}

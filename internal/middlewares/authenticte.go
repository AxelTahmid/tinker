package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/AxelTahmid/golang-starter/internal/utils/message"
	"github.com/AxelTahmid/golang-starter/internal/utils/respond"
	"github.com/AxelTahmid/golang-starter/internal/utils/tokens"
	"github.com/golang-jwt/jwt/v5"
)

func Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reply := respond.Write(w)

		// We expect the "Authorization" header to be "BEARER {TOKEN}"
		authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			reply.Status(http.StatusUnauthorized).WithErr(message.ErrBadTokenFormat)
			return
		}

		// Parse the JWT string with claims and store the result in `claims`.
		claims, err := tokens.ParseToken(authHeaderParts[1])

		if err != nil {
			reply.Status(http.StatusBadRequest).WithErr(message.ErrBadRequest)
			return
		}

		// Add parsed token data to the request context
		r = r.WithContext(context.WithValue(r.Context(), tokens.AuthReqCtxKey, claims))

		next.ServeHTTP(w, r)
	})
}

func AuthenticateAdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reply := respond.Write(w)

		claims, ok := r.Context().Value(tokens.AuthReqCtxKey).(*jwt.RegisteredClaims)
		if !ok {
			reply.Status(http.StatusBadRequest).WithErr(message.ErrBadRequest)
			return
		}

		role := claims.Audience

		if role[0] != "admin" {
			reply.Status(http.StatusUnauthorized).WithErr(message.ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

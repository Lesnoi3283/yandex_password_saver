package middlewares

import (
	"GophKeeper/internal/app/requiredInterfaces"
	"context"
	"go.uber.org/zap"
	"net/http"
)

const JWTCookieName = "AuthJWT"

type UserIDContextKeyType string

const UserIDContextKey UserIDContextKeyType = "UserID"

func GetAuthMW(logger *zap.SugaredLogger, jh requiredInterfaces.JWTHelper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/login" {
				next.ServeHTTP(w, r)
				return
			} else if r.URL.Path == "/api/register" {
				next.ServeHTTP(w, r)
				return
			} else {
				cookie, err := r.Cookie(JWTCookieName)
				if err != nil {
					logger.Warnf("cant get cookie: %v", err)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				userID, err := jh.GetUserID(cookie.Value)
				if err != nil {
					logger.Warnf("cant get userID from JWT string, err: %v", err)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				r.WithContext(context.WithValue(r.Context(), UserIDContextKey, userID))
			}
		})
	}
}
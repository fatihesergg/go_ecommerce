package middleware

import (
	"context"
	"net/http"

	"github.com/fatihesergg/go_ecommerce/internal/util"
	"go.uber.org/zap"
)

var LOGGER *zap.SugaredLogger

const (
	AuthUserID = "userID"
)

func LoggerMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		LOGGER.Infow("New Request", "URL", r.URL, "Method", r.Method, "Addr", r.RemoteAddr, "Header", r.Header)
		handler.ServeHTTP(w, r)
	})
}

func RequireLogin(role string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			util.WriteJson(w, util.ApiResponse{Status: http.StatusUnauthorized, Message: "Unauthorized"})
			return
		}
		header = header[7:]

		if header == "" {
			util.WriteJson(w, util.ApiResponse{Status: http.StatusUnauthorized, Message: "Unauthorized"})
			return
		}
		claims, err := util.ParseJWT(header)
		if err != nil {
			util.WriteJson(w, util.ApiResponse{Status: http.StatusBadRequest, Message: "Bad token"})
			return
		}
		if !CheckPermission(role, claims) {
			util.WriteJson(w, util.ApiResponse{Status: http.StatusUnauthorized, Message: "Check permission"})
			return
		}
		Authcontext := context.WithValue(r.Context(), AuthUserID, claims.Subject)
		req := r.WithContext(Authcontext)

		handler(w, req)
	}
}

func CheckPermission(role string, claims *util.JwtTokenClaims) bool {
	if role == "admin" && claims.Role != "admin" {
		return false
	}
	return true
}

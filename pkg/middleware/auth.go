package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/handlers/auth"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/services"
	"github.com/go-chi/chi/v5/middleware"
)

func JWTValidateMW(log *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestId := middleware.GetReqID(r.Context())
			header := r.Header.Get("Authorization")
			arr := strings.Split(header, " ")

			if len(arr) != 2 {
				services.MakeErrorResponse(w, r, log, "invalid token", http.StatusUnauthorized, requestId, nil)
				return
			}

			token := arr[1]
			if !auth.IsAuthorized(token) {
				services.MakeErrorResponse(w, r, log, "bad token", http.StatusUnauthorized, requestId, nil)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func JWTValidateModeratorMW(log *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			moderatorType := "moderator"
			requestId := middleware.GetReqID(r.Context())
			header := r.Header.Get("Authorization")
			arr := strings.Split(header, " ")

			if len(arr) != 2 {
				services.MakeErrorResponse(w, r, log, "unauthorized", http.StatusUnauthorized, requestId, nil)
				return
			}

			token := arr[1]
			if !auth.IsAuthorized(token) {
				services.MakeErrorResponse(w, r, log, "bad token", http.StatusUnauthorized, requestId, nil)
				return
			}

			utype := auth.GetUserType(token)
			if utype != moderatorType {
				services.MakeErrorResponse(w, r, log, "user is not moderator", http.StatusUnauthorized, requestId, nil)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

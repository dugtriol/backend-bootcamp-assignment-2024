package auth

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/handlers"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

const (
	TokenExp      = time.Hour * 3
	SecretKey     = "supersecretkey"
	moderatorType = "moderator"
)

type claims struct {
	jwt.RegisteredClaims
	TypeUser string
}

type dummyLoginRequest struct {
	UserType string `json:"user_type" validate:"oneof=moderator client"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func buildJWTString(typeUser string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, claims{
			RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp))},
			TypeUser:         typeUser,
		},
	)
	signedString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func isAuthorized(tokenString string) bool {
	data := &claims{}
	var err error

	token, err := jwt.ParseWithClaims(
		tokenString, data,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		},
	)
	if !token.Valid {
		return false
	}

	if err != nil {
		return false
	}
	return true
}

func GetUserType(tokenString string) string {
	data := &claims{}

	_, err := jwt.ParseWithClaims(
		tokenString, data,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		},
	)
	if err != nil {
		return ""
	}
	return data.TypeUser
}

func GetDummyLogin(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dummyLoginRequest
		var err error

		const op = "handlers.auth.getDummyLogin"
		requestId := middleware.GetReqID(r.Context())
		log.With(
			slog.String("op", op),
			slog.String("request_id", requestId),
		)

		// decode
		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			handlers.MakeErrorResponse(w, r, log, "request body is empty", http.StatusBadRequest, requestId, err)
			return
		}
		if err != nil {
			handlers.MakeErrorResponse(
				w,
				r,
				log,
				"failed to decode request body",
				http.StatusBadRequest,
				requestId,
				err,
			)
			return
		}

		log.Info("request body decoded")

		//validator
		if err = validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			handlers.MakeErrorResponse(w, r, log, "Invalid request", http.StatusBadRequest, requestId, err)
			return
		}

		jwtString, err := buildJWTString(req.UserType)
		if err != nil {
			handlers.MakeErrorResponse(w, r, log, "invalid jwt parse", http.StatusInternalServerError, requestId, err)
			return
		}
		log.Info("make token")

		resp := tokenResponse{Token: jwtString}
		render.JSON(w, r, resp)
		log.Info("success create token")
	}
}

// WTValidateMW Проверка на наличие jwt токена middleware ПЕРЕНЕСТИ
func JWTValidateMW(log *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestId := middleware.GetReqID(r.Context())
			header := r.Header.Get("Authorization")
			arr := strings.Split(header, " ")

			if len(arr) != 2 {
				handlers.MakeErrorResponse(w, r, log, "invalid token", http.StatusUnauthorized, requestId, nil)
				return
			}

			token := arr[1]
			if !isAuthorized(token) {
				handlers.MakeErrorResponse(w, r, log, "bad token", http.StatusUnauthorized, requestId, nil)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// JWTValidateModeratorMW Проверка на наличие токена модератора middleware ПЕРЕНЕСТИ
func JWTValidateModeratorMW(log *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestId := middleware.GetReqID(r.Context())
			header := r.Header.Get("Authorization")
			arr := strings.Split(header, " ")

			if len(arr) != 2 {
				handlers.MakeErrorResponse(w, r, log, "unauthorized", http.StatusUnauthorized, requestId, nil)
				return
			}

			token := arr[1]
			if !isAuthorized(token) {
				handlers.MakeErrorResponse(w, r, log, "bad token", http.StatusUnauthorized, requestId, nil)
				return
			}

			utype := GetUserType(token)
			if utype != moderatorType {
				handlers.MakeErrorResponse(w, r, log, "user is not moderator", http.StatusUnauthorized, requestId, nil)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

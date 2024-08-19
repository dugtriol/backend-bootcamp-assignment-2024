package auth

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/services"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

const (
	TokenExp  = time.Hour * 3
	SecretKey = "supersecretkey"
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

func BuildJWTString(typeUser string) (string, error) {
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

func IsAuthorized(tokenString string) bool {
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
			services.MakeErrorResponse(w, r, log, "request body is empty", http.StatusBadRequest, requestId, err)
			return
		}
		if err != nil {
			services.MakeErrorResponse(
				w,
				r,
				log,
				"failed to decode request body",
				http.StatusInternalServerError,
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
			services.MakeErrorResponse(w, r, log, "Invalid request", http.StatusBadRequest, requestId, err)
			return
		}

		jwtString, err := BuildJWTString(req.UserType)
		if err != nil {
			services.MakeErrorResponse(w, r, log, "invalid jwt parse", http.StatusInternalServerError, requestId, err)
			return
		}
		log.Info("make token")

		resp := tokenResponse{Token: jwtString}
		render.JSON(w, r, resp)
		log.Info("success create token")
	}
}

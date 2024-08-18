package auth

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/handlers"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/pkg/response"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/pkg/storage/structures"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type loginRequest struct {
	Id       string `json:"id" validate:"required,uuid"`
	Password string `json:"password" validate:"required"`
}

type getUser interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*structures.User, error)
}

func Login(ctx context.Context, log *slog.Logger, data getUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		requestId := middleware.GetReqID(r.Context())
		var err error
		const op = "handlers.auth.Login"
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
		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("Invalid request")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr, requestId))
			return
		}

		id, err := uuid.Parse(req.Id)
		if err != nil {
			log.Error("failed to decode id")
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
		user, err := data.GetUserById(ctx, id)
		if err != nil {
			handlers.MakeErrorResponse(w, r, log, "failed to find user by id", http.StatusBadRequest, requestId, err)
			return
		}

		if err = checkPassword(req.Password, user.Password); err != nil {
			handlers.MakeErrorResponse(w, r, log, "invalid password", http.StatusBadRequest, requestId, err)
			return
		}

		jwtString, err := buildJWTString(user.Type)
		if err != nil {
			handlers.MakeErrorResponse(w, r, log, "invalid jwt parse", http.StatusBadRequest, requestId, err)
			return
		}
		log.Info("make token")

		resp := tokenResponse{Token: jwtString}
		render.JSON(w, r, resp)
		log.Info("success create token")
	}
}

package auth

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/services"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/pkg/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type userRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	UserType string `json:"user_type" validate:"required,oneof=moderator client"`
}

type userResponse struct {
	Id string `json:"user_id"`
}

type userSaver interface {
	SaveUser(ctx context.Context, email, password, userType string) (uuid.UUID, error)
}

func Register(ctx context.Context, log *slog.Logger, saver userSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req userRequest
		var err error
		const op = "handlers.user.create"
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

		password, err := services.HashPassword(req.Password)
		if err != nil {
			services.MakeErrorResponse(w, r, log, "failed to hash password", http.StatusBadRequest, requestId, err)
			return
		}

		id, err := saver.SaveUser(ctx, req.Email, password, req.UserType)
		if err != nil {
			// email is already in db
			services.MakeErrorResponse(w, r, log, "failed to save user in DB", http.StatusBadRequest, requestId, err)
			return
		}

		// send user id
		resp := userResponse{Id: id.String()}
		render.JSON(w, r, &resp)
		log.Info("success create token")
	}
}

package flat

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/datasource/storage/structures"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/services"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/pkg/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const on_moderate = "on moderate"

type moderationRequest struct {
	Id     int    `json:"id" validate:"required,min=1"`
	Status string `json:"status" validate:"required,oneof='created' 'approved' 'declined' 'on moderation'"`
}

type updateModeration interface {
	UpdateStatus(ctx context.Context, id int, status string) error
	GetFlat(ctx context.Context, id int) (*structures.Flat, error)
}

func Moderate(ctx context.Context, log *slog.Logger, moderation updateModeration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req moderationRequest
		var err error
		const op = "handlers.flat.moderate"
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

		if err = validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("Invalid request")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr, requestId))
			return
		}

		notChangedFlat, err := moderation.GetFlat(ctx, req.Id)
		if err != nil {
			services.MakeErrorResponse(w, r, log, "failed to find flat", http.StatusBadRequest, requestId, err)
			return
		}

		if notChangedFlat.Status == on_moderate {
			services.MakeErrorResponse(
				w,
				r,
				log,
				"the flat on moderate be another moderator",
				http.StatusBadRequest,
				requestId,
				err,
			)
			return
		}

		err = moderation.UpdateStatus(ctx, req.Id, req.Status)
		if err != nil {
			services.MakeErrorResponse(w, r, log, "failed to update status", http.StatusBadRequest, requestId, err)
			return
		}

		flat, err := moderation.GetFlat(ctx, req.Id)
		if err != nil {
			services.MakeErrorResponse(w, r, log, "failed to find flat", http.StatusBadRequest, requestId, err)
			return
		}

		render.JSON(w, r, &flat)
	}
}

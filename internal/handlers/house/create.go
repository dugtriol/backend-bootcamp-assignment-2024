package house

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

type houseRequest struct {
	Address   string `json:"address" validate:"required"`
	Year      int    `json:"year" validate:"required,min=0"`
	Developer string `json:"developer" validate:"required"`
}

type houseSaver interface {
	SaveHouse(ctx context.Context, address, developer string, year int) (*structures.House, error)
	GetHouse(ctx context.Context, id int) (*structures.House, error)
}

func Create(ctx context.Context, log *slog.Logger, saver houseSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req houseRequest
		var err error
		const op = "handlers.house.create"
		requestId := middleware.GetReqID(r.Context())
		log.With(
			slog.String("op", op),
			slog.String("request_id", requestId),
		)

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

		house, err := saver.SaveHouse(ctx, req.Address, req.Developer, req.Year)
		if err != nil {
			services.MakeErrorResponse(w, r, log, "failed to save house to db", http.StatusBadRequest, requestId, err)
			return
		}

		render.JSON(w, r, &house)
	}
}

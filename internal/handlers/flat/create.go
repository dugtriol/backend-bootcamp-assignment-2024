package flat

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/handlers"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/pkg/response"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/pkg/storage/structures"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type flatRequest struct {
	HouseId int `json:"house_id" validate:"required"`
	Price   int `json:"price" validate:"required,min=0"`
	Rooms   int `json:"rooms" validate:"required,min=1"`
}

type houseSaver interface {
	SaveFlat(ctx context.Context, houseId, price, rooms int) (*structures.Flat, error)
	UpdateDate(ctx context.Context, time time.Time, id int) error
}

func Create(ctx context.Context, log *slog.Logger, saver houseSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req flatRequest
		var err error
		const op = "handlers.house.create"
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

		if err = validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("Invalid request")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr, requestId))

			return
		}

		flat, err := saver.SaveFlat(ctx, req.HouseId, req.Price, req.Rooms)
		if err != nil {
			handlers.MakeErrorResponse(w, r, log, "failed to save flat to db", http.StatusBadRequest, requestId, err)
			return
		}
		t := time.Now()
		err = saver.UpdateDate(ctx, t, flat.HouseId)
		if err != nil {
			handlers.MakeErrorResponse(
				w,
				r,
				log,
				"failed to update time for house in db",
				http.StatusInternalServerError,
				requestId,
				err,
			)
			return
		}
		render.JSON(w, r, &flat)
	}
}

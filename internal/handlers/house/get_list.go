package house

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/datasource/storage/structures"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/handlers/auth"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const moderatorType = "moderator"

type getList interface {
	GetHouse(ctx context.Context, id int) (*structures.House, error)
	GetListByClient(ctx context.Context, id int) (*[]structures.Flat, error)
	GetListByModerator(ctx context.Context, id int) (*[]structures.Flat, error)
}

type GetListResponse struct {
	Flats *[]structures.Flat `json:"flats"`
}

func GetList(ctx context.Context, log *slog.Logger, getListFlats getList) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.flat.moderate"
		requestId := middleware.GetReqID(r.Context())
		var err error
		log.With(
			slog.String("op", op),
			slog.String("request_id", requestId),
		)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			services.MakeErrorResponse(
				w,
				r,
				log,
				"failed to get id from url param",
				http.StatusBadRequest,
				requestId,
				err,
			)
			return
		}

		_, err = getListFlats.GetHouse(ctx, id)
		if err != nil {
			services.MakeErrorResponse(w, r, log, "failed to find house", http.StatusBadRequest, requestId, err)
			return
		}

		header := r.Header.Get("Authorization")
		token := strings.Split(header, " ")[1]

		var flats *[]structures.Flat
		utype := auth.GetUserType(token)
		if utype == moderatorType {
			list, e := getListFlats.GetListByModerator(ctx, id)
			err = e
			flats = list
		} else {
			list, e := getListFlats.GetListByClient(ctx, id)
			err = e
			flats = list
		}
		if err != nil {
			services.MakeErrorResponse(w, r, log, "failed to get flats", http.StatusBadRequest, requestId, err)
			return
		}
		fmt.Println("flats", flats)
		listResponse := GetListResponse{Flats: flats}
		render.JSON(w, r, &listResponse)
	}
}

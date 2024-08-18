package handlers

import (
	"log/slog"
	"net/http"

	"github.com/dugtriol/backend-bootcamp-assignment-2024/pkg/response"
	"github.com/go-chi/render"
)

func MakeErrorResponse(
	w http.ResponseWriter, r *http.Request, log *slog.Logger, str string, code int, requestId string, err error,
) {
	log.Error(str, err)
	w.WriteHeader(code)
	render.JSON(w, r, response.MakeResponse(str, requestId, code))
}

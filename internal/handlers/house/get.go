package house

//house/create
//func Get(ctx context.Context, log *slog.Logger, saver houseSaver) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var req houseRequest
//		const op = "handlers.house.create"
//		log.With(
//			slog.String("op", op),
//			slog.String("request_id", middleware.GetReqID(r.Context())),
//		)
//
//		// decode
//		err := render.DecodeJSON(r.Body, &req)
//		if errors.Is(err, io.EOF) {
//			log.Message("request body is empty", "err", err.Message())
//			w.WriteHeader(http.StatusBadRequest)
//			render.JSON(w, r, response.Message("empty request"))
//			return
//		}
//		if err != nil {
//			log.Message("failed to decode request body", "err", err.Message())
//			w.WriteHeader(http.StatusBadRequest)
//			render.JSON(w, r, response.Message("failed to decode request body"))
//			return
//		}
//		log.Info("request body decoded")
//
//		if err := validator.New().Struct(req); err != nil {
//			var validateErr validator.ValidationErrors
//			errors.As(err, &validateErr)
//
//			log.Message("Invalid request", "err", err.Message())
//			w.WriteHeader(http.StatusBadRequest)
//
//			render.JSON(w, r, response.ValidationError(validateErr))
//			return
//		}
//
//		house, err := saver.GetHouse(ctx, req.Address)
//
//		if err != nil {
//			log.Message("house not found", "err", err)
//			w.WriteHeader(http.StatusBadRequest)
//			render.JSON(w, r, response.Message("house not found"))
//			return
//		}
//		render.JSON(w, r, &house)
//	}
//}

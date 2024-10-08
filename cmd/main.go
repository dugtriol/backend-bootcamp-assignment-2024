package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/config"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/datasource/cache"
	strg "github.com/dugtriol/backend-bootcamp-assignment-2024/internal/datasource/storage"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/handlers/auth"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/handlers/flat"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/internal/handlers/house"
	"github.com/dugtriol/backend-bootcamp-assignment-2024/pkg/db"
	mwLogger "github.com/dugtriol/backend-bootcamp-assignment-2024/pkg/middleware"
	"github.com/go-chi/render"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.MustLoad()

	// logger
	log := setupLogger()
	log.Info("initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")
	log.Info("cfg ", cfg)

	// database
	database, err := db.NewDB(ctx)
	if err != nil {
		log.Error(err.Error())
	}
	defer database.GetPool(ctx).Close()

	//storage
	bigcache, err := bigcache.New(ctx, bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		log.Error(err.Error())
	}

	storage := cache.NewClient(log, bigcache, strg.New(database, log))
	_ = storage

	if err != nil {
		log.Error("failed to init storage ")
		os.Exit(1)
	}

	//router
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.URLFormat)
	router.Use(middleware.Recoverer)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Group(
		func(r chi.Router) {
			r.Get("/dummyLogin", auth.GetDummyLogin(log))
			r.Post("/register", auth.Register(ctx, log, storage))
			r.Post("/login", auth.Login(ctx, log, storage))
		},
	)

	router.Group(
		func(r chi.Router) {
			r.Use(mwLogger.JWTValidateMW(log))

			r.Post("/flat/create", flat.Create(ctx, log, storage))
			r.Get("/house/{id}", house.GetList(ctx, log, storage))

			r.Group(
				func(c chi.Router) {
					c.Use(mwLogger.JWTValidateModeratorMW(log))

					c.Post("/house/create", house.Create(ctx, log, storage))
					c.Post("/flat/update", flat.Moderate(ctx, log, storage))
					//c.Post("/house/get", house.Get(ctx, log, storage))
				},
			)
		},
	)

	// long time request change
	serv := &http.Server{
		Addr:         cfg.Address,
		Handler:      http.TimeoutHandler(router, 1*time.Second, "long time request"),
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("start server")
	if err = serv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		os.Exit(1)
	}
}

// env string
func setupLogger() *slog.Logger {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return log
}

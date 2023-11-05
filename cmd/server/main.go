package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/overgoy/url-shortener/internal/controller"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func Start(cfg *config.Configuration) {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	baseController := controller.NewBaseController(logger, cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", baseController.Route())

	logger.Info("Server started on " + cfg.ServerAddress)
	err = http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		logger.Error("Failed to start the server", zap.Error(err))
		os.Exit(1)
	}
}

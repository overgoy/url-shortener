package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/overgoy/url-shortener/internal/controller"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func Start(cfg *config.Configuration) {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	baseController := controller.NewBaseController(logger, cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", baseController.Route())

	logger.Info("Server started on " + cfg.ServerAddress)
	err := http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		fmt.Printf("server: %v", err)
		os.Exit(1)
	}
}

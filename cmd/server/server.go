package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/overgoy/url-shortener/internal/controller"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func Start(cfg *config.Configuration) {
	logger := log.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.InfoLevel)

	baseController := controller.NewBaseController(logger, cfg)
	r := chi.NewRouter()
	r.Mount("/", baseController.Route())

	logger.Info("Server started on " + cfg.ServerAddress)
	err := http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		fmt.Printf("server: %v", err)
		os.Exit(1)
	}
}

package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/overgoy/url-shortener/internal/controller"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func Start() {
	logger := log.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.InfoLevel)

	baseController := controller.NewBaseController(logger)
	r := chi.NewRouter()
	r.Mount("/", baseController.Route())

	logger.Info("Server started on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("server: %v", err)
	}
}

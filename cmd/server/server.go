package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/overgoy/url-shortener/internal/controller"
	"github.com/overgoy/url-shortener/internal/handler" // импортирование пакета handler
	"log"                                               // использование стандартного логгера
	"net/http"
	"os"
)

func Start(cfg *config.Configuration) {
	baseController := controller.NewBaseController(cfg)
	r := chi.NewRouter()
	r.Use(handler.RequestLogger) // использование RequestLogger middleware
	r.Mount("/", baseController.Route())

	log.Printf("Server started on %s", cfg.ServerAddress)
	err := http.ListenAndServe(cfg.ServerAddress, r)
	if err != nil {
		fmt.Printf("server: %v", err)
		os.Exit(1)
	}
}

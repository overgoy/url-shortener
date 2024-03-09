package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/overgoy/url-shortener/internal/gzip"
	"github.com/overgoy/url-shortener/internal/handler"
	"github.com/overgoy/url-shortener/internal/logging"
)

type BaseController struct {
	cfg        *config.Configuration
	logger     logging.Logger
	urlHandler *handler.Handler // Добавляем экземпляр обработчика URL
}

func NewBaseController(cfg *config.Configuration, logger logging.Logger, urlHandler *handler.Handler) *BaseController {
	return &BaseController{
		cfg:        cfg,
		logger:     logger,
		urlHandler: urlHandler,
	}
}

func (c *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Use(gzip.GzipMiddleware)
	r.Use(logging.NewStructuredLogger(c.logger))
	r.Post("/", c.urlHandler.HandlePost)
	r.Get("/{id:[a-zA-Z0-9]+}", c.urlHandler.HandleGet)
	r.Post("/api/shorten", c.urlHandler.ShortenEndpoint)
	return r
}

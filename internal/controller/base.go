package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/overgoy/url-shortener/internal/gzip"
	"github.com/overgoy/url-shortener/internal/handler"
	"github.com/overgoy/url-shortener/internal/logger"
	"go.uber.org/zap"
)

type BaseController struct {
	logger     *zap.Logger
	cfg        *config.Configuration
	urlHandler *handler.App // Добавляем экземпляр обработчика URL
}

func NewBaseController(logger *zap.Logger, cfg *config.Configuration) *BaseController {
	return &BaseController{
		logger:     logger,
		cfg:        cfg,
		urlHandler: handler.NewApp(cfg, logger), // Инициализируем обработчик с конфигурацией
	}
}

func (c *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Use(gzip.GzipMiddleware)
	r.Use(logger.NewStructuredLogger(c.logger)) // Добавляем logger для логирования
	r.Post("/", c.urlHandler.HandlePost)
	r.Get("/{id:[a-zA-Z0-9]+}", c.urlHandler.HandleGet)
	r.Post("/api/shorten", c.urlHandler.ShortenEndpoint)
	return r
}

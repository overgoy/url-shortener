package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/overgoy/url-shortener/internal/middleware"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/overgoy/url-shortener/internal/handler"
)

type BaseController struct {
	logger     *logrus.Logger
	cfg        *config.Configuration
	urlHandler *handler.App // Добавляем экземпляр обработчика URL
}

func NewBaseController(logger *logrus.Logger, cfg *config.Configuration) *BaseController {
	return &BaseController{
		logger:     logger,
		cfg:        cfg,
		urlHandler: handler.NewApp(cfg, logger), // Инициализируем обработчик с конфигурацией
	}
}

func (c *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.NewStructuredLogger(c.logger)) // Добавляем middleware для логирования
	r.Post("/", c.handleMain)
	r.Get("/{id:[a-zA-Z0-9]+}", c.handleName)
	return r
}

func (c *BaseController) handleMain(writer http.ResponseWriter, request *http.Request) {
	c.urlHandler.HandlePost(writer, request) // Обращаемся к обработчику напрямую
}

func (c *BaseController) handleName(writer http.ResponseWriter, request *http.Request) {
	c.urlHandler.HandleGet(writer, request) // Обращаемся к обработчику напрямую
}

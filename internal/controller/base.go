package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/overgoy/url-shortener/internal/config"
	"net/http"

	"github.com/overgoy/url-shortener/internal/handler"
	log "github.com/sirupsen/logrus"
)

type BaseController struct {
	logger     *log.Logger
	cfg        *config.Configuration
	urlHandler *handler.URLHandler // Добавляем экземпляр обработчика URL
}

func NewBaseController(logger *log.Logger, cfg *config.Configuration) *BaseController {
	return &BaseController{
		logger:     logger,
		cfg:        cfg,
		urlHandler: handler.NewURLHandler(cfg), // Инициализируем обработчик с конфигурацией
	}
}

func (c *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
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

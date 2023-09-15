package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/overgoy/url-shortener/internal/config"
	"net/http"

	"github.com/overgoy/url-shortener/internal/handlers"
	log "github.com/sirupsen/logrus"
)

type BaseController struct {
	logger *log.Logger
	cfg    *config.Configuration
}

func NewBaseController(logger *log.Logger, cfg *config.Configuration) *BaseController {
	return &BaseController{
		logger: logger,
		cfg:    cfg,
	}
}

func (c *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", c.handleMain)
	r.Get("/{id:[a-zA-Z0-9]+}", c.handleName)
	return r
}

func (c *BaseController) handleMain(writer http.ResponseWriter, request *http.Request) {
	handlers.HandlePost(writer, request, c.cfg)
}

func (c *BaseController) handleName(writer http.ResponseWriter, request *http.Request) {
	handlers.HandleGet(writer, request)
}

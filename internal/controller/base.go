package controller

import (
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/overgoy/url-shortener/internal/handlers"
	log "github.com/sirupsen/logrus"
)

type BaseController struct {
	logger *log.Logger
}

func NewBaseController(logger *log.Logger) *BaseController {
	return &BaseController{
		logger: logger,
	}
}

func (c *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", c.handleMain)
	r.Get("/{id:[a-zA-Z0-9]+}", c.handleName)
	return r
}

func (c *BaseController) handleMain(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodPost {
		handlers.HandlePost(writer, request)
		return
	}
	http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
}

func (c *BaseController) handleName(writer http.ResponseWriter, request *http.Request) {
	handlers.HandleGet(writer, request)
}

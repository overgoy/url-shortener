package controller

//
//import (
//	"github.com/go-chi/chi/v5"
//	"net/http"
//
//	log "github.com/sirupsen/logrus"
//)
//
//type BaseController struct {
//	logger log.Logger
//}
//
//func NewBaseController(logger log.Logger) *BaseController {
//	return &BaseController{
//		logger: logger,
//	}
//}
//
//func (c *BaseController) Route() *chi.Mux {
//	r := chi.NewRouter()
//	r.Get("/", c.handleMain)
//	r.Get("/{name}", c.handleName)
//	return r
//}
//
//func (c *BaseController) handleMain(writer http.ResponseWriter, reques *http.Request) {
//
//}
//func (c *BaseController) handleName(writer http.ResponseWriter, reques *http.Request) {
//
//}

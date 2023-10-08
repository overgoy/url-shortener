// server.go
package server

import (
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/overgoy/url-shortener/internal/handler"
	"net/http"
)

func Start(cfg *config.Configuration) {
	app := handler.NewApp(cfg)

	// Создание логгера
	logger := handler.Logger // предполагается, что функция Logger находится в пакете handler

	// Использование логгера как middleware
	http.Handle("/", logger(http.HandlerFunc(app.HandlePost)))

	http.ListenAndServe(":8080", nil)
}

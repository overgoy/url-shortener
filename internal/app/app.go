package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/overgoy/url-shortener/internal/controller"
	"github.com/overgoy/url-shortener/internal/handler"
	"github.com/overgoy/url-shortener/internal/logging"
	"github.com/overgoy/url-shortener/internal/storage"
	"github.com/overgoy/url-shortener/internal/storage/file"
	"github.com/overgoy/url-shortener/internal/storage/memory"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func Start() {
	// Инициализация конфигурации
	cfg := config.New()

	// Создание нового логгера
	myLogger, err := logging.NewMyLogger()
	if err != nil {
		println("Failed to create logging: %v\n", err)
		os.Exit(1)
	}

	// Создание нового приложения
	application := NewApp(cfg, myLogger)

	application.Logger.Info("Server started on " + cfg.ServerAddress)
	err = http.ListenAndServe(cfg.ServerAddress, application.GetRouter())
	if err != nil {
		application.Logger.Error("Failed to start the server", zap.Error(err))
		os.Exit(1)
	}
}

func NewApp(cfg *config.Configuration, logger logging.Logger) *App {
	var urlStorage storage.URLStorage

	if cfg.FileStoragePath != "" {
		// Use file storage
		urlStorage = file.NewFileURLStorage(cfg.FileStoragePath)
		if _, err := urlStorage.LoadURLsFromFile(); err != nil {
			logger.Error("Error loading URLs from file", zap.Error(err))
		}
	} else {
		// Use in-memory storage
		urlStorage = memory.NewInMemoryURLStorage()
	}

	baseController := controller.NewBaseController(cfg, logger, handler.NewHandler(logger, urlStorage))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", baseController.Route())

	return &App{URLStorage: urlStorage, Config: cfg, Logger: logger, Router: r}
}

type App struct {
	URLStorage storage.URLStorage
	Config     *config.Configuration
	Logger     logging.Logger
	Router     *chi.Mux
}

func (h *App) GetRouter() *chi.Mux {
	return h.Router
}

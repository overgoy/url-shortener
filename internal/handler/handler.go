package handler

import (
	"encoding/json"
	"fmt"
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/overgoy/url-shortener/internal/models"
	"github.com/overgoy/url-shortener/internal/util"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type App struct {
	URLStore map[string]string
	Mux      sync.Mutex
	Config   *config.Configuration
	Logger   *zap.Logger
}

func NewApp(cfg *config.Configuration, logger *zap.Logger) *App {
	return &App{
		URLStore: make(map[string]string),
		Config:   cfg,
		Logger:   logger,
	}
}

func (h *App) ShortenEndpoint(w http.ResponseWriter, r *http.Request) {
	// Парсим JSON-запрос в структуру ShortenRequest
	var request models.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidURL(request.URL) {
		h.Logger.Error("Invalid URL format received")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL format"))
		return
	}

	shortURL, err := h.generateShortURL(request.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем ответ в формате JSON
	response := models.ShortenResponse{Result: shortURL}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *App) HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	longURL, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || len(strings.TrimSpace(string(longURL))) == 0 {
		h.Logger.Error("Error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading request"))
		return
	}

	if !isValidURL(string(longURL)) {
		h.Logger.Error("Invalid URL format received")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL format"))
		return
	}

	shortURL, err := h.generateShortURL(string(longURL))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *App) HandleGet(w http.ResponseWriter, r *http.Request) {
	logger := h.Logger

	if len(r.URL.Path) < 2 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		logger.Error("Invalid request")
		return
	}

	id := r.URL.Path[1:]

	h.Mux.Lock()
	longURL, ok := h.URLStore[id]
	h.Mux.Unlock()

	if !ok {
		h.Logger.Error("URL not found", zap.String("id", id))
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *App) generateShortURL(longURL string) (string, error) {
	id := util.GenerateID()
	h.Mux.Lock()
	defer h.Mux.Unlock()
	h.URLStore[id] = longURL

	baseURL := h.Config.BaseURL
	baseURL = strings.TrimSuffix(baseURL, "/")

	shortURL := fmt.Sprintf("%s/%s", baseURL, id)
	return shortURL, nil
}

func isValidURL(u string) bool {
	parsedURL, err := url.Parse(u)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}

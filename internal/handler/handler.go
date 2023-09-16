package handler

import (
	"fmt"
	"github.com/overgoy/url-shortener/internal/config"
	logger "github.com/overgoy/url-shortener/internal/logging"
	"github.com/overgoy/url-shortener/internal/util"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type URLHandler struct {
	Store  map[string]string
	Mux    sync.Mutex
	Config *config.Configuration
	Logger logger.Logger
}

func NewURLHandler(cfg *config.Configuration, logger logger.Logger) *URLHandler {
	return &URLHandler{
		Store:  make(map[string]string),
		Config: cfg,
		Logger: logger,
	}
}

func (h *URLHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
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

	id := util.GenerateID()
	h.Mux.Lock()
	h.Store[id] = string(longURL)
	h.Mux.Unlock()

	baseURL := h.Config.BaseURL
	baseURL = strings.TrimSuffix(baseURL, "/")

	shortURL := fmt.Sprintf("%s/%s", baseURL, id)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *URLHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[1:]

	h.Mux.Lock()
	longURL, ok := h.Store[id]
	h.Mux.Unlock()

	if !ok {
		h.Logger.WithField("id", id).Error("URL not found")
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func isValidURL(u string) bool {
	parsedURL, err := url.Parse(u)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}

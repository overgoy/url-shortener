package handler

import (
	"fmt"
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/overgoy/url-shortener/internal/util"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		// Создаем запись для логирования ответа
		ww := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(ww, r)

		// Запись логов с использованием log
		log.Printf("%s %s %v", r.Method, r.URL, time.Since(startTime))
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

type App struct {
	URLStore map[string]string
	Mux      sync.Mutex
	Config   *config.Configuration
}

func NewApp(cfg *config.Configuration) *App {
	return &App{
		URLStore: make(map[string]string),
		Config:   cfg,
	}
}

func (h *App) HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	longURL, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || len(strings.TrimSpace(string(longURL))) == 0 {
		log.Println("Error reading request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading request"))
		return
	}

	if !isValidURL(string(longURL)) {
		log.Println("Invalid URL format received")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL format"))
		return
	}

	id := util.GenerateID()
	h.Mux.Lock()
	h.URLStore[id] = string(longURL)
	h.Mux.Unlock()

	baseURL := h.Config.BaseURL
	baseURL = strings.TrimSuffix(baseURL, "/")

	shortURL := fmt.Sprintf("%s/%s", baseURL, id)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *App) HandleGet(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) < 2 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id := r.URL.Path[1:]

	h.Mux.Lock()
	longURL, ok := h.URLStore[id]
	h.Mux.Unlock()

	if !ok {
		log.Printf("URL not found for id: %s\n", id)
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

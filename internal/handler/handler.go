package handler

import (
	"errors"
	"fmt"
	"github.com/overgoy/url-shortener/internal/config"
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
}

func NewURLHandler(cfg *config.Configuration) *URLHandler {
	return &URLHandler{
		Store:  make(map[string]string),
		Config: cfg,
	}
}

func (h *URLHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	longURL, err := readRequestBody(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidURL(longURL) {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	id := util.GenerateID()
	h.Mux.Lock()
	h.Store[id] = longURL
	h.Mux.Unlock()

	baseURL := h.Config.BaseURL
	if strings.HasSuffix(baseURL, "/") {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}

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
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func readRequestBody(body io.ReadCloser) (string, error) {
	content, err := io.ReadAll(body)
	if err != nil {
		return "", errors.New("Error reading request")
	}

	contentStr := strings.TrimSpace(string(content))
	if len(contentStr) == 0 {
		return "", errors.New("Request body is empty")
	}

	return contentStr, nil
}

func isValidURL(u string) bool {
	parsedURL, err := url.Parse(u)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}

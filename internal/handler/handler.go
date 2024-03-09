package handler

import (
	"encoding/json"
	"github.com/overgoy/url-shortener/internal/logging"
	"github.com/overgoy/url-shortener/internal/models"
	"github.com/overgoy/url-shortener/internal/storage"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Handler представляет собой HTTP-обработчик для сокращения URL
type Handler struct {
	Logger     logging.Logger
	URLStorage storage.URLStorage
}

// NewHandler создает новый экземпляр Handler с переданным логгером и хранилищем URL
func NewHandler(logger logging.Logger, urlStorage storage.URLStorage) *Handler {
	return &Handler{
		Logger:     logger,
		URLStorage: urlStorage,
	}
}

func (h *Handler) ShortenEndpoint(w http.ResponseWriter, r *http.Request) {
	// Парсим JSON-запрос в структуру ShortenRequest
	var request models.ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.Logger.Error("Error decoding JSON request", zap.Error(err))
		return
	}

	if !isValidURL(request.URL) {
		h.Logger.Info(`Invalid URL format received`, zap.String("url", request.URL))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL format"))
		return
	}

	shortURL, err := h.URLStorage.GenerateShortURL(request.URL)
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

	h.URLStorage.SaveURLToDisk(request.URL, shortURL)

	// Получаем объект http.Header из запроса
	header := w.Header()
	// Устанавливаем значение заголовка Content-Type
	header.Set("Content-Type", "application/json")
	// Устанавливаем статус ответа
	w.WriteHeader(http.StatusCreated)
	// Пишем ответ в теле
	w.Write(jsonResponse)
}

func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
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

	shortURL, err := h.URLStorage.GenerateShortURL(string(longURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.URLStorage.SaveURLToDisk(string(longURL), shortURL)

	// Получаем объект http.Header из запроса
	header := w.Header()
	// Устанавливаем значение заголовка Content-Type
	header.Set("Content-Type", "text/plain")
	// Устанавливаем статус ответа
	w.WriteHeader(http.StatusCreated)
	// Пишем ответ в теле
	w.Write([]byte(shortURL))
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {

	if len(r.URL.Path) < 2 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		h.Logger.Error("Invalid request")
		return
	}

	id := r.URL.Path[1:]

	// Загрузка URL из файла
	data, err := h.URLStorage.LoadURLsFromFile()
	if err != nil {
		h.Logger.Error("Error loading URLs from file", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	longURL, ok := data[id]
	if !ok {
		h.Logger.Error("URL not found", zap.String("id", id))
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Получаем объект http.Header из запроса
	header := w.Header()
	header.Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func isValidURL(u string) bool {
	parsedURL, err := url.ParseRequestURI(u)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}

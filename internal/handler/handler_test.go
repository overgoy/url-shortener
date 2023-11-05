package handler

import (
	"bytes"
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePost(t *testing.T) {
	cfg := &config.Configuration{
		ServerAddress: "localhost:8888",
		BaseURL:       "http://localhost:8000/",
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	h := NewApp(cfg, logger)

	type want struct {
		code         int
		contentType  string
		bodyContains string
	}

	tests := []struct {
		name     string
		inputURL string
		want     want
	}{
		{"valid URL", "https://practicum.yandex.ru/", want{http.StatusCreated, "text/plain", cfg.BaseURL}},
		{"empty URL", "", want{http.StatusBadRequest, "", ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/", bytes.NewBufferString(tt.inputURL))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "text/plain")

			rec := httptest.NewRecorder()

			h.HandlePost(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestHandleGet(t *testing.T) {
	cfg := &config.Configuration{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080/",
	}
	// Создаем экземпляр logrus.Logger для использования в тестах
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	h := NewApp(cfg, logger)

	// Добавляем тестовую ссылку в хранилище
	testID := "testID"
	h.URLStore[testID] = "https://practicum.yandex.ru/"

	type want struct {
		code     int
		location string
	}

	tests := []struct {
		name    string
		inputID string
		want    want
	}{
		{"valid ID", testID, want{http.StatusTemporaryRedirect, "https://practicum.yandex.ru/"}},
		{"invalid ID", "invalidID", want{http.StatusNotFound, ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/"+tt.inputID, nil)
			require.NoError(t, err)

			rec := httptest.NewRecorder()

			h.HandleGet(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.location, res.Header.Get("Location"))
		})
	}
}
func TestShortenEndpoint(t *testing.T) {
	cfg := &config.Configuration{
		ServerAddress: "localhost:8888",
		BaseURL:       "http://localhost:8000/"}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	h := NewApp(cfg, logger)

	type want struct {
		code         int
		contentType  string
		bodyContains string
	}

	tests := []struct {
		name     string
		inputURL string
		want     want
	}{
		{"valid URL", "https://practicum.yandex.ru/", want{http.StatusCreated, "application/json", "http://localhost:8000/"}},
		{"empty URL", "", want{http.StatusBadRequest, "", ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := `{"url":"` + tt.inputURL + `"}`
			req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBufferString(requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			h.ShortenEndpoint(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			if tt.want.bodyContains != "" {
				responseBody, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Contains(t, string(responseBody), tt.want.bodyContains)
			}
		})
	}
}

func TestGzipMiddleware(t *testing.T) {
	cfg := &config.Configuration{
		ServerAddress: "localhost:8888",
		BaseURL:       "http://localhost:8000/",
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	h := NewApp(cfg, logger)

	type want struct {
		code         int
		contentType  string
		bodyContains string
	}

	tests := []struct {
		name     string
		inputURL string
		want     want
	}{
		{"valid URL", "https://practicum.yandex.ru/", want{http.StatusCreated, "text/plain", cfg.BaseURL}},
		{"empty URL", "", want{http.StatusBadRequest, "", ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/", bytes.NewBufferString(tt.inputURL))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "text/plain")
			req.Header.Set("Accept-Encoding", "gzip") // Устанавливаем заголовок Accept-Encoding

			rec := httptest.NewRecorder()

			h.HandlePost(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

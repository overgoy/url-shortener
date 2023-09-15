package handlers

import (
	"bytes"
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePost(t *testing.T) {
	cfg := &config.Configuration{
		ServerAddress: "localhost:8888",
		BaseURL:       "http://localhost:8000/",
	}

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
		{"valid URL", "https://practicum.yandex.ru/", want{http.StatusCreated, "text/plain", cfg.BaseURL}}, // Теперь мы проверяем, что короткий URL начинается с BaseURL
		{"empty URL", "", want{http.StatusBadRequest, "", ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/", bytes.NewBufferString(tt.inputURL))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "text/plain")

			rec := httptest.NewRecorder()

			HandlePost(rec, req, cfg)

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
	// Добавляем тестовую ссылку в хранилище
	testID := "testID"
	urlStore[testID] = "https://practicum.yandex.ru/"

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

			HandleGet(rec, req, cfg)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.location, res.Header.Get("Location"))
		})
	}
}

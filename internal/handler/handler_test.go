package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/overgoy/url-shortener/internal/storage/memory"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/overgoy/url-shortener/internal/handler"
	"github.com/overgoy/url-shortener/internal/models"
	"github.com/overgoy/url-shortener/internal/storage"
)

func TestShortenEndpoint(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	urlStorage := memory.NewInMemoryURLStorage()
	h := handler.NewHandler(logger, urlStorage)

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

				var response models.ShortenResponse
				err = json.Unmarshal(responseBody, &response)
				require.NoError(t, err)

				assert.Contains(t, response.Result, tt.want.bodyContains)
			}
		})
	}
}

func TestHandlePost(t *testing.T) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	urlStorage := storage.NewMemoryStorage()
	h := handler.NewHandler(logger, urlStorage)

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
		{"valid URL", "https://practicum.yandex.ru/", want{http.StatusCreated, "text/plain", "http://localhost:8000/"}},
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
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	urlStorage := storage.NewMemoryStorage()
	h := handler.NewHandler(logger, urlStorage)

	testID := "testID"
	testURL := "https://practicum.yandex.ru/"
	urlStorage.SaveURLToDisk(testURL, testID)

	type want struct {
		code     int
		location string
	}

	tests := []struct {
		name    string
		inputID string
		want    want
	}{
		{"valid ID", testID, want{http.StatusTemporaryRedirect, testURL}},
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

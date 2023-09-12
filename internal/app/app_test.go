package app

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePost(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name     string
		inputURL string
		want     want
	}{
		{"valid URL", "https://practicum.yandex.ru/", want{http.StatusCreated, "text/plain"}},
		{"empty URL", "", want{http.StatusBadRequest, ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/", bytes.NewBufferString(tt.inputURL))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "text/plain")

			rec := httptest.NewRecorder()

			HandleRequest(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestHandleGet(t *testing.T) {
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

			HandleRequest(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.location, res.Header.Get("Location"))
		})
	}
}

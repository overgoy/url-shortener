package handlers

import (
	"fmt"
	"github.com/overgoy/url-shortener/internal/config"
	"github.com/overgoy/url-shortener/internal/util"
	"io"
	"net/http"
	"strings"
)

var urlStore = make(map[string]string)

func HandlePost(w http.ResponseWriter, r *http.Request, cfg *config.Configuration) {
	if r.URL.Path != "/" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	longURL, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil || len(strings.TrimSpace(string(longURL))) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading request"))
		return
	}

	id := util.GenerateID()
	urlStore[id] = string(longURL)

	shortURL := fmt.Sprintf("%s%s", cfg.BaseURL, id)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func HandleGet(w http.ResponseWriter, r *http.Request, cfg *config.Configuration) {
	id := r.URL.Path[1:]

	// Проверьте наличие длинного URL по идентификатору в хранилище
	longURL, ok := urlStore[id]
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

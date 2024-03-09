package memory

import (
	"github.com/overgoy/url-shortener/internal/util"
	"sync"
)

type InMemoryURLStorage struct {
	URLStore map[string]string
	Mux      sync.Mutex
}

func NewInMemoryURLStorage() *InMemoryURLStorage {
	return &InMemoryURLStorage{
		URLStore: make(map[string]string),
	}
}

func (s *InMemoryURLStorage) LoadURLsFromFile() (map[string]string, error) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	return s.URLStore, nil
}

func (s *InMemoryURLStorage) SaveURLToDisk(originalURL string, shortURL string) error {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	s.URLStore[shortURL] = originalURL

	return nil
}

func (s *InMemoryURLStorage) GenerateShortURL(longURL string) (string, error) {
	id := util.GenerateID()

	s.Mux.Lock()
	defer s.Mux.Unlock()
	s.URLStore[id] = longURL

	return id, nil
}

package file

import (
	"github.com/overgoy/url-shortener/internal/util"
	"io"
	"sync"
)

type FileURLStorage struct {
	FilePath string
	Mux      sync.Mutex
}

func NewFileURLStorage(filePath string) *FileURLStorage {
	return &FileURLStorage{
		FilePath: filePath,
	}
}

func (s *FileURLStorage) LoadURLsFromFile() (map[string]string, error) {

	if s.FilePath == "" {
		return nil, nil
	}

	consumer, err := util.NewConsumer(s.FilePath)
	if err != nil {
		return nil, err
	}
	defer consumer.Close()

	s.Mux.Lock()
	defer s.Mux.Unlock()

	urlStore := make(map[string]string)

	for {
		var urlData util.URLData
		err := consumer.Read(&urlData)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		urlStore[urlData.ShortURL] = urlData.OriginalURL
	}

	return urlStore, nil
}

func (s *FileURLStorage) SaveURLToDisk(originalURL string, shortURL string) error {
	if s.FilePath == "" {
		return nil
	}

	producer, err := util.NewProducer(s.FilePath)
	if err != nil {
		return err
	}
	defer producer.Close()

	s.Mux.Lock()
	defer s.Mux.Unlock()

	id := util.GenerateID()

	urlData := util.URLData{
		UUID:        id,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}

	return producer.Write(&urlData)
}

func (s *FileURLStorage) GenerateShortURL(longURL string) (string, error) {
	id := util.GenerateID()

	urlStore, err := s.LoadURLsFromFile()
	if err != nil {
		return "", err
	}

	urlStore[id] = longURL

	err = s.SaveURLToDisk(longURL, id)
	if err != nil {
		return "", err
	}

	return id, nil
}

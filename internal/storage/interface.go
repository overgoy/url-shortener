package storage

type URLStorage interface {
	LoadURLsFromFile() (map[string]string, error)
	SaveURLToDisk(shortURL string, originalURL string) error
	GenerateShortURL(longURL string) (string, error)
}

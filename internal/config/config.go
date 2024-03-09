package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type Configuration struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func New() *Configuration {
	config := &Configuration{}

	// Значения по умолчанию
	defaultAddr := "localhost:8080"
	defaultBaseURL := "http://localhost:8080/"
	defaultFileStoragePath := "/tmp/short-url-db.json"

	// Регистрация флагов
	addrFlag := flag.String("a", defaultAddr, "Address for the HTTP server to run on")
	baseURLFlag := flag.String("b", defaultBaseURL, "Base URL for the resulting shortened URL")
	fileStoragePathFlag := flag.String("f", defaultFileStoragePath, "Path to the file where URL data will be stored")

	// Анализировать переданные аргументы командной строки
	flag.Parse()

	if err := env.Parse(config); err != nil {
		log.WithError(err).Error("Ошибка при разборе переменных окружения")
	}

	// Если переменные окружения не установлены, используем значения из флагов
	if config.ServerAddress == "" {
		config.ServerAddress = *addrFlag
	}
	if config.BaseURL == "" {
		config.BaseURL = *baseURLFlag
	}
	if config.FileStoragePath == "" {
		config.FileStoragePath = *fileStoragePathFlag
	}

	log.WithFields(log.Fields{
		"ServerAddress": config.ServerAddress,
		"BaseURL":       config.BaseURL,
	}).Info("Используемые настройки конфигурации")

	return config
}

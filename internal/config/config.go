package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	logger "github.com/overgoy/url-shortener/internal/logging"
	log "github.com/sirupsen/logrus"
)

type Configuration struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func New(logger logger.Logger) *Configuration {
	config := &Configuration{}

	// Значения по умолчанию
	defaultAddr := "localhost:8080"
	defaultBaseURL := "http://localhost:8080/"

	// Регистрация флагов
	addrFlag := flag.String("a", defaultAddr, "Address for the HTTP server to run on")
	baseURLFlag := flag.String("b", defaultBaseURL, "Base URL for the resulting shortened URL")
	//// Регистрация флагов
	//flag.StringVar(&config.ServerAddress, "a", defaultAddr, "Address for the HTTP server to run on")
	//flag.StringVar(&config.BaseURL, "b", defaultBaseURL, "Base URL for the resulting shortened URL")

	// Анализировать переданные аргументы командной строки
	flag.Parse()

	if err := env.Parse(config); err != nil {
		logger.WithError(err).Error("Ошибка при разборе переменных окружения")
	}

	// Если переменные окружения не установлены, используем значения из флагов
	if config.ServerAddress == "" {
		config.ServerAddress = *addrFlag
	}
	if config.BaseURL == "" {
		config.BaseURL = *baseURLFlag
	}

	log.WithFields(log.Fields{
		"ServerAddress": config.ServerAddress,
		"BaseURL":       config.BaseURL,
	}).Info("Используемые настройки конфигурации")

	return config
}

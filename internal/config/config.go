package config

import (
	"flag"
)

type Configuration struct {
	ServerAddress string
	BaseURL       string
}

func New() *Configuration {
	config := &Configuration{}

	// Значения по умолчанию
	defaultAddr := "localhost:8888"
	defaultBaseURL := "http://localhost:8000/"

	// Регистрация флагов
	flag.StringVar(&config.ServerAddress, "a", defaultAddr, "Address for the HTTP server to run on")
	flag.StringVar(&config.BaseURL, "b", defaultBaseURL, "Base URL for the resulting shortened URL")

	// Анализировать переданные аргументы командной строки
	flag.Parse()

	return config
}

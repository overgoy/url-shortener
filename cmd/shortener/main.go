package main

import (
	"github.com/overgoy/url-shortener/cmd/server"
	"github.com/overgoy/url-shortener/internal/config"
)

func main() {
	cfg := config.New()
	server.Start(cfg)
}

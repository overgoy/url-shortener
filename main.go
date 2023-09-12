package main

import (
	"github.com/overgoy/url-shortener/cmd/client"
	"github.com/overgoy/url-shortener/cmd/server"
	"time"
)

func main() {
	go server.Start()
	time.Sleep(1 * time.Second)

	client.RequestURL()
}

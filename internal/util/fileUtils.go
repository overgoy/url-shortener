package util

import (
	"encoding/json"
	"os"
)

type URLData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Producer struct {
	file *os.File
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Producer{file: file}, nil
}

func (p *Producer) Write(data *URLData) error {
	encoder := json.NewEncoder(p.file)
	return encoder.Encode(data)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file *os.File
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &Consumer{file: file}, nil
}

func (c *Consumer) Read(data *URLData) error {
	decoder := json.NewDecoder(c.file)
	return decoder.Decode(data)
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

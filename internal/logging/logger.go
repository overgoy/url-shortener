package logger

import (
	log "github.com/sirupsen/logrus"
	"io"
)

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
	SetOutput(output io.Writer)
	SetLevel(level log.Level)
}

type LogrusAdapter struct {
	logger *log.Logger
}

func NewLogrusAdapter() Logger {
	return &LogrusAdapter{
		logger: log.StandardLogger(),
	}
}

func (l *LogrusAdapter) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *LogrusAdapter) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *LogrusAdapter) WithError(err error) Logger {
	return &LogrusAdapter{
		logger: l.logger.WithError(err).Logger,
	}
}

func (l *LogrusAdapter) WithField(key string, value interface{}) Logger {
	return &LogrusAdapter{
		logger: l.logger.WithField(key, value).Logger,
	}
}

func (l *LogrusAdapter) SetOutput(output io.Writer) {
	l.logger.SetOutput(output)
}

func (l *LogrusAdapter) SetLevel(level log.Level) {
	l.logger.SetLevel(level)
}

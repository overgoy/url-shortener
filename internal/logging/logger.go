package logger

import log "github.com/sirupsen/logrus"

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
}

type LogrusAdapter struct {
	*log.Entry
}

func NewLogrusAdapter() Logger {
	return &LogrusAdapter{
		Entry: log.NewEntry(log.StandardLogger()),
	}
}

func (l *LogrusAdapter) Info(args ...interface{}) {
	l.Entry.Info(args...)
}

func (l *LogrusAdapter) Error(args ...interface{}) {
	l.Entry.Error(args...)
}

func (l *LogrusAdapter) WithError(err error) Logger {
	return &LogrusAdapter{
		Entry: l.Entry.WithError(err),
	}
}

func (l *LogrusAdapter) WithField(key string, value interface{}) Logger {
	return &LogrusAdapter{
		Entry: l.Entry.WithField(key, value),
	}
}

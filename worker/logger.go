package worker

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) print(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).
        Str("proccess", "worker").
        Msg(fmt.Sprint(args...))
}

func (l *Logger) Debug(args ...interface{}) {
	l.print(zerolog.DebugLevel, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.print(zerolog.InfoLevel, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.print(zerolog.WarnLevel, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.print(zerolog.ErrorLevel, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.print(zerolog.FatalLevel, args...)
}

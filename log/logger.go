package log

import (
	"fmt"
	"os"
	"time"
)

const (
	defaultLevel      = Info
	defaultDateFormat = "2006-01-02 15:04:05"

	stderrLevel = Error
)

type Logger struct {
	Name       string
	Level      Level
	DateFormat string
}

func New(name string) *Logger {
	return &Logger{
		Name:       name,
		Level:      defaultLevel,
		DateFormat: defaultDateFormat,
	}
}

func (logger *Logger) Log(level Level, message string, arguments ...interface{}) {
	if level <= logger.Level {
		text := fmt.Sprintf("%s %s %s: %s\n",
			time.Now().Format(logger.DateFormat),
			level.String(),
			logger.Name,
			fmt.Sprintf(message, arguments...),
		)

		if level <= stderrLevel {
			_, _ = os.Stderr.WriteString(text)
		} else {
			_, _ = os.Stdout.WriteString(text)
		}
	}
}

func (logger *Logger) Fatal(message string, arguments ...interface{}) {
	logger.Log(Fatal, message, arguments...)
}

func (logger *Logger) Error(message string, arguments ...interface{}) {
	logger.Log(Error, message, arguments...)
}

func (logger *Logger) Warn(message string, arguments ...interface{}) {
	logger.Log(Warn, message, arguments...)
}

func (logger *Logger) Info(message string, arguments ...interface{}) {
	logger.Log(Info, message, arguments...)
}

func (logger *Logger) Debug(message string, arguments ...interface{}) {
	logger.Log(Debug, message, arguments...)
}

package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
)

type Logger struct {
	level   zerolog.Level
	verbose bool
	logger  zerolog.Logger
}

func NewLogger(verbose bool) *Logger {
	level := zerolog.InfoLevel
	if verbose {
		level = zerolog.DebugLevel
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
	logger := zerolog.New(output).With().Timestamp().Logger().Level(level)

	return &Logger{
		level:   level,
		verbose: verbose,
		logger:  logger,
	}
}

func (l *Logger) Info(msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	color.Green("[INFO] %s - %s", timestamp, msg)
}

func (l *Logger) Error(msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	color.Red("[ERROR] %s - %s", timestamp, msg)
}

func (l *Logger) Debug(msg string) {
	if l.verbose {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		color.Blue("[DEBUG] %s - %s", timestamp, msg)
	}
}

func (l *Logger) Warn(msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	color.Yellow("[WARN] %s - %s", timestamp, msg)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.verbose {
		l.Debug(fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

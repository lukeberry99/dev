package ui

import (
	"fmt"

	"github.com/fatih/color"
)

type Logger struct {
	verbose bool
}

func NewLogger(verbose bool) *Logger {
	return &Logger{
		verbose: verbose,
	}
}

func (l *Logger) Info(msg string) {
	fmt.Printf("ℹ️  %s\n", msg)
}

func (l *Logger) Error(msg string) {
	color.Red("❌ %s", msg)
}

func (l *Logger) Debug(msg string) {
	if l.verbose {
		color.Cyan("🔍 %s", msg)
	}
}

func (l *Logger) Warn(msg string) {
	color.Yellow("⚠️  %s", msg)
}

func (l *Logger) Success(msg string) {
	color.Green("✅ %s", msg)
}

func (l *Logger) Progress(msg string) {
	color.Blue("🔄 %s", msg)
}

func (l *Logger) Section(msg string) {
	fmt.Println()
	color.Cyan("🔧 %s", msg)
}

func (l *Logger) Step(msg string) {
	fmt.Printf("   • %s\n", msg)
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

package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

const (
	INFO    = "INFO"
	WARNING = "WARNING"
	ERROR   = "ERROR"
	DEBUG   = "DEBUG"
)

type ColoredLogger struct {
	*log.Logger
}

// NewColoredLogger creates a new logger. If out is nil, it defaults to os.Stdout.
func NewColoredLogger(prefix string, out io.Writer) *ColoredLogger {
	if out == nil {
		out = os.Stdout
	}
	return &ColoredLogger{
		Logger: log.New(out, prefix, log.Ltime),
	}
}

func (l *ColoredLogger) Info(format string, v ...interface{}) {
	l.Printf("%s[%s]%s %s", colorGreen, INFO, colorReset, fmt.Sprintf(format, v...))
}

func (l *ColoredLogger) Warning(format string, v ...interface{}) {
	l.Printf("%s[%s]%s %s", colorYellow, WARNING, colorReset, fmt.Sprintf(format, v...))
}

func (l *ColoredLogger) Error(format string, v ...interface{}) {
	l.Printf("%s[%s]%s %s", colorRed, ERROR, colorReset, fmt.Sprintf(format, v...))
}

func (l *ColoredLogger) Debug(format string, v ...interface{}) {
	l.Printf("%s[%s]%s %s", colorBlue, DEBUG, colorReset, fmt.Sprintf(format, v...))
}

var Default = NewColoredLogger("", os.Stdout)

package govenmo

import (
	"log"
	"os"
)

var logger *log.Logger

type nullWriter struct{}

func (nw nullWriter) Write(p []byte) (n int, err error) {
	return
}

func init() {
	logger = log.New(nullWriter{}, "", 0)
}

// Call EnableLogging to start logging. If you pass nil, a default logger will be used.
// Or, you can pass a Logger instance.
func EnableLogging(newLogger *log.Logger) {
	if newLogger != nil {
		logger = newLogger
	} else {
		logger = log.New(os.Stderr, "govenmo", log.LstdFlags)
	}
}

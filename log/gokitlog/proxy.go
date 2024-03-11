package gokitlog

import (
	"fmt"

	"github.com/go-kit/log"
)

type LoggerProxy struct {
	msgField string
	logger   log.Logger
}

// NewLoggerProxy returns a new wrapper for the given logger.
// Received data is rendered as string and passed to the given
// logger instance as structured data using the provided field.
func NewLoggerProxy(l log.Logger, f string) *LoggerProxy {
	result := &LoggerProxy{
		msgField: f,
		logger:   l,
	}

	return result
}

func (p *LoggerProxy) Println(v ...interface{}) {
	p.logger.Log(p.msgField, fmt.Sprintln(v...))
}

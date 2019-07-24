package httputil

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// LogWriter implements io.Writer and writes all incoming text out to the specified log level.
type LogWriter struct {
	Level     logrus.Level
	Component string
}

func (d LogWriter) Write(p []byte) (n int, err error) {
	logrus.WithField("component", d.Component).Log(d.Level, strings.TrimRight(string(p), "\n"))

	return len(p), nil
}

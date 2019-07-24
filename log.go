package httputil


import (
	"strings"

	"github.com/sirupsen/logrus"
)


// DebugLogWriter implements io.Writer and writes all incoming text out to log level info.
type DebugLogWriter struct {
	Component string
}

func (d DebugLogWriter) Write(p []byte) (n int, err error) {
	logrus.WithField("component", d.Component).Debug(strings.TrimRight(string(p), "\n"))

	return len(p), nil
}

/*
Copyright 2019 Christian Banse

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package httputil

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// LogWriter implements io.Writer and writes all incoming text out to the specified log level.
type LogWriter struct {
	Logger    *logrus.Logger
	Level     logrus.Level
	Component string
}

func (d LogWriter) Write(p []byte) (n int, err error) {
	var entry *logrus.Entry

	if d.Logger == nil {
		entry = logrus.WithField("component", d.Component)
	} else {
		entry = d.Logger.WithField("component", d.Component)
	}

	entry.Log(d.Level, strings.TrimRight(string(p), "\n"))

	return len(p), nil
}

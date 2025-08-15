// Copyright (c) 2024. Licensed under the MIT License.
package httputil

import (
	"net/http"
	"strings"
	"sync"
)

// LogBuffer stores recent log lines in memory.
type LogBuffer struct {
	mu    sync.Mutex
	lines []string
	max   int
}

// NewLogBuffer creates a new buffer that keeps up to max lines.
func NewLogBuffer(max int) *LogBuffer {
	return &LogBuffer{max: max}
}

// Write implements io.Writer and appends lines to the buffer.
func (lb *LogBuffer) Write(p []byte) (int, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	for _, line := range strings.Split(string(p), "\n") {
		if line == "" {
			continue
		}
		if len(lb.lines) >= lb.max {
			lb.lines = lb.lines[1:]
		}
		lb.lines = append(lb.lines, line)
	}
	return len(p), nil
}

// Lines returns a copy of stored log lines.
func (lb *LogBuffer) Lines() []string {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	out := make([]string, len(lb.lines))
	copy(out, lb.lines)
	return out
}

// GlobalLogBuffer collects logs for the UI.
var GlobalLogBuffer = NewLogBuffer(200)

// LogsHandler writes collected logs as plain text.
func LogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(strings.Join(GlobalLogBuffer.Lines(), "\n")))
}

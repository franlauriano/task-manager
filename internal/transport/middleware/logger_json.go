package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

// JSONLogFormatter implements chimw.LogFormatter and emits one JSON object per
// request on a single line (NDJSON).
type JSONLogFormatter struct {
	Out io.Writer
}

// NewJSONLogFormatter returns a JSONLogFormatter writing to w. If w is nil, os.Stdout is used.
func NewJSONLogFormatter(w io.Writer) *JSONLogFormatter {
	if w == nil {
		w = os.Stdout
	}
	return &JSONLogFormatter{Out: w}
}

// NewLogEntry implements chimw.LogFormatter.
func (f *JSONLogFormatter) NewLogEntry(r *http.Request) chimw.LogEntry {
	return &jsonLogEntry{
		formatter:  f,
		method:     r.Method,
		path:       r.URL.Path,
		rawQuery:   r.URL.RawQuery,
		remoteAddr: r.RemoteAddr,
		requestID:  chimw.GetReqID(r.Context()),
		userAgent:  r.UserAgent(),
	}
}

type jsonLogEntry struct {
	formatter  *JSONLogFormatter
	method     string
	path       string
	rawQuery   string
	remoteAddr string
	requestID  string
	userAgent  string
}

type jsonLogLine struct {
	Time       string  `json:"time"`
	Level      string  `json:"level"`
	Method     string  `json:"method,omitempty"`
	Path       string  `json:"path,omitempty"`
	Query      string  `json:"query,omitempty"`
	RemoteAddr string  `json:"remote_addr,omitempty"`
	RequestID  string  `json:"request_id,omitempty"`
	UserAgent  string  `json:"user_agent,omitempty"`
	Status     int     `json:"status,omitempty"`
	Bytes      int     `json:"bytes,omitempty"`
	ElapsedMs  float64 `json:"elapsed_ms,omitempty"`
}

// Write implements chimw.LogEntry.
func (e *jsonLogEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	line := jsonLogLine{
		Time:       time.Now().UTC().Format(time.RFC3339Nano),
		Level:      "info",
		Method:     e.method,
		Path:       e.path,
		RemoteAddr: e.remoteAddr,
		RequestID:  e.requestID,
		UserAgent:  e.userAgent,
		Status:     status,
		Bytes:      bytes,
		ElapsedMs:  float64(elapsed.Microseconds()) / 1000,
	}
	if e.rawQuery != "" {
		line.Query = e.rawQuery
	}
	_ = writeJSONLine(e.formatter.Out, &line)
}

// Panic implements chimw.LogEntry.
func (e *jsonLogEntry) Panic(v interface{}, stack []byte) {
	line := struct {
		Time  string `json:"time"`
		Level string `json:"level"`
		Panic string `json:"panic"`
		Stack string `json:"stack,omitempty"`
	}{
		Time:  time.Now().UTC().Format(time.RFC3339Nano),
		Level: "panic",
		Panic: fmt.Sprint(v),
		Stack: string(stack),
	}
	_ = writeJSONLine(e.formatter.Out, &line)
}

// writeJSONLine writes a JSON object to the writer.
func writeJSONLine(w io.Writer, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	_, err = w.Write(b)
	return err
}

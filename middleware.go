package main

import (
	"encoding/json"
	"os"
	"time"

	server "github.com/Elizabethppppp/tcp_server"
)

type ResponseLog struct {
	server.ResponseWriter
	size    int
	status  int
	headers map[string]string
	body    []byte
}

func (w *ResponseLog) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *ResponseLog) Write(p []byte) (int, error) {
	w.size += len(p)
	w.body = append(w.body, p...)
	return w.ResponseWriter.Write(p)
}

func (w *ResponseLog) SetHeader(key, value string) {
	if w.headers == nil {
		w.headers = make(map[string]string)
	}
	w.headers[key] = value
	w.ResponseWriter.SetHeader(key, value)
}

func (w *ResponseLog) GetHeaders() map[string]string {
	return w.headers
}

func (w *ResponseLog) GetBody() string {
	return string(w.body)
}

func LoggerMiddleware(hand server.HandlerFunc) server.HandlerFunc {
	return func(w server.ResponseWriter, r *server.Request) {

		start := time.Now()

		logFields := []any{"method", r.Method, "path", r.Path}

		if len(r.Headers) > 0 {
			logFields = append(logFields, "headers", r.Headers)
		}

		if r.Body != nil && len(r.Body) > 0 {
			bodyStr := r.Body
			logFields = append(logFields, "body", bodyStr)
		}

		writeRequest("REQUEST", logFields)

		rw := &ResponseLog{ResponseWriter: w,
			size:    0,
			status:  200,
			headers: make(map[string]string),
			body:    []byte{}}

		hand(rw, r)

		responseFields := map[string]interface{}{
			"method":      r.Method,
			"path":        r.Path,
			"status":      rw.status,
			"size":        rw.size,
			"duration_ms": time.Since(start).Milliseconds(),
		}
		if len(rw.GetHeaders()) > 0 {
			responseFields["headers"] = rw.GetHeaders()
		}

		if body := rw.GetBody(); body != "" {
			responseFields["body"] = json.RawMessage(body) // если тело в JSON
		}

		writeResponse("RESPONSE", responseFields)
	}
}

func writeRequest(label string, fields []any) {
	logMap := make(map[string]interface{})
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			if key, ok := fields[i].(string); ok {
				logMap[key] = fields[i+1]
			}
		}
	}
	logMap["time"] = time.Now().Format(time.RFC3339)
	writeResponse(label, logMap)
}

func writeResponse(label string, logMap map[string]interface{}) {
	jsonData, err := json.MarshalIndent(logMap, "", "  ")
	if err != nil {
		return
	}

	if err := os.MkdirAll("logs", 0755); err != nil {
		return
	}

	file, err := os.OpenFile("logs/requests.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	separator := "--" + label + "--\n"
	file.WriteString(separator)
	file.Write(append(jsonData, '\n'))

}

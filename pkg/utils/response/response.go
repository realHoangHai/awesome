package response

import (
	"encoding/json"
	"fmt"
	"github.com/realHoangHai/awesome/pkg/status"
	"net/http"
)

// Write write the status code and body on a http ResponseWriter
func Write(w http.ResponseWriter, contentType string, code int, body []byte) {
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(body)))
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	_, _ = w.Write(body)
}

// WriteError write the status code and the error in JSON format on a http ResponseWriter.
// For writing error as plain text or other formats, use Write.
func WriteError(w http.ResponseWriter, code int, err error) {
	Write(w, "application/json", code, status.JSON(err))
}

// WriteJSON write status and JSON data to http ResponseWriter.
func WriteJSON(w http.ResponseWriter, code int, data interface{}) {
	if err, ok := data.(error); ok {
		WriteError(w, code, err)
		return
	}
	b, err := json.Marshal(data)
	if err != nil {
		WriteError(w, code, status.Internal("http: write json, err: %v", err))
		return
	}
	Write(w, "application/json", code, b)
}

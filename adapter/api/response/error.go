package response

import (
	"encoding/json"
	"net/http"
)

// Error defines the structure of success for http responses
type Error struct {
	statusCode int
	Errors     []string `json:"errors,omitempty"`
}

// NewError creates new Error
func NewError(msg []string, status int) *Error {
	return &Error{
		statusCode: status,
		Errors:     msg,
	}
}

// Send returns a response with JSON format
func (e Error) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.statusCode)
	return json.NewEncoder(w).Encode(e)
}

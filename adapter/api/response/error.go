package response

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	statusCode int
	Errors     []map[string]string `json:"errors,omitempty"`
	Error      string              `json:"error,omitempty"`
}

func (e Error) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.statusCode)
	return json.NewEncoder(w).Encode(e)
}

func NewError(err error, status int) *Error {
	return &Error{
		statusCode: status,
		Error:      err.Error(),
	}
}

func NewErrors(errs map[string]string, status int) *Error {
	var errors = make([]map[string]string, 0)
	errors = append(errors, errs)

	return &Error{
		statusCode: status,
		Errors:     errors,
	}
}

package router

import "github.com/gorilla/mux"

// NewGorillaMux create new mux.Router
func NewGorillaMux() *mux.Router {
	return mux.NewRouter()
}

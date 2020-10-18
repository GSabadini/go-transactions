package router

import "github.com/gorilla/mux"

// NewGorrilaMux create new mux.Router
func NewGorillaMux() *mux.Router {
	return mux.NewRouter()
}

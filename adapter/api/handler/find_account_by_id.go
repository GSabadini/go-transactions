package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/GSabadini/go-transactions/adapter/api/response"
	"github.com/GSabadini/go-transactions/usecase"
	"github.com/gorilla/mux"
)

// FindAccountByIDHandler defines the dependencies of the HTTP handler for the use case
type FindAccountByIDHandler struct {
	uc  usecase.FindAccountByIDUseCase
	log *log.Logger
}

// NewFindAccountByIDHandler creates new FindAccountByIDHandler with its dependencies
func NewFindAccountByIDHandler(uc usecase.FindAccountByIDUseCase, log *log.Logger) FindAccountByIDHandler {
	return FindAccountByIDHandler{
		uc:  uc,
		log: log,
	}
}

// Handle handles http request
func (f FindAccountByIDHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["account_id"]

	if ID == "" {
		err := errors.New("invalid account id")
		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}

	output, err := f.uc.Execute(r.Context(), usecase.FindAccountByIDInput{ID: ID})
	if err != nil {
		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	response.NewSuccess(output, http.StatusOK).Send(w)
}

package handler

import (
	"log"
	"net/http"

	"github.com/GSabadini/go-transactions/adapter/api/response"
	"github.com/GSabadini/go-transactions/domain"
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
		response.NewError([]string{"invalid account id"}, http.StatusBadRequest).Send(w)
		return
	}

	output, err := f.uc.Execute(r.Context(), usecase.FindAccountByIDInput{ID: ID})
	if err != nil {
		f.log.Println("failed to find account:", err)
		switch err {
		case domain.ErrAccountNotFound:
			response.NewError([]string{err.Error()}, http.StatusNotFound).Send(w)
			return
		default:
			response.NewError([]string{err.Error()}, http.StatusInternalServerError).Send(w)
			return
		}
	}

	f.log.Println("success to find account")
	response.NewSuccess(output, http.StatusOK).Send(w)
}

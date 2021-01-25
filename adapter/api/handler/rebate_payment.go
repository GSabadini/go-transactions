package handler

import (
	"log"
	"net/http"

	"github.com/GSabadini/go-transactions/adapter/api/response"
	"github.com/GSabadini/go-transactions/domain"
	"github.com/GSabadini/go-transactions/usecase"
	"github.com/gorilla/mux"
)

// RebatePaymentHandler defines the dependencies of the HTTP handler for the use case
type RebatePaymentHandler struct {
	uc  usecase.RebatePaymentUseCase
	log *log.Logger
}

// NewRebatePaymentHandler creates new RebatePaymentHandler with its dependencies
func NewRebatePaymentHandler(uc usecase.RebatePaymentUseCase, log *log.Logger) RebatePaymentHandler {
	return RebatePaymentHandler{
		uc:  uc,
		log: log,
	}
}

// Handle handles http request
func (f RebatePaymentHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["account_id"]

	if ID == "" {
		response.NewError([]string{"invalid account id"}, http.StatusBadRequest).Send(w)
		return
	}

	err := f.uc.Execute(r.Context(), usecase.RebatePaymentInput{AccountID: ID})
	if err != nil {
		f.log.Println("failed to rebate payment account:", err)
		switch err {
		case domain.ErrAccountNotFound:
			response.NewError([]string{err.Error()}, http.StatusNotFound).Send(w)
			return
		default:
			response.NewError([]string{err.Error()}, http.StatusInternalServerError).Send(w)
			return
		}
	}

	f.log.Println("success to rebate payment account")
	response.NewSuccess(nil, http.StatusOK).Send(w)
}

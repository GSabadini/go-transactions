package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GSabadini/go-transactions/adapter/api/response"
	"github.com/GSabadini/go-transactions/domain"
	"github.com/GSabadini/go-transactions/infrastructure/validation"
	"github.com/GSabadini/go-transactions/usecase"

	"github.com/go-playground/validator/v10"
)

// CreateCashoutHandler defines the dependencies of the HTTP handler for the use case
type CreateCashoutHandler struct {
	uc        usecase.CreateAuthorizationUseCase
	log       *log.Logger
	validator *validator.Validate
}

// NewCreateCashoutHandler creates new CreateCashoutHandler with its dependencies
func NewCreateCashoutHandler(
	uc usecase.CreateAuthorizationUseCase,
	log *log.Logger,
	v *validator.Validate,
) CreateCashoutHandler {
	return CreateCashoutHandler{
		uc:        uc,
		log:       log,
		validator: v,
	}
}

// Handle exposes the http handler
func (c CreateCashoutHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateAuthorizationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		c.log.Println("failed to marshal message:", err)
		response.NewError([]string{err.Error()}, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	if err := c.validator.Struct(input); err != nil {
		errs := validation.ErrMessages(err)
		c.log.Println("invalid input:", errs)
		response.NewError(errs, http.StatusBadRequest).Send(w)
		return
	}

	output, err := c.uc.Execute(r.Context(), input)
	if err != nil {
		c.log.Println("failed to creating transaction:", err)
		switch err {
		case domain.ErrOperationInvalid:
			response.NewError([]string{err.Error()}, http.StatusUnprocessableEntity).Send(w)
			return
		case domain.ErrAccountInsufficientCreditLimit:
			response.NewError([]string{err.Error()}, http.StatusUnprocessableEntity).Send(w)
			return
		default:
			response.NewError([]string{err.Error()}, http.StatusInternalServerError).Send(w)
			return
		}
	}

	c.log.Println("success to creating cashout")
	response.NewSuccess(output, http.StatusCreated).Send(w)
}

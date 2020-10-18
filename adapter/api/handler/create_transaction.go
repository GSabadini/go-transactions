package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GSabadini/go-transactions/adapter/api/response"
	"github.com/GSabadini/go-transactions/infrastructure/validation"
	"github.com/GSabadini/go-transactions/usecase"
	"github.com/go-playground/validator/v10"
)

// CreateTransactionHandler defines the dependencies of the HTTP handler for the use case
type CreateTransactionHandler struct {
	uc        usecase.CreateTransactionUseCase
	log       *log.Logger
	validator *validator.Validate
}

// NewCreateTransactionHandler creates new CreateTransactionHandler with its dependencies
func NewCreateTransactionHandler(
	uc usecase.CreateTransactionUseCase,
	log *log.Logger,
	v *validator.Validate,
) CreateTransactionHandler {
	return CreateTransactionHandler{
		uc:        uc,
		log:       log,
		validator: v,
	}
}

// Handler exposes the http handler
func (c CreateTransactionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateTransactionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		c.log.Println("failed to marshal message:", err)
		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	if err := c.validator.Struct(input); err != nil {
		errs := validation.TranslateErr(err)
		c.log.Println("invalid input:", errs)
		response.NewErrors(errs, http.StatusBadRequest).Send(w)
		return
	}

	output, err := c.uc.Execute(r.Context(), input)
	if err != nil {
		c.log.Println("failed to creating account:", err)
		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	c.log.Println("success to creating account")
	response.NewSuccess(output, http.StatusCreated).Send(w)
}
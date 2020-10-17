package handler

import (
	"encoding/json"
	"github.com/GSabadini/go-transactions/adapter/api/response"
	"github.com/GSabadini/go-transactions/usecase"
	"log"
	"net/http"
)

// CreateAccountHandler
type CreateAccountHandler struct {
	uc  usecase.CreateAccountUseCase
	log *log.Logger
}

// NewCreateAccountHandler
func NewCreateAccountHandler(uc usecase.CreateAccountUseCase, log *log.Logger) CreateAccountHandler {
	return CreateAccountHandler{
		uc:  uc,
		log: log,
	}
}

// Handle
func (c CreateAccountHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var input usecase.AccountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		c.log.Println("failed to marshal message:", err)
		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	output, err := c.uc.Execute(r.Context(), input)
	if err != nil {
		c.log.Println("failed to creating account:", err)
		response.NewError(err, http.StatusInternalServerError).Send(w)
		return
	}

	c.log.Println("success to creating account")
	response.NewSuccess(output, http.StatusCreated).Send(w)
}

package presenter

import (
	"time"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/GSabadini/go-transactions/usecase"
)

type createCashoutPresenter struct{}

// NewCreateCashoutPresenter creates new createCashoutPresenter
func NewCreateCashoutPresenter() usecase.CreateAuthorizationPresenter {
	return createCashoutPresenter{}
}

// Output returns the cashout creation response
func (c createCashoutPresenter) Output(authorization domain.Authorization) usecase.CreateAuthorizationOutput {
	return usecase.CreateAuthorizationOutput{
		ID:        authorization.ID(),
		AccountID: authorization.AccountID(),
		Operation: usecase.CreateAuthorizationOperationOutput{
			ID:          authorization.Operation().ID(),
			Description: authorization.Operation().Description(),
			Type:        authorization.Operation().Type(),
		},
		Amount:    authorization.Amount(),
		CreatedAt: authorization.CreatedAt().Format(time.RFC3339),
	}
}

package presenter

import (
	"time"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/GSabadini/go-transactions/usecase"
)

type createAccountPresenter struct{}

// NewCreateAccountPresenter creates new createAccountPresenter
func NewCreateAccountPresenter() usecase.CreateAccountPresenter {
	return createAccountPresenter{}
}

// Output returns the account creation response
func (c createAccountPresenter) Output(account domain.Account) usecase.CreateAccountOutput {
	return usecase.CreateAccountOutput{
		ID: account.ID(),
		Document: usecase.CreateAccountDocumentOutput{
			Number: account.Document().Number(),
		},
		AvailableCreditLimit: account.AvailableCreditLimit(),
		CreatedAt:            account.CreatedAt().Format(time.RFC3339),
	}
}

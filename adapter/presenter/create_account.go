package presenter

import (
	"github.com/GSabadini/go-transactions/domain"
	"github.com/GSabadini/go-transactions/usecase"
	"time"
)

type createAccountPresenter struct{}

// NewCreateAccountPresenter creates new createAccountPresenter
func NewCreateAccountPresenter() usecase.CreateAccountPresenter {
	return createAccountPresenter{}
}

// Output returns the account creation response
func (c createAccountPresenter) Output(account domain.Account) usecase.AccountOutput {
	return usecase.AccountOutput{
		ID: account.ID(),
		Document: usecase.AccountDocumentOutput{
			Number: account.Document().Number(),
		},
		CreatedAt: account.CreatedAt().Format(time.RFC3339),
	}
}
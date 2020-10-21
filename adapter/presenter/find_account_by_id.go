package presenter

import (
	"time"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/GSabadini/go-transactions/usecase"
)

type findAccountByIDPresenter struct{}

// NewFindAccountByIDPresenter creates new findAccountByIDPresenter
func NewFindAccountByIDPresenter() usecase.FindAccountByIDPresenter {
	return findAccountByIDPresenter{}
}

// Output returns the account fetch response by ID
func (f findAccountByIDPresenter) Output(account domain.Account) usecase.FindAccountByIDOutput {
	return usecase.FindAccountByIDOutput{
		ID:                   account.ID(),
		AvailableCreditLimit: account.AvailableCreditLimit(),
		Document: usecase.FindAccountByIDDocumentOutput{
			Number: account.Document().Number(),
		},
		CreatedAt: account.CreatedAt().Format(time.RFC3339),
	}
}

package presenter

import (
	"time"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/GSabadini/go-transactions/usecase"
)

type createTransactionPresenter struct{}

// NewCreateTransactionPresenter creates new createTransactionPresenter
func NewCreateTransactionPresenter() usecase.CreateTransactionPresenter {
	return createTransactionPresenter{}
}

// Output returns the transaction creation response
func (c createTransactionPresenter) Output(transaction domain.Transaction) usecase.CreateTransactionOutput {
	return usecase.CreateTransactionOutput{
		ID:        transaction.ID(),
		AccountID: transaction.AccountID(),
		Operation: usecase.CreateTransactionOperationOutput{
			ID:          transaction.Operation().ID(),
			Description: transaction.Operation().Description(),
			Type:        transaction.Operation().Type(),
		},
		Amount:    transaction.Amount(),
		Balance:   transaction.Balance(),
		CreatedAt: transaction.CreatedAt().Format(time.RFC3339),
	}
}

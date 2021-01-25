package presenter

import (
	"github.com/GSabadini/go-transactions/usecase"
)

type rebatePaymentPresenter struct{}

// NewRebatePaymentPresenter creates new rebatePaymentPresenter
func NewRebatePaymentPresenter() usecase.RebatePaymentPresenter {
	return rebatePaymentPresenter{}
}

// Output returns the account fetch response by ID
func (f rebatePaymentPresenter) Output(err error) error {
	return err
}

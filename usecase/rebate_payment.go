package usecase

import (
	"context"
	"time"

	"github.com/GSabadini/go-transactions/domain"
)

type (
	// Input port
	RebatePaymentUseCase interface {
		Execute(context.Context, RebatePaymentInput) error
	}

	// Input data
	RebatePaymentInput struct {
		AccountID string
	}

	// Output port
	RebatePaymentPresenter interface {
		Output(error) error
	}

	rebatePaymentInteractor struct {
		repoTransactionUpdater domain.TransactionUpdater
		repoTransactionFinder  domain.TransactionFinder
		pre                    RebatePaymentPresenter
		ctxTimeout             time.Duration
	}
)

// NewRebatePaymentInteractor creates new rebatePaymentInteractor with its dependencies
func NewRebatePaymentInteractor(
	repoTransactionUpdater domain.TransactionUpdater,
	repoTransactionFinder domain.TransactionFinder,
	pre RebatePaymentPresenter,
	ctxTimeout time.Duration,
) RebatePaymentUseCase {
	return rebatePaymentInteractor{
		repoTransactionUpdater: repoTransactionUpdater,
		repoTransactionFinder:  repoTransactionFinder,
		pre:                    pre,
		ctxTimeout:             ctxTimeout,
	}
}

// Execute orchestrates the use case
func (r rebatePaymentInteractor) Execute(ctx context.Context, input RebatePaymentInput) error {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	_, err := r.repoTransactionFinder.FindByAccountID(ctx, input.AccountID)
	if err != nil {
		return r.pre.Output(err)
	}

	//for _, transaction := range transactions {
	//	if transaction.Operation().Type() == domain.Credit {
	//		for _, transaction2 := range transactions {
	//			if transaction2.Operation().Type() == domain.Debit && transaction.Balance() > 0 {
	//				balance := transaction.Balance() + transaction2.Balance()
	//				if err := r.repoTransactionUpdater.UpdateBalance(ctx, transaction2.ID(), balance); err != nil {
	//					return err
	//				}
	//			}
	//		}
	//	}
	//}

	return nil
}

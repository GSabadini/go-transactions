package usecase

import (
	"context"
	"time"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/google/uuid"
)

type (
	//  Input port
	CreateTransactionUseCase interface {
		Execute(context.Context, CreateTransactionInput) (CreateTransactionOutput, error)
	}

	// Input data
	CreateTransactionInput struct {
		AccountID   string `json:"account_id" validate:"required"`
		OperationID string `json:"operation_id" validate:"required"`
		Amount      int64  `json:"amount" validate:"required,gt=0"`
	}

	// Output port
	CreateTransactionPresenter interface {
		Output(domain.Transaction) CreateTransactionOutput
	}

	// Output data
	CreateTransactionOutput struct {
		ID        string                           `json:"id"`
		AccountID string                           `json:"account_id"`
		Operation CreateTransactionOperationOutput `json:"operation"`
		Amount    int64                            `json:"amount"`
		Balance   int64                            `json:"balance"`
		CreatedAt string                           `json:"created_at"`
	}

	// Output data
	CreateTransactionOperationOutput struct {
		ID          string `json:"id"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}

	createTransactionInteractor struct {
		repoTransactionCreator domain.TransactionCreator
		repoAccountFinder      domain.AccountFinder
		repoAccountUpdater     domain.AccountUpdater
		pre                    CreateTransactionPresenter
		ctxTimeout             time.Duration
	}
)

// NewCreateTransactionInteractor creates new createTransactionInteractor with its dependencies
func NewCreateTransactionInteractor(
	repoTransactionCreator domain.TransactionCreator,
	repoAccountFinder domain.AccountFinder,
	repoAccountUpdater domain.AccountUpdater,
	pre CreateTransactionPresenter,
	ctxTimeout time.Duration,
) CreateTransactionUseCase {
	return createTransactionInteractor{
		repoTransactionCreator: repoTransactionCreator,
		repoAccountFinder:      repoAccountFinder,
		repoAccountUpdater:     repoAccountUpdater,
		pre:                    pre,
		ctxTimeout:             ctxTimeout,
	}
}

// Execute orchestrates the use case
func (c createTransactionInteractor) Execute(ctx context.Context, i CreateTransactionInput) (CreateTransactionOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	var (
		account     domain.Account
		transaction domain.Transaction
		err         error
	)

	op, err := domain.NewOperation(i.OperationID)
	if err != nil {
		return c.pre.Output(domain.Transaction{}), err
	}

	err = c.repoTransactionCreator.WithTransaction(ctx, func(ctxTx context.Context) error {
		account, err = c.repoAccountFinder.FindByID(ctxTx, i.AccountID)
		if err != nil {
			return err
		}

		if err = account.PaymentOperation(i.Amount, op.Type()); err != nil {
			return err
		}

		if err = c.repoAccountUpdater.UpdateCreditLimit(ctxTx, account.ID(), account.AvailableCreditLimit()); err != nil {
			return err
		}

		balance := i.Amount

		transaction, err = c.repoTransactionCreator.Create(ctxTx, domain.NewTransaction(
			uuid.New().String(),
			i.AccountID,
			op,
			i.Amount,
			balance,
			time.Now(),
		))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return c.pre.Output(domain.Transaction{}), err
	}

	return c.pre.Output(transaction), nil
}

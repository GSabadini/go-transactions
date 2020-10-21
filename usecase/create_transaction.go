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
		AccountID   string  `json:"account_id" validate:"required"`
		OperationID string  `json:"operation_id" validate:"required"`
		Amount      float64 `json:"amount" validate:"required,gt=0"`
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
		Amount    float64                          `json:"amount"`
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

	op, err := domain.NewOperation(i.OperationID)
	if err != nil {
		return c.pre.Output(domain.Transaction{}), err
	}

	account, err := c.repoAccountFinder.FindByID(ctx, i.AccountID)
	if err != nil {
		return c.pre.Output(domain.Transaction{}), err
	}

	if err := account.PaymentOperation(i.Amount, op.Type()); err != nil {
		return c.pre.Output(domain.Transaction{}), err
	}

	if err := c.repoAccountUpdater.UpdateCreditLimit(ctx, account.ID(), account.AvailableCreditLimit()); err != nil {
		return c.pre.Output(domain.Transaction{}), err
	}

	transaction, err := c.repoTransactionCreator.Create(ctx, domain.NewTransaction(
		uuid.New().String(),
		i.AccountID,
		op,
		i.Amount,
		time.Now(),
	))
	if err != nil {
		return c.pre.Output(domain.Transaction{}), err
	}

	return c.pre.Output(transaction), nil
}

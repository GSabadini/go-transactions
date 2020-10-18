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
		repo       domain.CreateTransactionRepository
		pre        CreateTransactionPresenter
		ctxTimeout time.Duration
	}
)

// NewCreateTransactionInteractor creates new createTransactionInteractor with its dependencies
func NewCreateTransactionInteractor(
	repo domain.CreateTransactionRepository,
	pre CreateTransactionPresenter,
	ctxTimeout time.Duration,
) CreateTransactionUseCase {
	return createTransactionInteractor{
		repo:       repo,
		pre:        pre,
		ctxTimeout: ctxTimeout,
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

	transaction, err := c.repo.Create(ctx, domain.NewTransaction(
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

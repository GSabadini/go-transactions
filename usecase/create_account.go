package usecase

import (
	"context"
	"time"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/google/uuid"
)

type (
	// Input port
	CreateAccountUseCase interface {
		Execute(context.Context, CreateAccountInput) (CreateAccountOutput, error)
	}

	// Input data
	CreateAccountInput struct {
		Document struct {
			Number string `json:"number" validate:"required,max=30"`
		}
		AvailableCreditLimit int64 `json:"available_credit_limit" validate:"required,gt=0"`
	}

	// Output port
	CreateAccountPresenter interface {
		Output(domain.Account) CreateAccountOutput
	}

	// Output data
	CreateAccountOutput struct {
		ID                   string                      `json:"id"`
		AvailableCreditLimit int64                       `json:"available_credit_limit"`
		Document             CreateAccountDocumentOutput `json:"document"`
		CreatedAt            string                      `json:"created_at"`
	}

	// Output data
	CreateAccountDocumentOutput struct {
		Number string `json:"number"`
	}

	createAccountInteractor struct {
		repo       domain.AccountCreator
		pre        CreateAccountPresenter
		ctxTimeout time.Duration
	}
)

// NewCreateAccountInteractor creates new createAccountInteractor with its dependencies
func NewCreateAccountInteractor(
	repo domain.AccountCreator,
	pre CreateAccountPresenter,
	ctxTimeout time.Duration,
) CreateAccountUseCase {
	return createAccountInteractor{
		repo:       repo,
		pre:        pre,
		ctxTimeout: ctxTimeout,
	}
}

// Execute orchestrates the use case
func (c createAccountInteractor) Execute(ctx context.Context, i CreateAccountInput) (CreateAccountOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	account, err := c.repo.Create(ctx, domain.NewAccount(
		uuid.New().String(),
		i.Document.Number,
		i.AvailableCreditLimit,
		time.Now(),
	))
	if err != nil {
		return c.pre.Output(domain.Account{}), err
	}

	return c.pre.Output(account), nil
}

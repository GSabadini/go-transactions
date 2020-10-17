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
		Execute(context.Context, AccountInput) (AccountOutput, error)
	}

	// Input data
	AccountInput struct {
		Document struct {
			Number string
		}
	}

	//Output port
	CreateAccountPresenter interface {
		Output(domain.Account) AccountOutput
	}

	// Output data
	AccountOutput struct {
		ID        string                `json:"id"`
		Document  AccountDocumentOutput `json:"document"`
		CreatedAt string                `json:"created_at"`
	}

	// Output data
	AccountDocumentOutput struct {
		Number string `json:"number"`
	}

	createAccountInteractor struct {
		repo       domain.CreateAccountRepository
		pre        CreateAccountPresenter
		ctxTimeout time.Duration
	}
)

// NewCreateAccountInteractor creates new createAccountInteractor with its dependencies
func NewCreateAccountInteractor(
	repo domain.CreateAccountRepository,
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
func (c createAccountInteractor) Execute(ctx context.Context, i AccountInput) (AccountOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	account, err := c.repo.Create(ctx, domain.NewAccount(uuid.New().String(), i.Document.Number, time.Now()))
	if err != nil {
		return c.pre.Output(domain.Account{}), err
	}

	return c.pre.Output(account), nil
}

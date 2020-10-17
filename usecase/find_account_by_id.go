package usecase

import (
	"context"
	"time"

	"github.com/GSabadini/go-transactions/domain"
)

type (
	// Input port
	FindAccountByIDUseCase interface {
		Execute(context.Context, FindAccountByIDInput) (FindAccountByIDOutput, error)
	}

	// Input data
	FindAccountByIDInput struct {
		ID string
	}

	// Output port
	FindAccountByIDPresenter interface {
		Output(domain.Account) FindAccountByIDOutput
	}

	// Output data
	FindAccountByIDOutput struct {
		ID        string                        `json:"id"`
		Document  FindAccountByIDDocumentOutput `json:"document"`
		CreatedAt string                        `json:"created_at"`
	}

	// Output data
	FindAccountByIDDocumentOutput struct {
		Number string `json:"number"`
	}

	findAccountByIDInteractor struct {
		repo       domain.FindAccountByIDRepository
		pre        FindAccountByIDPresenter
		ctxTimeout time.Duration
	}
)

// NewFindAccountByIDInteractor creates new findAccountByIDInteractor with its dependencies
func NewFindAccountByIDInteractor(
	repo domain.FindAccountByIDRepository,
	pre FindAccountByIDPresenter,
	ctxTimeout time.Duration,
) FindAccountByIDUseCase {
	return findAccountByIDInteractor{
		repo:       repo,
		pre:        pre,
		ctxTimeout: ctxTimeout,
	}
}

// Execute orchestrates the use case
func (f findAccountByIDInteractor) Execute(ctx context.Context, i FindAccountByIDInput) (FindAccountByIDOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, f.ctxTimeout)
	defer cancel()

	account, err := f.repo.FindByID(ctx, i.ID)
	if err != nil {
		return f.pre.Output(domain.Account{}), err
	}

	return f.pre.Output(account), nil
}

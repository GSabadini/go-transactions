package usecase

import (
	"context"
	"time"

	"github.com/GSabadini/go-transactions/domain"

	"github.com/google/uuid"
)

type (
	Evaluator interface {
		Rules(context.Context, string) error
	}

	External interface {
		Acquirer(context.Context, string) error
	}

	Event interface {
		Audit(context.Context, domain.Authorization)
		Transactions(context.Context, domain.Authorization)
	}

	//  Input port
	CreateAuthorizationUseCase interface {
		Execute(context.Context, CreateAuthorizationInput) (CreateAuthorizationOutput, error)
	}

	// Input data
	CreateAuthorizationInput struct {
		AccountID   string `json:"account_id" validate:"required"`
		OperationID string `json:"operation_id" validate:"required"`
		Amount      int64  `json:"amount" validate:"required,gt=0"`
		External    bool   `json:"external" validate:"required,bool"`
	}

	// Output port
	CreateAuthorizationPresenter interface {
		Output(domain.Authorization) CreateAuthorizationOutput
	}

	// Output data
	CreateAuthorizationOutput struct {
		ID        string                             `json:"id"`
		AccountID string                             `json:"account_id"`
		Operation CreateAuthorizationOperationOutput `json:"operation"`
		Amount    int64                              `json:"amount"`
		CreatedAt string                             `json:"created_at"`
	}

	// Output data
	CreateAuthorizationOperationOutput struct {
		ID          string `json:"id"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}

	createAuthorizationInteractor struct {
		repoAuthorizationCreator domain.AuthorizationCreator
		repoAccountFinder        domain.AccountFinder
		repoAccountUpdater       domain.AccountUpdater
		evaluator                Evaluator
		external                 External
		event                    Event
		pre                      CreateAuthorizationPresenter
		ctxTimeout               time.Duration
	}
)

// NewCreateAuthorizationInteractor creates new createAuthorizationInteractor with its dependencies
func NewCreateAuthorizationInteractor(
	repoAuthorizationCreator domain.AuthorizationCreator,
	repoAccountFinder domain.AccountFinder,
	repoAccountUpdater domain.AccountUpdater,
	evaluator Evaluator,
	external External,
	event Event,
	pre CreateAuthorizationPresenter,
	ctxTimeout time.Duration,
) CreateAuthorizationUseCase {
	return createAuthorizationInteractor{
		repoAuthorizationCreator: repoAuthorizationCreator,
		repoAccountFinder:        repoAccountFinder,
		repoAccountUpdater:       repoAccountUpdater,
		evaluator:                evaluator,
		external:                 external,
		event:                    event,
		pre:                      pre,
		ctxTimeout:               ctxTimeout,
	}
}

// Execute orchestrates the use case
func (c createAuthorizationInteractor) Execute(ctx context.Context, i CreateAuthorizationInput) (CreateAuthorizationOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	var (
		account       domain.Account
		authorization domain.Authorization
		err           error
	)

	op, err := domain.NewOperation(i.OperationID)
	if err != nil {
		return c.pre.Output(domain.Authorization{}), err
	}

	err = c.repoAuthorizationCreator.WithTransaction(ctx, func(ctxTx context.Context) error {
		account, err = c.repoAccountFinder.FindByID(ctxTx, i.AccountID)
		if err != nil {
			return err
		}

		if err = c.evaluator.Rules(ctx, i.AccountID); err != nil {
			return err
		}

		if i.External {
			if err = c.external.Acquirer(ctx, i.AccountID); err != nil {
				return err
			}

			authorization, err = c.repoAuthorizationCreator.Create(ctxTx, domain.NewAuthorization(
				uuid.New().String(),
				i.AccountID,
				op,
				i.Amount,
				time.Now(),
			))
			if err != nil {
				return err
			}

			return nil
		}

		if err = account.PaymentOperation(i.Amount, op.Type()); err != nil {
			return err
		}

		if err = c.repoAccountUpdater.UpdateCreditLimit(ctxTx, account.ID(), account.AvailableCreditLimit()); err != nil {
			return err
		}

		authorization, err = c.repoAuthorizationCreator.Create(ctxTx, domain.NewAuthorization(
			uuid.New().String(),
			i.AccountID,
			op,
			i.Amount,
			time.Now(),
		))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return c.pre.Output(domain.Authorization{}), err
	}

	go c.event.Audit(ctx, authorization)
	go c.event.Transactions(ctx, authorization)

	return c.pre.Output(authorization), nil
}

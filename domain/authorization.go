package domain

import (
	"context"
	"time"
)

type (
	// AuthorizationCreator defines the operation of creating an authorization entity
	AuthorizationCreator interface {
		Create(context.Context, Authorization) (Authorization, error)
		WithTransaction(context.Context, func(context.Context) error) error
	}

	Authorization struct {
		id        string
		accountID string
		operation Operation
		amount    int64
		createdAt time.Time
	}
)

func NewAuthorization(id string, accID string, op Operation, amount int64, createdAt time.Time) Authorization {
	return Authorization{
		id:        id,
		accountID: accID,
		operation: op,
		amount:    amount,
		createdAt: createdAt,
	}
}

func (a Authorization) ID() string {
	return a.id
}

func (a Authorization) Amount() int64 {
	return a.amount
}

func (a Authorization) AccountID() string {
	return a.accountID
}

func (a Authorization) Operation() Operation {
	return a.operation
}

func (a Authorization) CreatedAt() time.Time {
	return a.createdAt
}

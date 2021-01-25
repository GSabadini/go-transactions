package domain

import (
	"context"
	"time"
)

type (
	// TransactionCreator defines the operation of creating a transaction entity
	TransactionCreator interface {
		Create(context.Context, Transaction) (Transaction, error)
		WithTransaction(context.Context, func(context.Context) error) error
	}

	// TransactionFinder defines the search operation for a transaction entity
	TransactionFinder interface {
		FindByAccountID(context.Context, string) ([]Transaction, error)
	}

	// TransactionUpdater defines the update operation for a transaction entity
	TransactionUpdater interface {
		UpdateBalance(context.Context, string, int64) error
	}

	// Transaction defines the transaction entity
	Transaction struct {
		id        string
		accountID string
		operation Operation
		amount    int64
		balance   int64
		createdAt time.Time
	}
)

// NewTransaction creates new Transaction
func NewTransaction(id string, accID string, op Operation, amount int64, balance int64, createdAt time.Time) Transaction {
	if op.opType == Debit {
		if amount < 0 && balance < 0 {

		} else {
			amount = -amount
			balance = -balance
		}
	}

	return Transaction{
		id:        id,
		accountID: accID,
		operation: op,
		amount:    amount,
		balance:   balance,
		createdAt: createdAt,
	}
}

// ID returns the id property
func (t Transaction) ID() string {
	return t.id
}

// AccountID returns the accountID property
func (t Transaction) AccountID() string {
	return t.accountID
}

// Operation returns the operation property
func (t Transaction) Operation() Operation {
	return t.operation
}

// Amount returns the amount property
func (t Transaction) Amount() int64 {
	return t.amount
}

// Balance returns the balance property
func (t Transaction) Balance() int64 {
	return t.balance
}

// CreatedAt returns the createdAt property
func (t Transaction) CreatedAt() time.Time {
	return t.createdAt
}

package domain

import (
	"context"
	"time"
)

type (
	// TransactionCreator defines the operation of creating a transaction entity
	TransactionCreator interface {
		Create(context.Context, Transaction) (Transaction, error)
	}

	// Transaction defines the transaction entity
	Transaction struct {
		id        string
		accountID string
		operation Operation
		amount    float64
		createdAt time.Time
	}
)

// NewTransaction creates new Transaction
func NewTransaction(id string, accID string, op Operation, amount float64, createdAt time.Time) Transaction {
	if op.opType == Debit {
		amount = -amount
	}

	return Transaction{
		id:        id,
		accountID: accID,
		operation: op,
		amount:    amount,
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
func (t Transaction) Amount() float64 {
	return t.amount
}

// CreatedAt returns the createdAt property
func (t Transaction) CreatedAt() time.Time {
	return t.createdAt
}

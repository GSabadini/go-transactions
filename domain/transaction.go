package domain

import (
	"context"
	"time"
)

type (
	CreateTransaction interface {
		Create(context.Context, Transaction)
	}

	//TransactionInput struct {
	//	AccountID string
	//	OperationTypeID string
	//	Amount int64
	//}

	Transaction struct {
		id              string
		accountID       string
		operationTypeID string
		amount          int64
		eventDate       time.Time
	}
)

func NewTransaction(id string, accID string, opID string, amount int64, eventDate time.Time) Transaction {
	return Transaction{
		id:              id,
		accountID:       accID,
		operationTypeID: opID,
		amount:          amount,
		eventDate:       eventDate,
	}
}

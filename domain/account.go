package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrAccountAlreadyExists           = errors.New("account already exists")
	ErrAccountNotFound                = errors.New("account not found")
	ErrAccountInsufficientCreditLimit = errors.New("credit limit insufficient")
)

type (
	// AccountCreator defines the operation of creating a account entity
	AccountCreator interface {
		Create(context.Context, Account) (Account, error)
	}

	// AccountFinder defines the search operation for a account entity
	AccountFinder interface {
		FindByID(context.Context, string) (Account, error)
	}

	// AccountFinder defines the update operation for a account entity
	AccountUpdater interface {
		UpdateCreditLimit(context.Context, string, int64) error
	}

	// Account defines the account entity
	Account struct {
		id                   string
		document             Document
		availableCreditLimit int64
		createdAt            time.Time
	}

	// Document defines document property
	Document struct {
		number string
	}
)

// NewAccount creates new Account
func NewAccount(ID string, docNumber string, avCreditLimit int64, createdAt time.Time) Account {
	return Account{
		id: ID,
		document: Document{
			number: docNumber,
		},
		availableCreditLimit: avCreditLimit,
		createdAt:            createdAt,
	}
}

// PaymentOperation
func (a *Account) PaymentOperation(amount int64, opType string) error {
	if opType == Debit {
		return a.Withdraw(amount)
	}

	a.Deposit(amount)
	return nil
}

// Deposit
func (a *Account) Deposit(amount int64) {
	a.availableCreditLimit += amount
}

// Withdraw
func (a *Account) Withdraw(amount int64) error {
	if a.availableCreditLimit < amount {
		return ErrAccountInsufficientCreditLimit
	}
	a.availableCreditLimit -= amount
	return nil
}

// ID returns the id property
func (a Account) ID() string {
	return a.id
}

// Document returns the document property
func (a Account) Document() Document {
	return a.document
}

// CreatedAt returns the createdAt property
func (a Account) CreatedAt() time.Time {
	return a.createdAt
}

// AvailableCreditLimit returns the availableCreditLimit property
func (a Account) AvailableCreditLimit() int64 {
	return a.availableCreditLimit
}

// Number returns the number property
func (d Document) Number() string {
	return d.number
}

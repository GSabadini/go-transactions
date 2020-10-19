package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrAccountNotFound      = errors.New("account not found")
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

	// Account defines the account entity
	Account struct {
		id        string
		document  Document
		createdAt time.Time
	}

	// Document defines document property
	Document struct {
		number string
	}
)

// NewAccount creates new Account
func NewAccount(ID string, docNumber string, createdAt time.Time) Account {
	return Account{
		id: ID,
		document: Document{
			number: docNumber,
		},
		createdAt: createdAt,
	}
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

// Number returns the number property
func (d Document) Number() string {
	return d.number
}

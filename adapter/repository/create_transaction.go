package repository

import (
	"context"
	"database/sql"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/pkg/errors"
)

type createTransactionRepository struct {
	db *sql.DB
}

// NewCreateTransactionRepository creates new createTransactionRepository with its dependencies
func NewCreateTransactionRepository(db *sql.DB) domain.TransactionCreator {
	return createTransactionRepository{
		db: db,
	}
}

// Create performs insert into the database
func (c createTransactionRepository) Create(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	if _, err := c.db.ExecContext(
		ctx,
		`INSERT INTO transactions (id, account_id, operation_id, amount, created_at) VALUES (?, ?, ?, ?, ?)`,
		transaction.ID(),
		transaction.AccountID(),
		transaction.Operation().ID(),
		transaction.Amount(),
		transaction.CreatedAt(),
	); err != nil {
		return domain.Transaction{}, errors.Wrap(err, errDatabase.Error())
	}

	return transaction, nil
}

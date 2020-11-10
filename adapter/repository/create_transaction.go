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
	tx, ok := ctx.Value("TransactionContextKey").(*sql.Tx)
	if !ok {
		var err error
		tx, err = c.db.BeginTx(ctx, nil)
		if err != nil {
			return domain.Transaction{}, errors.Wrap(err, errDatabase.Error())
		}
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO transactions (id, account_id, operation_id, amount, balance, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		transaction.ID(),
		transaction.AccountID(),
		transaction.Operation().ID(),
		transaction.Amount(),
		transaction.Balance(),
		transaction.CreatedAt(),
	); err != nil {
		return domain.Transaction{}, errors.Wrap(err, errDatabase.Error())
	}

	return transaction, nil
}

func (c createTransactionRepository) WithTransaction(ctx context.Context, fn func(ctxFn context.Context) error) error {
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, errDatabase.Error())
	}

	ctxTx := context.WithValue(ctx, "TransactionContextKey", tx)
	err = fn(ctxTx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Wrap(err, "rollback error")
		}
		return err
	}

	return tx.Commit()
}

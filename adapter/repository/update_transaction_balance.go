package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/pkg/errors"
)

type updateTransactionBalanceRepository struct {
	db *sql.DB
}

// NewUpdateTransactionBalanceRepository creates new updateTransactionBalanceRepository with its dependencies
func NewUpdateTransactionBalanceRepository(db *sql.DB) domain.TransactionUpdater {
	return updateTransactionBalanceRepository{
		db: db,
	}
}

// UpdateBalance performs update into the database
func (u updateTransactionBalanceRepository) UpdateBalance(ctx context.Context, ID string, balance int64) error {
	tx, ok := ctx.Value("TransactionContextKey").(*sql.Tx)
	if !ok {
		var err error
		tx, err = u.db.BeginTx(ctx, nil)
		if err != nil {
			return errors.Wrap(err, errDatabase.Error())
		}
	}

	fmt.Println(balance, ID)
	if _, err := tx.ExecContext(
		ctx,
		`UPDATE transactions SET balance = ? WHERE id = ?`,
		balance,
		ID,
	); err != nil {
		return errors.Wrap(err, errDatabase.Error())
	}

	return nil
}

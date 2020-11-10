package repository

import (
	"context"
	"database/sql"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/pkg/errors"
)

type updateAccountCreditLimitRepository struct {
	db *sql.DB
}

func NewUpdateAccountCreditLimitRepository(db *sql.DB) domain.AccountUpdater {
	return updateAccountCreditLimitRepository{
		db: db,
	}
}

func (u updateAccountCreditLimitRepository) UpdateCreditLimit(ctx context.Context, ID string, amount int64) error {
	tx, ok := ctx.Value("TransactionContextKey").(*sql.Tx)
	if !ok {
		var err error
		tx, err = u.db.BeginTx(ctx, nil)
		if err != nil {
			return errors.Wrap(err, errDatabase.Error())
		}
	}

	if _, err := tx.ExecContext(
		ctx,
		`UPDATE accounts SET available_credit_limit = ? WHERE id = ?`,
		amount,
		ID,
	); err != nil {
		return errors.Wrap(err, errDatabase.Error())
	}

	return nil
}

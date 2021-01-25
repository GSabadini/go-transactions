package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/pkg/errors"
)

type findTransactionByAccountIDRepository struct {
	db *sql.DB
}

// NewFindTransactionByAccountIDRepository creates new findTransactionByAccountIDRepository with its dependencies
func NewFindTransactionByAccountIDRepository(db *sql.DB) domain.TransactionFinder {
	return findTransactionByAccountIDRepository{
		db: db,
	}
}

// FindByAccountID performs select into the database
func (f findTransactionByAccountIDRepository) FindByAccountID(ctx context.Context, ID string) ([]domain.Transaction, error) {
	tx, ok := ctx.Value("TransactionContextKey").(*sql.Tx)
	if !ok {
		var err error
		tx, err = f.db.BeginTx(ctx, nil)
		if err != nil {
			return []domain.Transaction{}, errors.Wrap(err, errDatabase.Error())
		}
	}

	rows, err := tx.QueryContext(
		ctx,
		"SELECT * FROM transactions WHERE account_id = ? AND NOT operation_id = ?",
		ID,
		domain.Pagamento,
	)

	var transactions = make([]domain.Transaction, 0)
	for rows.Next() {
		var (
			ID          string
			accountID   string
			operationID string
			amount      int64
			balance     int64
			createdAt   time.Time
		)

		if err = rows.Scan(&ID, &accountID, &operationID, &amount, &balance, &createdAt); err != nil {
			return []domain.Transaction{}, errors.Wrap(err, "error listing transfers")
		}

		op, err := domain.NewOperation(operationID)
		if err != nil {
			continue
		}

		transactions = append(transactions, domain.NewTransaction(
			ID,
			accountID,
			op,
			amount,
			balance,
			createdAt,
		))
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return []domain.Transaction{}, err
	}

	return transactions, nil
}

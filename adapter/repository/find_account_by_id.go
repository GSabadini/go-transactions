package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/pkg/errors"
)

type findAccountByIDRepository struct {
	db *sql.DB
}

// NewCreateAccountRepository creates new findAccountByIDRepository with its dependencies
func NewAccountByIDRepository(db *sql.DB) domain.AccountFinder {
	return findAccountByIDRepository{
		db: db,
	}
}

// FindByID performs select into the database
func (f findAccountByIDRepository) FindByID(ctx context.Context, ID string) (domain.Account, error) {
	tx, ok := ctx.Value("TransactionContextKey").(*sql.Tx)
	if !ok {
		var err error
		tx, err = f.db.BeginTx(ctx, nil)
		if err != nil {
			return domain.Account{}, errors.Wrap(err, errDatabase.Error())
		}
	}

	var (
		id            string
		docNumber     string
		avCreditLimit int64
		createdAt     time.Time
	)

	err := tx.QueryRowContext(
		ctx,
		"SELECT * FROM accounts WHERE id = ?",
		ID,
	).Scan(&id, &docNumber, &avCreditLimit, &createdAt)
	switch {
	case err == sql.ErrNoRows:
		return domain.Account{}, domain.ErrAccountNotFound
	default:
		return domain.NewAccount(id, docNumber, avCreditLimit, createdAt), errors.Wrap(err, errDatabase.Error())
	}
}

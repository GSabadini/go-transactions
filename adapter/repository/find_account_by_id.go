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

// NewAccountByIDRepository NewCreateAccountRepository creates new findAccountByIDRepository with its dependencies
func NewAccountByIDRepository(db *sql.DB) domain.AccountFinder {
	return findAccountByIDRepository{
		db: db,
	}
}

// FindByID performs select into the database
func (f findAccountByIDRepository) FindByID(ctx context.Context, ID string) (domain.Account, error) {
	tx, ok := ctx.Value("TxKey").(*sql.Tx)
	if !ok {
		var err error
		tx, err = f.db.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			return domain.Account{}, errors.Wrap(err, errUnknown.Error())
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
		return domain.NewAccount(id, docNumber, avCreditLimit, createdAt), errors.Wrap(err, errUnknown.Error())
	}
}

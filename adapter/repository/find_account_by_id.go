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
	var (
		id            string
		docNumber     string
		avCreditLimit float64
		createdAt     time.Time
	)

	err := f.db.QueryRowContext(
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

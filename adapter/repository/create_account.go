package repository

import (
	"context"
	"database/sql"

	"github.com/GSabadini/go-transactions/domain"
)

type createAccountRepository struct {
	db *sql.DB
}

func NewCreateAccountRepository(db *sql.DB) domain.CreateAccountRepository {
	return createAccountRepository{
		db: db,
	}
}

func (c createAccountRepository) Create(ctx context.Context, account domain.Account) (domain.Account, error) {
	if _, err := c.db.ExecContext(
		ctx,
		`INSERT INTO accounts (id, document_number, created_at) VALUES (?, ?, ?)`,
		account.ID(),
		account.Document().Number(),
		account.CreatedAt(),
	); err != nil {
		//return domain.Account{}, errors.Wrap(err, "error creating account")
		return domain.Account{}, err
	}

	return account, nil
}

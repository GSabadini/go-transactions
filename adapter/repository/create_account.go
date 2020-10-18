package repository

import (
	"context"
	"database/sql"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type createAccountRepository struct {
	db *sql.DB
}

// NewCreateAccountRepository creates new createAccountRepository with its dependencies
func NewCreateAccountRepository(db *sql.DB) domain.CreateAccountRepository {
	return createAccountRepository{
		db: db,
	}
}

// Create performs insert into the database
func (c createAccountRepository) Create(ctx context.Context, account domain.Account) (domain.Account, error) {
	if _, err := c.db.ExecContext(
		ctx,
		`INSERT INTO accounts (id, document_number, created_at) VALUES (?, ?, ?)`,
		account.ID(),
		account.Document().Number(),
		account.CreatedAt(),
	); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == errDupEntry {
				return domain.Account{}, domain.ErrAccountAlreadyExists
			}
		}

		return domain.Account{}, errors.Wrap(err, errDatabase.Error())
	}

	return account, nil
}

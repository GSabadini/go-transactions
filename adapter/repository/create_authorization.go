package repository

import (
	"context"
	"database/sql"

	"github.com/GSabadini/go-transactions/domain"

	"github.com/pkg/errors"
)

type createAuthorizationRepository struct {
	db *sql.DB
}

// NewCreateAuthorizationRepository creates new createAuthorizationRepository with its dependencies
func NewCreateAuthorizationRepository(db *sql.DB) domain.AuthorizationCreator {
	return createAuthorizationRepository{
		db: db,
	}
}

// Create performs insert into the database
func (c createAuthorizationRepository) Create(ctx context.Context, authorization domain.Authorization) (domain.Authorization, error) {
	tx, ok := ctx.Value("TxKey").(*sql.Tx)
	if !ok {
		var err error
		tx, err = c.db.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			return domain.Authorization{}, errors.Wrap(err, errUnknown.Error())
		}
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO authorizations (id, account_id, operation_id, amount, balance, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		authorization.ID(),
		authorization.AccountID(),
		authorization.Operation().ID(),
		authorization.Amount(),
		authorization.CreatedAt(),
	); err != nil {
		return domain.Authorization{}, errors.Wrap(err, errUnknown.Error())
	}

	return authorization, nil
}

func (c createAuthorizationRepository) WithTransaction(ctx context.Context, fn func(ctxFn context.Context) error) error {
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Wrap(err, errUnknown.Error())
	}

	ctxTx := context.WithValue(ctx, "TxKey", tx)
	err = fn(ctxTx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errors.Wrap(err, "rollback error")
		}
		return err
	}

	return tx.Commit()
}

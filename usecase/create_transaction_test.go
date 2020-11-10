package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/GSabadini/go-transactions/domain"
)

type stubCreateTransactionRepo struct {
	result domain.Transaction
	err    error
}

func (s stubCreateTransactionRepo) WithTransaction(_ context.Context, fn func(context.Context) error) error {
	if err := fn(context.Background()); err != nil {
		return err
	}

	return nil
}

func (s stubCreateTransactionRepo) Create(_ context.Context, _ domain.Transaction) (domain.Transaction, error) {
	return s.result, s.err
}

type stubCreateTransactionPresenter struct{}

func (s stubCreateTransactionPresenter) Output(transaction domain.Transaction) CreateTransactionOutput {
	return CreateTransactionOutput{
		ID:        transaction.ID(),
		AccountID: transaction.AccountID(),
		Operation: CreateTransactionOperationOutput{
			ID:          transaction.Operation().ID(),
			Description: transaction.Operation().Description(),
			Type:        transaction.Operation().Type(),
		},
		Amount:    transaction.Amount(),
		Balance:   transaction.Balance(),
		CreatedAt: transaction.CreatedAt().String(),
	}
}

type stubFindUserByRepo struct {
	result domain.Account
	err    error
}

func (s stubFindUserByRepo) FindByID(_ context.Context, _ string) (domain.Account, error) {
	return s.result, s.err
}

type stubUpdateCreditLimitRepo struct {
	err error
}

func (s stubUpdateCreditLimitRepo) UpdateCreditLimit(_ context.Context, _ string, _ int64) error {
	return s.err
}

func Test_createTransactionInteractor_Execute(t *testing.T) {
	var opCompraAVista, _ = domain.NewOperation("fd426041-0648-40f6-9d04-5284295c5095")

	type fields struct {
		repo               domain.TransactionCreator
		repoAccountFinder  domain.AccountFinder
		repoAccountUpdater domain.AccountUpdater
		pre                CreateTransactionPresenter
		ctxTimeout         time.Duration
	}
	type args struct {
		ctx context.Context
		i   CreateTransactionInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    CreateTransactionOutput
		wantErr bool
	}{
		{
			name: "Create successful transaction",
			fields: fields{
				repo: stubCreateTransactionRepo{
					result: domain.NewTransaction(
						"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
						"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
						opCompraAVista,
						10025,
						0,
						time.Time{},
					),
					err: nil,
				},
				repoAccountFinder: stubFindUserByRepo{
					result: domain.NewAccount(
						"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
						"12345678900",
						10025,
						time.Time{},
					),
					err: nil,
				},
				repoAccountUpdater: stubUpdateCreditLimitRepo{err: nil},
				pre:                stubCreateTransactionPresenter{},
				ctxTimeout:         time.Second,
			},
			args: args{
				ctx: context.Background(),
				i: CreateTransactionInput{
					AccountID:   "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					OperationID: "fd426041-0648-40f6-9d04-5284295c5095",
					Amount:      10025,
				},
			},
			want: CreateTransactionOutput{
				ID:        "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				AccountID: "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				Operation: CreateTransactionOperationOutput{
					ID:          domain.CompraAVista,
					Description: "COMPRA A VISTA",
					Type:        domain.Debit,
				},
				Amount:    -10025,
				Balance:   0,
				CreatedAt: time.Time{}.String(),
			},
			wantErr: false,
		},
		{
			name: "Error create transaction insufficient credit limit",
			fields: fields{
				repo: stubCreateTransactionRepo{
					result: domain.NewTransaction(
						"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
						"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
						opCompraAVista,
						10025,
						10025,
						time.Time{},
					),
					err: nil,
				},
				repoAccountFinder: stubFindUserByRepo{
					result: domain.NewAccount(
						"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
						"12345678900",
						1,
						time.Time{},
					),
					err: nil,
				},
				repoAccountUpdater: stubUpdateCreditLimitRepo{err: nil},
				pre:                stubCreateTransactionPresenter{},
				ctxTimeout:         time.Second,
			},
			args: args{
				ctx: context.Background(),
				i: CreateTransactionInput{
					AccountID:   "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					OperationID: "fd426041-0648-40f6-9d04-5284295c5095",
					Amount:      10025,
				},
			},
			want: CreateTransactionOutput{
				CreatedAt: time.Time{}.String(),
			},
			wantErr: true,
		},
		{
			name: "Error creating transaction with invalid operation",
			fields: fields{
				repo: stubCreateTransactionRepo{
					result: domain.Transaction{},
					err:    nil,
				},
				repoAccountFinder: stubFindUserByRepo{
					result: domain.NewAccount(
						"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
						"12345678900",
						1,
						time.Time{},
					),
					err: nil,
				},
				repoAccountUpdater: stubUpdateCreditLimitRepo{err: nil},
				pre:                stubCreateTransactionPresenter{},
				ctxTimeout:         time.Second,
			},
			args: args{
				ctx: context.Background(),
				i: CreateTransactionInput{
					AccountID:   "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					OperationID: "invalid operation",
					Amount:      10025,
				},
			},
			want: CreateTransactionOutput{
				CreatedAt: time.Time{}.String(),
			},
			wantErr: true,
		},
		{
			name: "Repository error when create transaction",
			fields: fields{
				repo: stubCreateTransactionRepo{
					result: domain.Transaction{},
					err:    errors.New("db_error"),
				},
				repoAccountFinder: stubFindUserByRepo{
					result: domain.NewAccount(
						"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
						"12345678900",
						10025,
						time.Time{},
					),
					err: nil,
				},
				repoAccountUpdater: stubUpdateCreditLimitRepo{err: nil},
				pre:                stubCreateTransactionPresenter{},
				ctxTimeout:         time.Second,
			},
			args: args{
				ctx: context.Background(),
				i: CreateTransactionInput{
					AccountID:   "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					OperationID: domain.CompraAVista,
					Amount:      10025,
				},
			},
			want: CreateTransactionOutput{
				CreatedAt: time.Time{}.String(),
			},
			wantErr: true,
		},
		{
			name: "Error find account by id repository",
			fields: fields{
				repo: stubCreateTransactionRepo{
					result: domain.Transaction{},
					err:    nil,
				},
				repoAccountFinder: stubFindUserByRepo{
					result: domain.Account{},
					err:    errors.New("db_error"),
				},
				repoAccountUpdater: stubUpdateCreditLimitRepo{err: nil},
				pre:                stubCreateTransactionPresenter{},
				ctxTimeout:         time.Second,
			},
			args: args{
				ctx: context.Background(),
				i: CreateTransactionInput{
					AccountID:   "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					OperationID: domain.CompraAVista,
					Amount:      10025,
				},
			},
			want: CreateTransactionOutput{
				CreatedAt: time.Time{}.String(),
			},
			wantErr: true,
		},
		{
			name: "Error update credit limit by id repository",
			fields: fields{
				repo: stubCreateTransactionRepo{
					result: domain.Transaction{},
					err:    nil,
				},
				repoAccountFinder: stubFindUserByRepo{
					result: domain.NewAccount(
						"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
						"12345678900",
						10025,
						time.Time{},
					),
					err: nil,
				},
				repoAccountUpdater: stubUpdateCreditLimitRepo{err: errors.New("db_error")},
				pre:                stubCreateTransactionPresenter{},
				ctxTimeout:         time.Second,
			},
			args: args{
				ctx: context.Background(),
				i: CreateTransactionInput{
					AccountID:   "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					OperationID: domain.CompraAVista,
					Amount:      10025,
				},
			},
			want: CreateTransactionOutput{
				CreatedAt: time.Time{}.String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interactor := NewCreateTransactionInteractor(
				tt.fields.repo,
				tt.fields.repoAccountFinder,
				tt.fields.repoAccountUpdater,
				tt.fields.pre,
				tt.fields.ctxTimeout,
			)

			got, err := interactor.Execute(tt.args.ctx, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("[TestCase '%s'] Err: '%v' | WantErr: '%v'", tt.name, err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}

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
		CreatedAt: transaction.CreatedAt().String(),
	}
}

func Test_createTransactionInteractor_Execute(t *testing.T) {
	var opCompraAVista, _ = domain.NewOperation("fd426041-0648-40f6-9d04-5284295c5095")

	type fields struct {
		repo       domain.CreateTransactionRepository
		pre        CreateTransactionPresenter
		ctxTimeout time.Duration
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
						100.25,
						time.Time{},
					),
					err: nil,
				},
				pre:        stubCreateTransactionPresenter{},
				ctxTimeout: time.Second,
			},
			args: args{
				ctx: context.Background(),
				i: CreateTransactionInput{
					AccountID:   "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					OperationID: "fd426041-0648-40f6-9d04-5284295c5095",
					Amount:      100.25,
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
				Amount:    -100.25,
				CreatedAt: time.Time{}.String(),
			},
			wantErr: false,
		},
		{
			name: "Error creating transaction with invalid operation",
			fields: fields{
				repo: stubCreateTransactionRepo{
					result: domain.Transaction{},
					err:    nil,
				},
				pre:        stubCreateTransactionPresenter{},
				ctxTimeout: time.Second,
			},
			args: args{
				ctx: context.Background(),
				i: CreateTransactionInput{
					AccountID:   "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					OperationID: "invalid operation",
					Amount:      100.25,
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
				pre:        stubCreateTransactionPresenter{},
				ctxTimeout: time.Second,
			},
			args: args{
				ctx: context.Background(),
				i: CreateTransactionInput{
					AccountID:   "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					OperationID: domain.CompraAVista,
					Amount:      100.25,
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

package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/GSabadini/go-transactions/domain"
)

type stubCreateAccountRepo struct {
	result domain.Account
	err    error
}

func (s stubCreateAccountRepo) Create(_ context.Context, _ domain.Account) (domain.Account, error) {
	return s.result, s.err
}

type stubCreateAccountPresenter struct {
	result AccountOutput
}

func (s stubCreateAccountPresenter) Output(_ domain.Account) AccountOutput {
	return s.result
}

func Test_createAccountInteractor_Execute(t *testing.T) {
	type fields struct {
		repo       domain.CreateAccountRepository
		pre        CreateAccountPresenter
		ctxTimeout time.Duration
	}
	type args struct {
		ctx   context.Context
		input AccountInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    AccountOutput
		wantErr bool
	}{
		{
			name: "Create account success",
			fields: fields{
				repo: stubCreateAccountRepo{
					result: domain.NewAccount(
						"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
						"12345678900",
						time.Time{},
					),
					err: nil,
				},
				pre: stubCreateAccountPresenter{
					result: AccountOutput{
						ID: "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
						Document: AccountDocumentOutput{
							Number: "12345678900",
						},
						CreatedAt: time.Time{}.String(),
					},
				},
				ctxTimeout: time.Second,
			},
			args: args{
				ctx: context.Background(),
				input: AccountInput{
					Document: struct {
						Number string
					}{
						Number: "12345678900",
					},
				},
			},
			want: AccountOutput{
				ID: "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				Document: AccountDocumentOutput{
					Number: "12345678900",
				},
				CreatedAt: time.Time{}.String(),
			},
			wantErr: false,
		},
		{
			name: "Create account repository error",
			fields: fields{
				repo: stubCreateAccountRepo{
					result: domain.Account{},
					err:    errors.New("db_error"),
				},
				pre:        stubCreateAccountPresenter{},
				ctxTimeout: time.Second,
			},
			args: args{
				ctx: context.Background(),
				input: AccountInput{
					Document: struct {
						Number string
					}{
						Number: "12345678900",
					},
				},
			},
			want:    AccountOutput{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interactor := NewCreateAccountInteractor(tt.fields.repo, tt.fields.pre, tt.fields.ctxTimeout)

			got, err := interactor.Execute(tt.args.ctx, tt.args.input)
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

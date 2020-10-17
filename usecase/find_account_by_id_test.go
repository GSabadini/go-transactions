package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/GSabadini/go-transactions/domain"
)

type stubFindAccountByIDRepo struct {
	result domain.Account
	err    error
}

func (s stubFindAccountByIDRepo) FindByID(_ context.Context, _ string) (domain.Account, error) {
	return s.result, s.err
}

type stubFindAccountByIDPresenter struct {
	result FindAccountByIDOutput
}

func (s stubFindAccountByIDPresenter) Output(_ domain.Account) FindAccountByIDOutput {
	return s.result
}

func Test_findAccountByIDInteractor_Execute(t *testing.T) {
	type fields struct {
		repo       domain.FindAccountByIDRepository
		pre        FindAccountByIDPresenter
		ctxTimeout time.Duration
	}
	type args struct {
		ctx context.Context
		i   FindAccountByIDInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    FindAccountByIDOutput
		wantErr bool
	}{
		{
			name: "Find account by id success",
			fields: fields{
				repo: stubFindAccountByIDRepo{
					result: domain.NewAccount(
						"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
						"12345678900",
						time.Time{},
					),
					err: nil,
				},
				pre: stubFindAccountByIDPresenter{
					result: FindAccountByIDOutput{
						ID: "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
						Document: struct {
							Number string `json:"number"`
						}{
							Number: "12345678900",
						},
						CreatedAt: time.Time{}.String(),
					},
				},
				ctxTimeout: time.Second,
			},
			args: args{
				ctx: context.Background(),
				i: FindAccountByIDInput{
					ID: "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				},
			},
			want: FindAccountByIDOutput{
				ID: "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				Document: struct {
					Number string `json:"number"`
				}{
					Number: "12345678900",
				},
				CreatedAt: time.Time{}.String(),
			},
			wantErr: false,
		},
		{
			name: "Find account by id account not found",
			fields: fields{
				repo: stubFindAccountByIDRepo{
					result: domain.Account{},
					err:    domain.ErrAccountNotFound,
				},
				pre: stubFindAccountByIDPresenter{
					result: FindAccountByIDOutput{},
				},
				ctxTimeout: time.Second,
			},
			args: args{
				ctx: context.Background(),
				i: FindAccountByIDInput{
					ID: "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				},
			},
			want:    FindAccountByIDOutput{},
			wantErr: true,
		},
		{
			name: "Find account by id account repository error",
			fields: fields{
				repo: stubFindAccountByIDRepo{
					result: domain.Account{},
					err:    errors.New("db_error"),
				},
				pre: stubFindAccountByIDPresenter{
					result: FindAccountByIDOutput{},
				},
				ctxTimeout: time.Second,
			},
			args: args{
				ctx: context.Background(),
				i: FindAccountByIDInput{
					ID: "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				},
			},
			want:    FindAccountByIDOutput{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interactor := NewFindAccountByIDInteractor(tt.fields.repo, tt.fields.pre, tt.fields.ctxTimeout)

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

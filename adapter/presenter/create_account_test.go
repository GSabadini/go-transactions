package presenter

import (
	"reflect"
	"testing"
	"time"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/GSabadini/go-transactions/usecase"
)

func Test_createAccountPresenter_Output(t *testing.T) {
	type args struct {
		account domain.Account
	}
	tests := []struct {
		name string
		args args
		want usecase.CreateAccountOutput
	}{
		{
			name: "Create account output",
			args: args{
				account: domain.NewAccount(
					"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					"12345678900",
					100,
					time.Time{},
				),
			},
			want: usecase.CreateAccountOutput{
				ID:                   "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				AvailableCreditLimit: 100,
				Document: usecase.CreateAccountDocumentOutput{
					Number: "12345678900",
				},
				CreatedAt: "0001-01-01T00:00:00Z",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre := NewCreateAccountPresenter()
			if got := pre.Output(tt.args.account); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}

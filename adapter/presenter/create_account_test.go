package presenter

import (
	"github.com/GSabadini/go-transactions/domain"
	"github.com/GSabadini/go-transactions/usecase"
	"reflect"
	"testing"
	"time"
)

func Test_createAccountPresenter_Output(t *testing.T) {
	type args struct {
		account domain.Account
	}
	tests := []struct {
		name string
		args args
		want usecase.AccountOutput
	}{
		{
			name: "Create account output",
			args: args{
				account: domain.NewAccount(
					"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					"12345678900",
					time.Time{},
				),
			},
			want: usecase.AccountOutput{
				ID: "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				Document: usecase.AccountDocumentOutput{
					Number: "12345678900",
				},
				CreatedAt: "0001-01-01T00:00:00Z",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCreateAccountPresenter()
			if got := c.Output(tt.args.account); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}

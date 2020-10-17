package presenter

import (
	"github.com/GSabadini/go-transactions/domain"
	"github.com/GSabadini/go-transactions/usecase"
	"reflect"
	"testing"
	"time"
)

func Test_createTransactionPresenter_Output(t *testing.T) {
	var op, _ = domain.NewOperation("fd426041-0648-40f6-9d04-5284295c5095")

	type args struct {
		transaction domain.Transaction
	}
	tests := []struct {
		name string
		args args
		want usecase.CreateTransactionOutput
	}{
		{
			name: "Create transaction output",
			args: args{
				transaction: domain.NewTransaction(
					"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					"eae0bbf7-19ee-46d6-8244-77bccd64ab93",
					op,
					100.25,
					time.Time{},
				),
			},
			want: usecase.CreateTransactionOutput{
				ID:        "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				AccountID: "eae0bbf7-19ee-46d6-8244-77bccd64ab93",
				Operation: usecase.CreateTransactionOperationOutput{
					ID:          "fd426041-0648-40f6-9d04-5284295c5095",
					Description: "COMPRA A VISTA",
					Type:        "DEBIT",
				},
				Amount:    -100.25,
				CreatedAt: "0001-01-01T00:00:00Z",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pre := NewCreateTransactionPresenter()
			if got := pre.Output(tt.args.transaction); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}

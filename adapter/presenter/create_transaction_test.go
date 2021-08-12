package presenter

import (
	"reflect"
	"testing"
	"time"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/GSabadini/go-transactions/usecase"
)

func Test_createTransactionPresenter_Output(t *testing.T) {
	var (
		opCompraAVista, _ = domain.NewOperation(domain.CompraAVista)
		opPagamento, _    = domain.NewOperation(domain.Pagamento)
	)

	type args struct {
		transaction domain.Transaction
	}
	tests := []struct {
		name string
		args args
		want usecase.CreateTransactionOutput
	}{
		{
			name: "Create transaction operation type compra a vista output",
			args: args{
				transaction: domain.NewTransaction(
					"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					"eae0bbf7-19ee-46d6-8244-77bccd64ab93",
					opCompraAVista,
					10025,
					10025,
					time.Time{},
				),
			},
			want: usecase.CreateTransactionOutput{
				ID:        "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				AccountID: "eae0bbf7-19ee-46d6-8244-77bccd64ab93",
				Operation: usecase.CreateTransactionOperationOutput{
					ID:          "1",
					Description: "COMPRA A VISTA",
					Type:        "DEBIT",
				},
				Amount:    -10025,
				Balance:   10025,
				CreatedAt: "0001-01-01T00:00:00Z",
			},
		},
		{
			name: "Create transaction operation type pagamento output",
			args: args{
				transaction: domain.NewTransaction(
					"fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
					"eae0bbf7-19ee-46d6-8244-77bccd64ab93",
					opPagamento,
					10025,
					10025,
					time.Time{},
				),
			},
			want: usecase.CreateTransactionOutput{
				ID:        "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				AccountID: "eae0bbf7-19ee-46d6-8244-77bccd64ab93",
				Operation: usecase.CreateTransactionOperationOutput{
					ID:          "4",
					Description: "PAGAMENTO",
					Type:        "CREDIT",
				},
				Amount:    10025,
				Balance:   10025,
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

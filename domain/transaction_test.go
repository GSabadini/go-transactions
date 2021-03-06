package domain

import (
	"testing"
	"time"
)

func TestTransaction_Amount(t *testing.T) {
	type fields struct {
		id        string
		accountID string
		operation Operation
		amount    int64
		balance   int64
		createdAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "Return amount operation type Debit",
			fields: fields{
				id:        "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				accountID: "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				operation: Operation{
					id:          CompraAVista,
					description: "COMPRA A VISTA",
					opType:      Debit,
				},
				amount:    10024,
				createdAt: time.Time{},
			},
			want: -10024,
		},
		{
			name: "Return amount operation type Credit",
			fields: fields{
				id:        "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				accountID: "fc95e907-e0eb-4ef8-927e-3eaad3a4d9a8",
				operation: Operation{
					id:          Pagamento,
					description: "PAGAMENTO",
					opType:      Credit,
				},
				amount:    10024,
				createdAt: time.Time{},
			},
			want: 10024,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			transaction := NewTransaction(
				tt.fields.id,
				tt.fields.accountID,
				tt.fields.operation,
				tt.fields.amount,
				tt.fields.balance,
				tt.fields.createdAt,
			)
			if got := transaction.Amount(); got != tt.want {
				t.Errorf("[TestCase '%s'] Got: '%+v' | Want: '%+v'", tt.name, got, tt.want)
			}
		})
	}
}

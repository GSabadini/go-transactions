package domain

import (
	"reflect"
	"testing"
)

func TestNewOperation(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    Operation
		wantErr bool
	}{
		{
			name: "Create operation compra a vista",
			args: args{
				id: "1",
			},
			want: Operation{
				id:          CompraAVista,
				description: "COMPRA A VISTA",
				opType:      Debit,
			},
			wantErr: false,
		},
		{
			name: "Create operation compra parcelada",
			args: args{
				id: "2",
			},
			want: Operation{
				id:          CompraParcelada,
				description: "COMPRA PARCELADA",
				opType:      Debit,
			},
			wantErr: false,
		},
		{
			name: "Create operation Saque",
			args: args{
				id: "3",
			},
			want: Operation{
				id:          Saque,
				description: "SAQUE",
				opType:      Debit,
			},
			wantErr: false,
		},
		{
			name: "Create operation Pagamento",
			args: args{
				id: "4",
			},
			want: Operation{
				id:          Pagamento,
				description: "PAGAMENTO",
				opType:      Credit,
			},
			wantErr: false,
		},
		{
			name: "Error operation type invalid",
			args: args{
				id: "123",
			},
			want:    Operation{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOperation(tt.args.id)
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

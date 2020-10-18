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
				id: "fd426041-0648-40f6-9d04-5284295c5095",
			},
			want: Operation{
				id:          CompraAVista,
				description: "COMPRA A VISTA",
				opType:      debit,
			},
			wantErr: false,
		},
		{
			name: "Create operation compra parcelada",
			args: args{
				id: "b03dcb59-006f-472f-a8f1-58651990dea6",
			},
			want: Operation{
				id:          CompraParcelada,
				description: "COMPRA PARCELADA",
				opType:      debit,
			},
			wantErr: false,
		},
		{
			name: "Create operation Saque",
			args: args{
				id: "3f973e5b-cb9f-475c-b27d-8f855a0b90b0",
			},
			want: Operation{
				id:          Saque,
				description: "SAQUE",
				opType:      debit,
			},
			wantErr: false,
		},
		{
			name: "Create operation Pagamento",
			args: args{
				id: "976f88ea-eb2f-4325-a106-26f9cb35810d",
			},
			want: Operation{
				id:          Pagamento,
				description: "PAGAMENTO",
				opType:      credit,
			},
			wantErr: false,
		},
		{
			name: "Create invalid operation",
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

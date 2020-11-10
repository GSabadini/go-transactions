package domain

import (
	"testing"
	"time"
)

func TestAccount_Deposit(t *testing.T) {
	type fields struct {
		id                   string
		document             Document
		availableCreditLimit int64
		createdAt            time.Time
	}
	type args struct {
		amount int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name: "Successful depositing amount",
			fields: fields{
				id: "123",
				document: Document{
					number: "123",
				},
				availableCreditLimit: 100,
				createdAt:            time.Time{},
			},
			args: args{
				amount: 50,
			},
			want: 150,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := NewAccount(
				tt.fields.id,
				tt.fields.document.number,
				tt.fields.availableCreditLimit,
				tt.fields.createdAt,
			)

			account.Deposit(tt.args.amount)
			if account.AvailableCreditLimit() != tt.want {
				t.Errorf("[TestCase '%s'] Result: '%v' | Expected: '%v'",
					tt.name,
					account.availableCreditLimit,
					tt.want,
				)
			}
		})
	}
}

func TestAccount_PaymentOperation(t *testing.T) {
	type fields struct {
		id                   string
		document             Document
		availableCreditLimit int64
		createdAt            time.Time
	}
	type args struct {
		amount int64
		opType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "Operation type debit",
			fields: fields{
				id: "123",
				document: Document{
					number: "123",
				},
				availableCreditLimit: 100,
				createdAt:            time.Time{},
			},
			args: args{
				amount: 100,
				opType: Debit,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "Operation type credit",
			fields: fields{
				id: "123",
				document: Document{
					number: "123",
				},
				availableCreditLimit: 0,
				createdAt:            time.Time{},
			},
			args: args{
				amount: 100,
				opType: Credit,
			},
			want:    100,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := NewAccount(
				tt.fields.id,
				tt.fields.document.number,
				tt.fields.availableCreditLimit,
				tt.fields.createdAt,
			)
			if err := account.PaymentOperation(tt.args.amount, tt.args.opType); (err != nil) != tt.wantErr {
				t.Errorf("[TestCase '%s'] Error: '%v' | WantErr: '%v'",
					tt.name,
					err,
					tt.wantErr,
				)
				return
			}

			if account.AvailableCreditLimit() != tt.want {
				t.Errorf("[TestCase '%s'] Got: '%v' | Want: '%v'",
					tt.name,
					account.availableCreditLimit,
					tt.want,
				)
			}
		})
	}
}

func TestAccount_Withdraw(t *testing.T) {
	type fields struct {
		id                   string
		document             Document
		availableCreditLimit int64
		createdAt            time.Time
	}
	type args struct {
		amount int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "Success in withdrawing balance",
			fields: fields{
				id: "123",
				document: Document{
					number: "123",
				},
				availableCreditLimit: 100,
				createdAt:            time.Time{},
			},
			args: args{
				amount: 30,
			},
			want:    70,
			wantErr: false,
		},
		{
			name: "Error when withdrawing account without sufficient credit limit",
			fields: fields{
				id: "123",
				document: Document{
					number: "123",
				},
				availableCreditLimit: 100,
				createdAt:            time.Time{},
			},
			args: args{
				amount: 150,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := NewAccount(
				tt.fields.id,
				tt.fields.document.number,
				tt.fields.availableCreditLimit,
				tt.fields.createdAt,
			)
			if err := account.Withdraw(tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("[TestCase '%s'] Got: '%v' | WantErr: '%v'",
					tt.name,
					err,
					tt.wantErr,
				)
				return
			}

			if !tt.wantErr && account.AvailableCreditLimit() != tt.want {
				t.Errorf("[TestCase '%s'] Got: '%v' | Want: '%v'",
					tt.name,
					account.availableCreditLimit,
					tt.want,
				)
			}
		})
	}
}

package handler

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GSabadini/go-transactions/domain"
	"github.com/GSabadini/go-transactions/infrastructure/logger"
	"github.com/GSabadini/go-transactions/infrastructure/validation"
	"github.com/GSabadini/go-transactions/usecase"
	"github.com/go-playground/validator/v10"
)

type stubCreateTransactionUseCase struct {
	result usecase.CreateTransactionOutput
	err    error
}

func (s stubCreateTransactionUseCase) Execute(_ context.Context, _ usecase.CreateTransactionInput) (usecase.CreateTransactionOutput, error) {
	return s.result, s.err
}

func TestCreateTransactionHandler_Handle(t *testing.T) {
	logFake := logger.NewLogFake()
	v := validation.NewValidator()

	type fields struct {
		uc        usecase.CreateTransactionUseCase
		log       *log.Logger
		validator *validator.Validate
	}
	tests := []struct {
		name           string
		fields         fields
		rawPayload     []byte
		wantBody       string
		wantStatusCode int
	}{
		{
			name: "Create transaction successfully",
			fields: fields{
				uc: stubCreateTransactionUseCase{
					result: usecase.CreateTransactionOutput{
						ID:        "aef3836b-5ea4-4890-80ad-e13337ccf47f",
						AccountID: "92c82203-cdba-4932-9860-bce2e6140267",
						Operation: usecase.CreateTransactionOperationOutput{
							ID:          domain.CompraAVista,
							Description: "COMPRA A VISTA",
							Type:        domain.Debit,
						},
						Amount:    -1074,
						CreatedAt: "2020-10-16T17:50:39Z",
					},
					err: nil,
				},
				log:       logFake,
				validator: v,
			},
			rawPayload:     []byte(`{"account_id": "92c82203-cdba-4932-9860-bce2e6140267","operation_id": "1","amount": 1074}`),
			wantBody:       `{"id":"aef3836b-5ea4-4890-80ad-e13337ccf47f","account_id":"92c82203-cdba-4932-9860-bce2e6140267","operation":{"id":"1","description":"COMPRA A VISTA","type":"DEBIT"},"amount":-1074,"balance":0,"created_at":"2020-10-16T17:50:39Z"}`,
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "Error operation type invalid",
			fields: fields{
				uc: stubCreateTransactionUseCase{
					result: usecase.CreateTransactionOutput{},
					err:    domain.ErrOperationInvalid,
				},
				log:       logFake,
				validator: v,
			},
			rawPayload:     []byte(`{"account_id": "92c82203-cdba-4932-9860-bce2e6140267","operation_id": "invalid","amount": 1074}`),
			wantBody:       `{"errors":["operation type invalid"]}`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Error required fields",
			fields: fields{
				uc: stubCreateTransactionUseCase{
					result: usecase.CreateTransactionOutput{},
					err:    nil,
				},
				log:       logFake,
				validator: v,
			},
			rawPayload:     []byte(`{}`),
			wantBody:       `{"errors":["account_id is a required field","operation_id is a required field","amount is a required field"]}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Error invalid amount",
			fields: fields{
				uc: stubCreateTransactionUseCase{
					result: usecase.CreateTransactionOutput{},
					err:    nil,
				},
				log:       logFake,
				validator: v,
			},
			rawPayload:     []byte(`{"account_id": "92c82203-cdba-4932-9860-bce2e6140267","operation_id": "fd426041-0648-40f6-9d04-5284295c509","amount": -1074}`),
			wantBody:       `{"errors":["amount must be greater than 0"]}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Repository error when create account",
			fields: fields{
				uc: stubCreateTransactionUseCase{
					result: usecase.CreateTransactionOutput{},
					err:    errors.New("db_error"),
				},
				log:       logFake,
				validator: v,
			},
			rawPayload:     []byte(`{"account_id": "92c82203-cdba-4932-9860-bce2e6140267","operation_id": "fd426041-0648-40f6-9d04-5284295c509","amount": 1074}`),
			wantBody:       `{"errors":["db_error"]}`,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Error credit limit insufficient",
			fields: fields{
				uc: stubCreateTransactionUseCase{
					result: usecase.CreateTransactionOutput{},
					err:    domain.ErrAccountInsufficientCreditLimit,
				},
				log:       logFake,
				validator: v,
			},
			rawPayload:     []byte(`{"account_id": "92c82203-cdba-4932-9860-bce2e6140267","operation_id": "1","amount": 1074}`),
			wantBody:       `{"errors":["credit limit insufficient"]}`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(
				http.MethodPost,
				"/transactions",
				bytes.NewReader(tt.rawPayload),
			)
			if err != nil {
				t.Fatal(err)
			}

			var (
				w       = httptest.NewRecorder()
				handler = NewCreateTransactionHandler(tt.fields.uc, tt.fields.log, tt.fields.validator)
			)

			handler.Handle(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf(
					"[TestCase '%s'] Got status code: '%v' | Want status code: '%v'",
					tt.name,
					w.Code,
					tt.wantStatusCode,
				)
			}

			var got = strings.TrimSpace(w.Body.String())
			if !strings.EqualFold(got, tt.wantBody) {
				t.Errorf(
					"[TestCase '%s'] Got body: '%v' |\n Want body: '%v'",
					tt.name,
					got,
					tt.wantBody,
				)
			}
		})
	}
}

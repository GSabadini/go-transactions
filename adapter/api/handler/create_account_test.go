package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/GSabadini/go-transactions/domain"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GSabadini/go-transactions/infrastructure/logger"
	"github.com/GSabadini/go-transactions/infrastructure/validation"
	"github.com/GSabadini/go-transactions/usecase"
	"github.com/go-playground/validator/v10"
)

type stubCreateAccountUseCase struct {
	result usecase.CreateAccountOutput
	err    error
}

func (s stubCreateAccountUseCase) Execute(_ context.Context, _ usecase.CreateAccountInput) (usecase.CreateAccountOutput, error) {
	return s.result, s.err
}

func TestCreateAccountHandler_Handle(t *testing.T) {
	logFake := logger.NewLogFake()
	v := validation.NewValidator()

	type fields struct {
		uc        usecase.CreateAccountUseCase
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
			name: "Create account successfully",
			fields: fields{
				uc: stubCreateAccountUseCase{
					result: usecase.CreateAccountOutput{
						ID:                   "cfd3c0e0-cfa7-4220-8e62-069657874aba",
						AvailableCreditLimit: 100,
						Document: usecase.CreateAccountDocumentOutput{
							Number: "12345678900",
						},
						CreatedAt: "2020-10-16T17:50:39Z",
					},
					err: nil,
				},
				log:       logFake,
				validator: v,
			},
			rawPayload:     []byte(`{"document": {"number": "12345678900"}, "available_credit_limit": 100}`),
			wantBody:       `{"id":"cfd3c0e0-cfa7-4220-8e62-069657874aba","available_credit_limit":100,"document":{"number":"12345678900"},"created_at":"2020-10-16T17:50:39Z"}`,
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "Error creating existing account",
			fields: fields{
				uc: stubCreateAccountUseCase{
					result: usecase.CreateAccountOutput{},
					err:    domain.ErrAccountAlreadyExists,
				},
				log:       logFake,
				validator: v,
			},
			rawPayload:     []byte(`{"document": {"number": "12345678900"}, "available_credit_limit": 100}`),
			wantBody:       `{"errors":["account already exists"]}`,
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Error required field",
			fields: fields{
				uc: stubCreateAccountUseCase{
					result: usecase.CreateAccountOutput{},
					err:    nil,
				},
				log:       logFake,
				validator: v,
			},
			rawPayload:     []byte(`{"document": {}, "available_credit_limit": 100}`),
			wantBody:       `{"errors":["number is a required field"]}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Error required field",
			fields: fields{
				uc: stubCreateAccountUseCase{
					result: usecase.CreateAccountOutput{},
					err:    nil,
				},
				log:       logFake,
				validator: v,
			},
			rawPayload:     []byte(`{"document": {"number": "123456"}}`),
			wantBody:       `{"errors":["available_credit_limit is a required field"]}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Error required field",
			fields: fields{
				uc: stubCreateAccountUseCase{
					result: usecase.CreateAccountOutput{},
					err:    nil,
				},
				log:       logFake,
				validator: v,
			},
			rawPayload:     []byte(`{"document": {"number": "123456"}, "available_credit_limit": -100}`),
			wantBody:       `{"errors":["available_credit_limit must be greater than 0"]}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Error exceeded the number of characters field",
			fields: fields{
				uc: stubCreateAccountUseCase{
					result: usecase.CreateAccountOutput{},
					err:    nil,
				},
				log:       logFake,
				validator: v,
			},
			rawPayload:     []byte(`{"document": {"number": "1234567899876545646455432103215648721212156451546456451205554564564564564"}, "available_credit_limit": 100}`),
			wantBody:       `{"errors":["number must be a maximum of 30 characters in length"]}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Repository error when create account",
			fields: fields{
				uc: stubCreateAccountUseCase{
					result: usecase.CreateAccountOutput{},
					err:    errors.New("db_error"),
				},
				log:       logFake,
				validator: v,
			},
			rawPayload:     []byte(`{"document": {"number": "12345678900"}, "available_credit_limit": 100}`),
			wantBody:       `{"errors":["db_error"]}`,
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(
				http.MethodPost,
				"/accounts",
				bytes.NewReader(tt.rawPayload),
			)
			if err != nil {
				t.Fatal(err)
			}

			var (
				w       = httptest.NewRecorder()
				handler = NewCreateAccountHandler(tt.fields.uc, tt.fields.log, tt.fields.validator)
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
					"[TestCase '%s'] Got body: '%v' | Want body: '%v'",
					tt.name,
					got,
					tt.wantBody,
				)
			}
		})
	}
}

package handler

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/GSabadini/go-transactions/usecase"
)

type stubCreateAccountUseCase struct {
	result usecase.AccountOutput
	err    error
}

func (s stubCreateAccountUseCase) Execute(_ context.Context, _ usecase.AccountInput) (usecase.AccountOutput, error) {
	return s.result, s.err
}

func TestCreateAccountHandler_Handle(t *testing.T) {
	logDummy := log.New(os.Stdout, "", log.LstdFlags)

	type fields struct {
		uc  usecase.CreateAccountUseCase
		log *log.Logger
	}

	tests := []struct {
		name           string
		fields         fields
		rawPayload     []byte
		wantBody       string
		wantStatusCode int
	}{
		{
			name: "Create account success",
			fields: fields{
				uc: stubCreateAccountUseCase{
					result: usecase.AccountOutput{
						ID: "cfd3c0e0-cfa7-4220-8e62-069657874aba",
						Document: usecase.AccountDocumentOutput{
							Number: "12345678900",
						},
						CreatedAt: "2020-10-16T17:50:39Z",
					},
					err: nil,
				},
				log: logDummy,
			},
			rawPayload:     []byte(`{"document": {"number": "12345678900"}}`),
			wantBody:       `{"id":"cfd3c0e0-cfa7-4220-8e62-069657874aba","document":{"number":"12345678900"},"created_at":"2020-10-16T17:50:39Z"}`,
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "Create account repository error",
			fields: fields{
				uc: stubCreateAccountUseCase{
					result: usecase.AccountOutput{},
					err:    errors.New("db_error"),
				},
				log: logDummy,
			},
			rawPayload:     []byte(`{"document": {"number": "12345678900"}}`),
			wantBody:       `{"errors":["db_error"]}`,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Create account failed to marshal",
			fields: fields{
				uc: stubCreateAccountUseCase{
					result: usecase.AccountOutput{},
					err:    nil,
				},
				log: logDummy,
			},
			rawPayload:     []byte(`{"document":`),
			wantBody:       `{"errors":["unexpected EOF"]}`,
			wantStatusCode: http.StatusBadRequest,
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
				handler = NewCreateAccountHandler(tt.fields.uc, tt.fields.log)
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

package handler

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GSabadini/go-transactions/usecase"
	"github.com/go-playground/validator/v10"
)

func TestCreateTransactionHandler_Handle(t *testing.T) {
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
		{},
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
					"[TestCase '%s'] Got body: '%v' | Want body: '%v'",
					tt.name,
					got,
					tt.wantBody,
				)
			}
		})
	}
}

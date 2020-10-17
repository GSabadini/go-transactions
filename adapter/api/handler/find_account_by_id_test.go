package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/GSabadini/go-transactions/usecase"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type stubFindAccountByIDUseCase struct {
	result usecase.FindAccountByIDOutput
	err    error
}

func (s stubFindAccountByIDUseCase) Execute(_ context.Context, _ usecase.FindAccountByIDInput) (usecase.FindAccountByIDOutput, error) {
	return s.result, s.err
}

func TestFindAccountByIDHandler_Handle(t *testing.T) {
	logDummy := log.New(fakeWriter{}, "", log.LstdFlags)

	type fields struct {
		uc  usecase.FindAccountByIDUseCase
		log *log.Logger
	}
	type args struct {
		ID string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantBody       string
		wantStatusCode int
	}{
		{
			name: "Find account by id success",
			fields: fields{
				uc: stubFindAccountByIDUseCase{
					result: usecase.FindAccountByIDOutput{
						ID: "cfd3c0e0-cfa7-4220-8e62-069657874aba",
						Document: usecase.FindAccountByIDDocumentOutput{
							Number: "123456789000",
						},
						CreatedAt: time.Time{}.String(),
					},
					err: nil,
				},
				log: logDummy,
			},
			args: args{
				ID: "cfd3c0e0-cfa7-4220-8e62-069657874aba",
			},
			wantBody:       `{"id":"cfd3c0e0-cfa7-4220-8e62-069657874aba","document":{"number":"123456789000"},"created_at":"0001-01-01 00:00:00 +0000 UTC"}`,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Find account by id use case error",
			fields: fields{
				uc: stubFindAccountByIDUseCase{
					result: usecase.FindAccountByIDOutput{},
					err:    errors.New("use case error"),
				},
				log: logDummy,
			},
			args: args{
				ID: "cfd3c0e0-cfa7-4220-8e62-069657874aba",
			},
			wantBody:       `{"errors":{"error":"use case error"}}`,
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Find account by id invalid account id",
			fields: fields{
				uc: stubFindAccountByIDUseCase{
					result: usecase.FindAccountByIDOutput{},
					err:    nil,
				},
				log: logDummy,
			},
			args: args{
				ID: "",
			},
			wantBody:       `{"errors":{"error":"invalid account id"}}`,
			wantStatusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri := fmt.Sprintf("/accounts/%s", tt.args.ID)
			req, _ := http.NewRequest(http.MethodGet, uri, nil)
			req = mux.SetURLVars(req, map[string]string{"account_id": tt.args.ID})

			var (
				w       = httptest.NewRecorder()
				handler = NewFindAccountByIDHandler(tt.fields.uc, tt.fields.log)
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

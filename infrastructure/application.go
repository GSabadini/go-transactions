package infrastructure

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GSabadini/go-transactions/adapter/api/handler"
	"github.com/GSabadini/go-transactions/adapter/presenter"
	"github.com/GSabadini/go-transactions/adapter/repository"
	"github.com/GSabadini/go-transactions/infrastructure/database"
	"github.com/GSabadini/go-transactions/infrastructure/logger"
	"github.com/GSabadini/go-transactions/infrastructure/router"
	"github.com/GSabadini/go-transactions/infrastructure/validation"
	"github.com/GSabadini/go-transactions/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Application define an application structure
type Application struct {
	database  *sql.DB
	logger    *log.Logger
	router    *mux.Router
	validator *validator.Validate
}

// NewApplication creates new Application with its dependencies
func NewApplication() *Application {
	return &Application{
		database:  database.NewMySQLConnection(),
		logger:    logger.NewLog(),
		router:    router.NewGorillaMux(),
		validator: validation.NewValidator(),
	}
}

// Start run the application
func (a Application) Start() {
	api := a.router.PathPrefix("/v1").Subrouter()

	api.Handle("/accounts", a.createAccountHandler()).Methods(http.MethodPost)
	api.Handle("/accounts/{account_id}", a.findAccountByIDHandler()).Methods(http.MethodGet)

	api.Handle("/transactions", a.createTransactionHandler()).Methods(http.MethodPost)

	api.HandleFunc("/health", healthCheck).Methods(http.MethodGet)

	server := &http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler:      a.router,
	}

	a.logger.Println("Starting HTTP Server in port:", os.Getenv("APP_PORT"))
	a.logger.Fatal(server.ListenAndServe())
}

func (a Application) createAccountHandler() http.HandlerFunc {
	uc := usecase.NewCreateAccountInteractor(
		repository.NewCreateAccountRepository(a.database),
		presenter.NewCreateAccountPresenter(),
		5*time.Second,
	)

	return handler.NewCreateAccountHandler(uc, a.logger, a.validator).Handle
}

func (a Application) findAccountByIDHandler() http.HandlerFunc {
	uc := usecase.NewFindAccountByIDInteractor(
		repository.NewAccountByIDRepository(a.database),
		presenter.NewFindAccountByIDPresenter(),
		5*time.Second,
	)

	return handler.NewFindAccountByIDHandler(uc, a.logger).Handle
}

func (a Application) createTransactionHandler() http.HandlerFunc {
	uc := usecase.NewCreateTransactionInteractor(
		repository.NewCreateTransactionRepository(a.database),
		presenter.NewCreateTransactionPresenter(),
		5*time.Second,
	)

	return handler.NewCreateTransactionHandler(uc, a.logger, a.validator).Handle
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{Status: http.StatusText(http.StatusOK)})
}

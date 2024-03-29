package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/GSabadini/go-transactions/adapter/api/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

// HTTPServer define an application structure
type HTTPServer struct {
	database  *sql.DB
	logger    *log.Logger
	router    *mux.Router
	validator *validator.Validate
}

// NewHTTPServer creates new HTTPServer with its dependencies
func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		database:  database.NewMySQLConnection(),
		logger:    logger.NewLog(),
		router:    router.NewGorillaMux(),
		validator: validation.NewValidator(),
	}
}

// Start run the application
func (a HTTPServer) Start() {
	api := a.router.PathPrefix("/v1").Subrouter()

	api.Use(middleware.NewCorrelationID().Execute)

	api.Handle("/accounts", a.createAccountHandler()).Methods(http.MethodPost)
	api.Handle("/accounts/{account_id}", a.findAccountByIDHandler()).Methods(http.MethodGet)

	api.Handle("/transactions", a.createTransactionHandler()).Methods(http.MethodPost)

	//api.Handle("/cashout", a.createCashoutHandler()).Methods(http.MethodPost)
	//api.Handle("/cashin", a.createTransactionHandler()).Methods(http.MethodPost)
	//api.Handle("/peer-too-peer", a.createTransactionHandler()).Methods(http.MethodPost)

	api.HandleFunc("/health", healthCheck).Methods(http.MethodGet)

	server := &http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler:      a.router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		a.logger.Println("Starting HTTP Server in port:", os.Getenv("APP_PORT"))
		a.logger.Fatal(server.ListenAndServe())
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		a.logger.Fatal("Server Shutdown Failed")
	}

	a.logger.Println("Service down")
}

func (a HTTPServer) createAccountHandler() http.HandlerFunc {
	uc := usecase.NewCreateAccountInteractor(
		repository.NewCreateAccountRepository(a.database),
		presenter.NewCreateAccountPresenter(),
		5*time.Second,
	)

	return handler.NewCreateAccountHandler(uc, a.logger, a.validator).Handle
}

func (a HTTPServer) findAccountByIDHandler() http.HandlerFunc {
	uc := usecase.NewFindAccountByIDInteractor(
		repository.NewAccountByIDRepository(a.database),
		presenter.NewFindAccountByIDPresenter(),
		5*time.Second,
	)

	return handler.NewFindAccountByIDHandler(uc, a.logger).Handle
}

func (a HTTPServer) createTransactionHandler() http.HandlerFunc {
	uc := usecase.NewCreateTransactionInteractor(
		repository.NewCreateTransactionRepository(a.database),
		repository.NewAccountByIDRepository(a.database),
		repository.NewUpdateAccountCreditLimitRepository(a.database),
		presenter.NewCreateTransactionPresenter(),
		5*time.Second,
	)

	return handler.NewCreateTransactionHandler(uc, a.logger, a.validator).Handle
}

//func (a HTTPServer) createCashoutHandler() http.HandlerFunc {
//	uc := usecase.NewCreateAuthorizationInteractor(
//		repository.NewCreateAuthorizationRepository(a.database),
//		repository.NewAccountByIDRepository(a.database),
//		repository.NewUpdateAccountCreditLimitRepository(a.database),
//		presenter.NewCreateCashoutPresenter(),
//		10*time.Second,
//	)
//
//	return handler.NewCreateTransactionHandler(uc, a.logger, a.validator).Handle
//}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{Status: http.StatusText(http.StatusOK)})
}

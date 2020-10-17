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
	"github.com/GSabadini/go-transactions/adapter/api/middleware"
	"github.com/GSabadini/go-transactions/adapter/presenter"
	"github.com/GSabadini/go-transactions/adapter/repository"
	"github.com/GSabadini/go-transactions/infrastructure/database"
	"github.com/GSabadini/go-transactions/usecase"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Application struct {
	database   *sql.DB
	logger     *log.Logger
	router     *mux.Router
	middleware *negroni.Negroni
}

func NewApplication() *Application {
	return &Application{
		database:   database.NewMySQLConnection(),
		logger:     log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds),
		router:     mux.NewRouter(),
		middleware: negroni.New(),
	}
}

func (a Application) Start() {
	api := a.router.PathPrefix("/v1").Subrouter()

	api.Handle("/accounts", a.buildCreateAccountHandler()).Methods(http.MethodPost)
	api.HandleFunc("/health", HealthCheck).Methods(http.MethodGet)

	a.middleware.UseHandler(a.router)

	server := &http.Server{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler:      a.middleware,
	}

	a.logger.Println("Starting HTTP Server in port:", os.Getenv("APP_PORT"))
	a.logger.Fatal(server.ListenAndServe())
}

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{Status: http.StatusText(http.StatusOK)})
}

func (a Application) buildCreateAccountHandler() *negroni.Negroni {
	var handlerFn http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
		var (
			uc = usecase.NewCreateAccountInteractor(
				repository.NewCreateAccountRepository(a.database),
				presenter.NewCreateAccountPresenter(),
				5*time.Second,
			)
			h = handler.NewCreateAccountHandler(uc, a.logger)
		)

		h.Handle(res, req)
	}

	return negroni.New(
		negroni.HandlerFunc(middleware.NewCorrelationID().Execute),
		negroni.NewRecovery(),
		negroni.Wrap(handlerFn),
	)
}

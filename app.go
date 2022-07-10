package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"hermes/db"
	"net/http"
	"os"
)

type App struct {
	Router *mux.Router
}

func (app *App) Initialize(dbConfig db.Config) {
	err := db.InitialiseConnection(dbConfig)
	if err != nil {
		log.Err(err).
			Str("Host", dbConfig.Host).
			Int32("Port", dbConfig.Port).
			Str("User", dbConfig.User).
			Str("Catalog", dbConfig.Catalog).
			Msg("Could not initialise a connection to the database")

		os.Exit(1)
	}

	app.Router = mux.NewRouter()
	app.Router.Use(attachRequestIdentifier)

	app.initializeRoutes()
}

func (app *App) Run(addr string) {
	_ = http.ListenAndServe(addr, app.Router)
}

func attachRequestIdentifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), "requestIdentifier", uuid.New().String()))
		next.ServeHTTP(w, r)
	})
}

func errorResponse(w http.ResponseWriter, code int, message string) {
	jsonResponse(w, code, map[string]string{"error": message})
}

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *App) getCustomers(w http.ResponseWriter, r *http.Request) {
	log.Debug().
		Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
		Str("handler", "getCustomers")

	dbConn, err := db.GetConnection()

	if err != nil {
		log.Err(err).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Could not retrieve connection to the database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	var customers []db.Customer
	if dbResult := dbConn.Find(&customers); dbResult.Error != nil {
		log.Err(err).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Error while fetching from database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	jsonResponse(w, http.StatusOK, customers)
}

func (app *App) getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := uuid.Parse(vars["id"])
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	log.Debug().
		Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
		Str("customerID", customerID.String()).
		Str("handler", "getCustomer")

	dbConn, err := db.GetConnection()

	if err != nil {
		log.Err(err).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Could not retrieve connection to the database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	var customer db.Customer
	if dbResult := dbConn.Where("ID = ?", customerID.String()).First(&customer); dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			jsonResponse(w, http.StatusNotFound, struct{}{})
			return
		}

		log.Err(err).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Error while fetching from database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	jsonResponse(w, http.StatusOK, customer)
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/customers", app.getCustomers).Methods("GET")
	app.Router.HandleFunc("/customers/{id}", app.getCustomer).Methods("GET")
}

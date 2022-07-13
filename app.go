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
	"time"
)

type App struct {
	Router *mux.Router
}

type customerCreateUpdateStruct struct {
	Name string `json:"name"`
}

type customerNotificationChannelCreateUpdateStruct struct {
	ContactCustomer bool   `json:"contact_customer"`
	LookupKey       string `json:"lookup_key"`
}

type notification struct {
	From string `json:"from"`
	Text string `json:"text"`
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Write(response)
	}
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
		log.Err(dbResult.Error).
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

		log.Err(dbResult.Error).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Error while fetching from database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	jsonResponse(w, http.StatusOK, customer)
}

func (app *App) addCustomer(w http.ResponseWriter, r *http.Request) {
	var customer customerCreateUpdateStruct
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Could not unmarshal body")
		return
	}

	log.Debug().
		Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
		Str("customerName", customer.Name).
		Str("handler", "addCustomer")

	dbConn, err := db.GetConnection()

	if err != nil {
		log.Err(err).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Could not retrieve connection to the database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	var newCustomer = db.Customer{ID: uuid.New(), Name: customer.Name}
	if dbResult := dbConn.Create(&newCustomer); dbResult.Error != nil {
		log.Err(dbResult.Error).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Error while saving to database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	jsonResponse(w, http.StatusCreated, newCustomer)
}

func (app *App) updateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := uuid.Parse(vars["id"])
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var customer customerCreateUpdateStruct
	err = json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Could not unmarshal body")
		return
	}

	log.Debug().
		Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
		Str("customerName", customer.Name).
		Str("handler", "updateCustomer")

	dbConn, err := db.GetConnection()

	if err != nil {
		log.Err(err).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Could not retrieve connection to the database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	if dbResult := dbConn.Model(&db.Customer{}).Where("ID = ?", customerID.String()).Update("name", customer.Name); dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			jsonResponse(w, http.StatusNotFound, struct{}{})
			return
		}

		log.Err(dbResult.Error).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Error while saving to database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	jsonResponse(w, http.StatusOK, nil)
}

func (app *App) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := uuid.Parse(vars["id"])
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	log.Debug().
		Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
		Str("customerID", customerID.String()).
		Str("handler", "deleteCustomer")

	dbConn, err := db.GetConnection()

	if err != nil {
		log.Err(err).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Could not retrieve connection to the database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	if dbResult := dbConn.Where("customer_id = ?", customerID.String()).Delete(&db.CustomerNotificationChannels{}); dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			jsonResponse(w, http.StatusNotFound, struct{}{})
			return
		}

		log.Err(dbResult.Error).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Error when deleting from database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	if dbResult := dbConn.Where("ID = ?", customerID.String()).Delete(&db.Customer{}); dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			jsonResponse(w, http.StatusNotFound, struct{}{})
			return
		}

		log.Err(dbResult.Error).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Error when deleting from database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	jsonResponse(w, http.StatusOK, nil)
}

func (app *App) getCustomerNotificationChannels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerID, err := uuid.Parse(vars["id"])
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	log.Debug().
		Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
		Str("customerID", customerID.String()).
		Str("handler", "getCustomerNotificationChannels")

	dbConn, err := db.GetConnection()

	if err != nil {
		log.Err(err).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Could not retrieve connection to the database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	var customerNotificationChannels []db.CustomerNotificationChannels
	if dbResult := dbConn.Where("customer_id = ?", customerID.String()).Find(&customerNotificationChannels); dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			jsonResponse(w, http.StatusNotFound, struct{}{})
			return
		}

		log.Err(dbResult.Error).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Error while fetching from database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	jsonResponse(w, http.StatusOK, customerNotificationChannels)
}

func (app *App) updateCustomerNotificationChannels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	customerID, err := uuid.Parse(vars["id"])
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	notificationChannelTypeID, err := uuid.Parse(vars["notificationChannelTypeID"])
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid notificationChannelTypeID")
		return
	}

	var customerNotificationChannel customerNotificationChannelCreateUpdateStruct
	err = json.NewDecoder(r.Body).Decode(&customerNotificationChannel)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Could not unmarshal body")
		return
	}

	log.Debug().
		Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
		Str("customerID", customerID.String()).
		Str("notificationChannelTypeID", notificationChannelTypeID.String()).
		Str("handler", "updateCustomerNotificationChannels")

	dbConn, err := db.GetConnection()

	if err != nil {
		log.Err(err).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Could not retrieve connection to the database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	var customerNotificationChannels []db.CustomerNotificationChannels
	dbResult := dbConn.
		Find(&customerNotificationChannels).
		Where("customer_id = ? AND notification_channel_type_id = ?", customerID.String(), notificationChannelTypeID.String()).
		Updates(map[string]interface{}{"contact_customer": customerNotificationChannel.ContactCustomer, "notification_channel_lookup_key": customerNotificationChannel.LookupKey})

	if dbResult.Error != nil {
		log.Err(dbResult.Error).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Error while updating database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	// nothing updated so it doesn't exist; insert it
	if dbResult.RowsAffected == 0 {
		if dbResult := dbConn.Create(&db.CustomerNotificationChannels{
			CustomerID:                   customerID.String(),
			NotificationChannelTypeID:    notificationChannelTypeID.String(),
			NotificationChannelLookupKey: customerNotificationChannel.LookupKey,
			ContactCustomer:              customerNotificationChannel.ContactCustomer,
		}); dbResult.Error != nil {
			log.Err(dbResult.Error).
				Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
				Msg("Error while inserting to database")

			errorResponse(w, http.StatusInternalServerError, "Error occurred")
			return
		}
	}

	jsonResponse(w, http.StatusOK, nil)
}

func (app *App) addNotification(w http.ResponseWriter, r *http.Request) {
	var notification notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Could not unmarshal body")
		return
	}

	log.Debug().
		Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
		Str("handler", "addNotification")

	dbConn, err := db.GetConnection()

	if err != nil {
		log.Err(err).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Could not retrieve connection to the database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	newNotification := db.Notification{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		From:      notification.From,
		Text:      notification.Text,
	}

	if dbResult := dbConn.Create(newNotification); dbResult.Error != nil {
		log.Err(dbResult.Error).
			Str("requestIdentifier", r.Context().Value("requestIdentifier").(string)).
			Msg("Error while updating database")

		errorResponse(w, http.StatusInternalServerError, "Error occurred")
		return
	}

	jsonResponse(w, http.StatusCreated, newNotification)
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/customers", app.getCustomers).Methods(http.MethodGet)
	app.Router.HandleFunc("/customers", app.addCustomer).Methods(http.MethodPost)
	app.Router.HandleFunc("/customers/{id}", app.getCustomer).Methods(http.MethodGet)
	app.Router.HandleFunc("/customers/{id}", app.updateCustomer).Methods(http.MethodPut)
	app.Router.HandleFunc("/customers/{id}", app.deleteCustomer).Methods(http.MethodDelete)

	app.Router.HandleFunc("/customers/{id}/notification-channels", app.getCustomerNotificationChannels).Methods(http.MethodGet)
	app.Router.HandleFunc("/customers/{id}/notification-channels/{notificationChannelTypeID}", app.updateCustomerNotificationChannels).Methods(http.MethodPost)

	app.Router.HandleFunc("/notifications", app.addNotification).Methods(http.MethodPost)
}

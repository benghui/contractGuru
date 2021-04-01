package router

import (
	"net/http"

	"github.com/contractGuru/pkg/db"
	"github.com/contractGuru/pkg/handlers"
	"github.com/contractGuru/pkg/middleware"
	"github.com/gorilla/mux"
)

// GetRouter handles routing.
func GetRouter(db *db.DB) *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/users/login", handlers.LoginUser(db)).Methods(http.MethodPost)
	api.HandleFunc("/users/logout", handlers.LogoutUser(db)).Methods(http.MethodPost)

	api.HandleFunc("/requests", handlers.GetPendingRequests(db)).Methods(http.MethodGet)
	api.HandleFunc("/requests", handlers.CreateRequest(db)).Methods(http.MethodPost)

	api.HandleFunc("/requests/{id:[0-9]+}/actions", handlers.GetPendingRequestActions(db)).Methods(http.MethodGet)
	api.HandleFunc("/requests/{id:[0-9]+}/transitions", handlers.CreateRequestTransition(db)).Methods(http.MethodPost)


	api.Use(middleware.LoggingMiddleware)

	return api
}

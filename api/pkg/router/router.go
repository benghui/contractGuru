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

	api.Use(middleware.LoggingMiddleware)

	return api
}

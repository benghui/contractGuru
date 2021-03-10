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

	api.HandleFunc("/users", handlers.GetUsers(db)).Methods(http.MethodGet)

	api.Use(middleware.LoggingMiddleware)

	return api
}

package router

import (
	"net/http"

	"github.com/contractGuru/pkg/application"
	"github.com/contractGuru/pkg/handlers"
	"github.com/contractGuru/pkg/middleware"
	"github.com/gorilla/mux"
)

// GetRouter handles routing.
func GetRouter(app *application.Application) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/users", handlers.GetUsers(app)).Methods(http.MethodGet)

	router.Use(middleware.LoggingMiddleware)

	return router
}

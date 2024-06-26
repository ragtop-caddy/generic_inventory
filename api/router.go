package api

import (
	"generic_inventory/auth"
	"generic_inventory/logger"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter - Function to create a new router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)
		handler = auth.ValidateSession(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}

package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route - Struct to hold route information
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes - Struct to hold multiple route definitions
type Routes []Route

// NewRouter - Function called by main to create a new router
func NewRouter() *mux.Router {

	entries = append(entries, Entry{
		SKU: "1",
		Header: &Header{
			Type:        "pant",
			Description: "boys small pant",
			Stock:       5,
			History: []Transaction{
				Transaction{
					ISODate:  "1234",
					Campus:   "duffy",
					Students: 0,
					Action:   "create",
				},
			},
		},
		Details: &Detail{
			Gender: "m",
			Color:  "blk",
			Size:   "small",
			Style:  "slacks",
			Fit:    "loose",
		},
	},
	)

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"All Entries",
		"GET",
		"/inventory",
		GetEntries,
	},
	Route{
		"Get Entry",
		"GET",
		"/inventory/{sku}",
		GetEntry,
	},
	Route{
		"Create Entry",
		"POST",
		"/inventory/{sku}",
		CreateEntry,
	},
}

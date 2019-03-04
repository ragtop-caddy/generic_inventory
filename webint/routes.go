package main

import "net/http"

// Route - Struct to hold route information
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes - Struct to hold multiple route definitions
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		GetIndex,
	},
}

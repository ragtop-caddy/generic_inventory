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
		"All Entries",
		"GET",
		"/api/inventory",
		GetEntries,
	},
	Route{
		"Get Entry",
		"GET",
		"/api/inventory/{sku}",
		GetEntry,
	},
	Route{
		"Create Entry",
		"POST",
		"/api/inventory/{sku}",
		CreateEntry,
	},
	Route{
		"Delete Entry",
		"DELETE",
		"/api/inventory/{sku}",
		DeleteEntry,
	},
	Route{
		"Index",
		"GET",
		"/index",
		GetIndex,
	},
	Route{
		"Stylesheet File",
		"GET",
		"/css/{cssfile}",
		GetCSS,
	},
	Route{
		"Javascript File",
		"GET",
		"/js/{jsfile}",
		GetJS,
	},
	Route{
		"Image File",
		"GET",
		"/img/{imgfile}",
		GetIMG,
	},
}

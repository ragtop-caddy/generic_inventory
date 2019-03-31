package api

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
		"Get Entry",
		"GET",
		"/api/inventory/{action}/{sku}",
		CrudHandle,
	},
	Route{
		"Create Entry",
		"POST",
		"/api/inventory/{action}/{sku}",
		CrudHandle,
	},
	Route{
		"Delete Entry",
		"DELETE",
		"/api/inventory/{action}/{sku}",
		CrudHandle,
	},
	Route{
		"Index",
		"GET",
		"/",
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

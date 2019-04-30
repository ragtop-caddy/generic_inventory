package api

import (
	"generic_inventory/auth"
	"net/http"
)

// Complex YAML for route/router config
//type Service struct {
//    APIVersion string `yaml:"apiVersion"`
//    Kind       string `yaml:"kind"`
//    Metadata   struct {
//        Name      string `yaml:"name"`
//        Namespace string `yaml:"namespace"`
//        Labels    struct {
//            RouterDeisIoRoutable string `yaml:"router.deis.io/routable"`
//        } `yaml:"labels"`
//        Annotations struct {
//            RouterDeisIoDomains string `yaml:"router.deis.io/domains"`
//        } `yaml:"annotations"`
//    } `yaml:"metadata"`
//    Spec struct {
//        Type     string `yaml:"type"`
//        Selector struct {
//            App string `yaml:"app"`
//        } `yaml:"selector"`
//        Ports []struct {
//            Name       string `yaml:"name"`
//            Port       int    `yaml:"port"`
//            TargetPort int    `yaml:"targetPort"`
//            NodePort   int    `yaml:"nodePort,omitempty"`
//        } `yaml:"ports"`
//    } `yaml:"spec"`
//}

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
		"Admin",
		"GET",
		"/admin",
		auth.AdminPanel,
	},
	Route{
		"Login",
		"POST",
		"/login",
		auth.Login,
	},
	Route{
		"Logout",
		"GET",
		"/logout",
		auth.Logout,
	},
	Route{
		"Get User",
		"GET",
		"/api/user/{action}/{uid}",
		auth.AdminCrudHandle,
	},
	Route{
		"Create User",
		"POST",
		"/api/user/{action}/{uid}",
		auth.AdminCrudHandle,
	},
	Route{
		"Delete User",
		"DELETE",
		"/api/user/{action}/{uid}",
		auth.AdminCrudHandle,
	},
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
		"Index",
		"GET",
		"/{path}/{file}.{ext}",
		StaticHandle,
	},
}

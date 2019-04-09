package main

import (
	"generic_inventory/api"
	"generic_inventory/conf"
	"generic_inventory/dao"
	"log"
	"net/http"
)

// main - our main function
func main() {
	var c = &conf.MyConfig
	c.ParseConfig()
	c.Router = api.NewRouter()

	// Configure DB Connection
	dao.ConfigDB(c)

	// Configure Server using TLS
	c.ConfigureTLS()
	c.ConfigServer()

	// Startup Standard HTTP Listener
	go http.ListenAndServe(":80", http.HandlerFunc(api.RedirectToTLS))

	// Startup TLS Listener
	log.Fatal(conf.MyConfig.SrvConf.ListenAndServeTLS(c.TLSCert, c.TLSKey))
}

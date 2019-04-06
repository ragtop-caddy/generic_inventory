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
	c.DBHost = "localhost"
	c.DBPort = "27017"
	c.DBName = "inventory"
	c.TLSCert = "F:/Docker/generic_inventory/ssl/cert.pem"
	c.TLSKey = "F:/Docker/generic_inventory/ssl/key.pem"
	c.StaticPath = "F:/Docker/generic_inventory/static/"
	c.TmplPath = "F:/Docker/generic_inventory/templates/"
	c.Addr = ":443"
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

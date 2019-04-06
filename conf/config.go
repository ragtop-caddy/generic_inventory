package conf

import (
	"crypto/tls"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// ServerConf - Object to hold server configuration values
type ServerConf struct {
	Router     *mux.Router
	TLSConf    *tls.Config
	SrvConf    *http.Server
	Addr       string
	TLSCert    string
	TLSKey     string
	StaticPath string
	TmplPath   string
	DBHost     string
	DBPort     string
	DBName     string
	DBClient   *mongo.Database
}

// MyConfig - Exported variable to hold server configuration data
var MyConfig ServerConf

// ConfigureTLS - Configure TLS settings
func (conf *ServerConf) ConfigureTLS() {
	conf.TLSConf = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
}

// ConfigServer - A function to setup basic server configuration
func (conf *ServerConf) ConfigServer() {
	conf.SrvConf = &http.Server{
		Addr:         conf.Addr,
		Handler:      conf.Router,
		TLSConfig:    conf.TLSConf,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
}

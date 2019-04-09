package conf

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v2"
)

// ServerConf - Object to hold server configuration values
type ServerConf struct {
	ConfigFile string
	Router     *mux.Router
	TLSConf    *tls.Config
	SrvConf    *http.Server
	DBClient   *mongo.Database
	Addr       string `yaml:"ssl_addr,omitempty"`
	TLSCert    string `yaml:"cert,omitempty"`
	TLSKey     string `yaml:"key,omitempty"`
	StaticPath string `yaml:"static_path,omitempty"`
	TmplPath   string `yaml:"template_path,omitempty"`
	DBHost     string `yaml:"dbhost,omitempty"`
	DBPort     string `yaml:"dbport,omitempty"`
	DBName     string `yaml:"dbname,omitempty"`
}

// MyConfig - Exported variable to hold server configuration data. Is initialized with defaults.
var MyConfig ServerConf

// ParseConfig - Function to parse the configuration file
func (conf *ServerConf) ParseConfig() {
	conf.ConfigFile = "/etc/generic_inventory/config.yaml"
	conf.Addr = ":443"
	conf.TLSCert = "/etc/generic_inventory/ssl/cert.pem"
	conf.TLSKey = "/etc/generic_inventory/ssl/key.pem"
	conf.StaticPath = "/etc/generic_inventory/static/"
	conf.TmplPath = "/etc/generic_inventory/templates/"
	conf.DBHost = "localhost"
	conf.DBPort = "27017"
	conf.DBName = "inventory"

	file, ok := os.LookupEnv("INV_CONF_FILE")
	if !ok {
		file = conf.ConfigFile
	}
	yamlFile, err := ioutil.ReadFile(file)
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		panic(err)
	}
}

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

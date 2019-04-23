package conf

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/yaml.v2"
)

// ServerConf - Object to hold server configuration values
type ServerConf struct {
	ConfigFile    string
	Router        *mux.Router
	TLSConf       *tls.Config
	SrvConf       *http.Server
	MongoClient   *mongo.Client
	DBClient      *mongo.Database
	Addr          string `yaml:"ssl_addr,omitempty"`
	TLSCert       string `yaml:"cert,omitempty"`
	TLSKey        string `yaml:"key,omitempty"`
	RootCAs       string `yaml:"root_ca_bundle,omitempty"`
	StaticPath    string `yaml:"static_path,omitempty"`
	TmplPath      string `yaml:"template_path,omitempty"`
	DBHost        string `yaml:"dbhost,omitempty"`
	DBPort        string `yaml:"dbport,omitempty"`
	DBName        string `yaml:"dbname,omitempty"`
	ClientTLScert string `yaml:"ssl_client_cert,omitempty"`
	ClientTLSkey  string `yaml:"ssl_client_key,omitempty"`
	ClientTLS     *tls.Config
}

// MyConfig - Exported variable initialized with default configuration data.
var MyConfig = ServerConf{
	ConfigFile:    "/etc/generic_inventory/config.yaml",
	Addr:          ":443",
	RootCAs:       "/etc/pki/tls/cert.pem",
	TLSCert:       "/etc/generic_inventory/ssl/cert.pem",
	TLSKey:        "/etc/generic_inventory/ssl/key.pem",
	StaticPath:    "/etc/generic_inventory/static/",
	TmplPath:      "/etc/generic_inventory/templates/",
	DBHost:        "localhost",
	DBPort:        "27017",
	DBName:        "inventory",
	ClientTLScert: "/etc/generic_inventory/ssl/client_cert.pem",
	ClientTLSkey:  "/etc/generic_inventory/ssl/client_key.pem",
}

// ParseConfig - Function to parse the configuration file
func (conf *ServerConf) ParseConfig() {
	file, ok := os.LookupEnv("INV_CONF_FILE")
	if !ok {
		file = conf.ConfigFile
	}
	yamlFile, err := ioutil.ReadFile(file)
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Got %s Opening Configuration file", err)
	}
}

// ConfigureClientTLS - Configure TLS settings for client connection
func (conf *ServerConf) ConfigureClientTLS() {
	rootCA := x509.NewCertPool()
	capem, err := ioutil.ReadFile(conf.RootCAs)
	if err != nil {
		log.Fatalf("Got %s opening RootCA Bundle %s", conf.RootCAs, err)
	}
	rootCA.AppendCertsFromPEM(capem)
	clientCert, err := tls.LoadX509KeyPair(conf.ClientTLScert, conf.ClientTLSkey)
	if err != nil {
		log.Fatalf("Got %s loading certificate %s", conf.ClientTLScert, err)
	}

	conf.ClientTLS = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		RootCAs:      rootCA,
		Certificates: []tls.Certificate{clientCert},
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

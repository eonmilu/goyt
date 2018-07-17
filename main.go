package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	// FilePath : path to the files to be served
	FilePath = "/var/www/public"
	// CertPath : path to the TLS certificate file
	CertPath = "/etc/letsencrypt/archive/oxygenrain.com/cert1.pem"
	// KeyPath : path to the TLS private key file
	KeyPath = "/etc/letsencrypt/archive/oxygenrain.com/privkey1.pem"
	// RootDomain : A-record of the domain
	RootDomain = "oxygenrain.com"
	// ConfigPath : path to the user, password and database
	ConfigPath = "/etc/postgresql/dietpi.cfg"
)

var (
	dbinfo string
	err    error
	db     *sql.DB
	cfg    Config
	ok     bool
)

func init() {
	raw, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(raw, &cfg)
	if err != nil {
		log.Panic(err)
	}
	dbinfo = fmt.Sprintf("user=%s password=%s dbname=%s", cfg.User, cfg.Password, cfg.Database)
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Panic(err)
	}
}

// Config : configuration used for instantiating the database
type Config struct {
	User     string
	Password string
	Database string
}

// Redirect the incoming HTTP request to HTTPS
func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	target := RootDomain + r.URL.RequestURI()
	http.Redirect(w, r, "https://"+target, http.StatusMovedPermanently)
	log.Printf("REDIRECT %s FROM %s TO %s", r.RemoteAddr, "http://"+target, "https://"+target)
}

func main() {
	defer db.Close()
	http.Handle("/", http.FileServer(http.Dir("/var/www/public")))
	http.HandleFunc("/yourtime/search", timemarksHandler)
	// Listen to HTTP trafic and redirect it to HTTPS
	go log.Panic(http.ListenAndServe(":8080", http.HandlerFunc(redirectToHTTPS)))
	log.Panic(http.ListenAndServeTLS(":8443", CertPath, KeyPath, nil))
}

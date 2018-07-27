package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	// FilePath : path to the files to be served
	FilePath = "/var/www/public/"
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
	db     *sql.DB
	cfg    Config
	err    error
	ok     bool
)

func init() {
	// Read credentials and open connection to the database
	log.Println("Opening connection to the database...")
	defer log.Println("Done")
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

func main() {
	defer db.Close()
	// Redirect the incoming HTTP request to HTTPS
	go http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer log.Printf("HTTP redirect enabled")
		target := RootDomain + r.URL.RequestURI()
		http.Redirect(w, r, "https://"+target, http.StatusMovedPermanently)
		log.Printf("REDIRECT %s FROM %s TO %s", r.RemoteAddr, "http://"+target, "https://"+target)
	}))
	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir(FilePath)))
	r.HandleFunc("/yourtime/search", searchTimemarksHandler).Methods("GET", "OPTIONS")

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8443",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second}
	log.Panic(srv.ListenAndServeTLS(CertPath, KeyPath))
}

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
	dbcfg  Config
)

func init() {
	// Read credentials and open connection to the database
	log.Println("Opening connection to the database...")
	defer log.Println("Done")
	raw, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(raw, &dbcfg)
	if err != nil {
		log.Panic(err)
	}
	dbinfo = fmt.Sprintf("user=%s password=%s dbname=%s", dbcfg.User, dbcfg.Password, dbcfg.Database)
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Panic(err)
	}
}

// Config : configuration used to initialize the database
type Config struct {
	User     string
	Password string
	Database string
}

func main() {
	defer db.Close()
	// Redirect the incoming HTTP request to HTTPS
	go http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := RootDomain + r.URL.RequestURI()
		http.Redirect(w, r, "https://"+target, http.StatusMovedPermanently)
		log.Printf("REDIRECT %s FROM %s TO %s", r.RemoteAddr, "http://"+target, "https://"+target)
	}))
	r := mux.NewRouter()
	r.HandleFunc("/yourtime/search", searchYourTimeAPI)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(FilePath)))

	log.Panic(http.ListenAndServeTLS(":8443", CertPath, KeyPath, r))
}

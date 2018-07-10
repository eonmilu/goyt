package main

import (
	"database/sql"
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
)

var (
	db, err = sql.Open("postgres", "yt.db")
)

// Redirect the incoming HTTP request to HTTPS
func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	target := RootDomain + r.URL.RequestURI()
	http.Redirect(w, r, "https://"+target, http.StatusMovedPermanently)
	log.Printf("REDIRECT %s FROM %s TO %s", r.RemoteAddr, "http://"+target, "https://"+target)
}

func main() {
	// Check for any errors opening the database
	if err != nil {
		log.Panic(err)
	}
	http.Handle("/", http.FileServer(http.Dir("/var/www/public")))
	http.HandleFunc("/yourtime/search", timemarksHandler)
	// Listen to HTTP trafic and redirect it to HTTPS
	go log.Panic(http.ListenAndServe(":8080", http.HandlerFunc(redirectToHTTPS)))
	log.Panic(http.ListenAndServeTLS(":8443", CertPath, KeyPath, nil))
}

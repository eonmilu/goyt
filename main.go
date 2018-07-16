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
	ConfigPath = "/etc/postgres/dietpi.cfg"
)

var (
	dbinfo string
	err    error
	db, _  = sql.Open("postgres", "")
	cfg    Config
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

// Timemark : data structure for a timemark
type Timemark struct {
	timemarkID int64
	author     string
	authorURL  string
	time       int64
	content    string
	votes      int64
	date       int64
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

func timemarksHandler(w http.ResponseWriter, r *http.Request) {
	t := Timemark{}
	videoID := r.URL.Query().Get("v")
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "0"
	}
	rows, err := db.Query("SELECT id, author, authorURL, time, content, votes, date FROM timemarks WHERE id = $1 ORDER BY votes OFFSET $2 LIMIT $3", videoID, offset, limit)
	defer rows.Close()
	if err != nil {
		fmt.Fprintf(w, "")
		log.Printf("DATABASE CONNECTION FAILED %s IP %s", err, r.RemoteAddr)
	} else {
		for rows.Next() {
			err = rows.Scan(&t.timemarkID, &t.author, &t.authorURL, &t.time, &t.content, &t.votes, &t.date)
		}
		s, err := json.Marshal(t)
		if err != nil {
			log.Printf("QUERY ERROR %s IP %s", err, r.RemoteAddr)
		} else {
			fmt.Fprintf(w, string(s))
		}
	}
}

func main() {
	defer db.Close()
	http.Handle("/", http.FileServer(http.Dir("/var/www/public")))
	http.HandleFunc("/yourtime/search", timemarksHandler)
	// Listen to HTTP trafic and redirect it to HTTPS
	go log.Panic(http.ListenAndServe(":8080", http.HandlerFunc(redirectToHTTPS)))
	log.Panic(http.ListenAndServeTLS(":8443", CertPath, KeyPath, nil))
}

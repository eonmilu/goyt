package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Timemark : data structure for a timemark
type Timemark struct {
	TimemarkID int64  `json:"timemarkID"`
	Author     string `json:"author"`
	AuthorURL  string `json:"authorURL"`
	Timemark   int64  `json:"timemark"`
	Content    string `json:"content"`
	Votes      int64  `json:"votes"`
	Date       int64  `json:"date"`
}

const (
	// SCFound : Status code for timemarks found
	SCFound = "200"
	// SCNotFound : Status code for NO timemarks found
	SCNotFound = "210"
	// SCError : Status code for an internal server error
	SCError = "220"
)

const (
	// tmrkDefOffset : the default offset for a timemark lookup
	tmrkDefOffset = "0"
	// tmrkDefLimit : the default limit for a timemark lookup
	tmrkDefLimit = "10"
)

func searchTimemarksHandler(w http.ResponseWriter, r *http.Request) {
	var (
		offset string
		limit  string
		ok     bool
	)
	t := Timemark{}
	videoID := r.URL.Query().Get("v")
	// Get and set offset to default value if not specified
	if offset = r.URL.Query().Get("offset"); offset == "" {
		offset = tmrkDefOffset
	}
	// Get and set limit to default value if not specified
	if limit = r.URL.Query().Get("limit"); limit == "" {
		limit = tmrkDefLimit
	}
	rows, err := db.Query("SELECT id, author, authorURL, timemark, content, votes, date FROM timemarks WHERE videoid = $1 ORDER BY votes OFFSET $2 LIMIT $3", videoID, offset, limit)
	defer rows.Close()
	if err != nil {
		fmt.Fprintf(w, SCError)
		log.Printf("DATABASE CONNECTION FAILED %s IP %s", err, r.RemoteAddr)
		return
	}
	for rows.Next() {
		ok = true
		err = rows.Scan(&t.TimemarkID, &t.Author, &t.AuthorURL, &t.Timemark, &t.Content, &t.Votes, &t.Date)
	}
	if err != nil {
		log.Printf("QUERY ERROR %s IP %s", err, r.RemoteAddr)
		fmt.Fprintf(w, SCError)
		return
	}
	// Check ok to see if there were any results
	if !ok {
		fmt.Fprintf(w, SCNotFound)
		return
	}
	s, err := json.Marshal(t)
	if err != nil {
		log.Printf("JSON ERROR %s IP %s", err, r.RemoteAddr)
		fmt.Fprintf(w, SCError)
		return
	}
	fmt.Fprintf(w, SCFound+string(s))
}

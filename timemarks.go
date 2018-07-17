package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

func timemarksHandler(w http.ResponseWriter, r *http.Request) {
	t := Timemark{}
	videoID := r.URL.Query().Get("v")
	// Get and set offset to default value
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}
	// Get and set limit to default value
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "10"
	}
	rows, err := db.Query("SELECT id, author, authorURL, time, content, votes, date FROM timemarks WHERE id = $1 ORDER BY votes OFFSET $2 LIMIT $3", videoID, offset, limit)
	defer rows.Close()
	if err != nil {
		fmt.Fprintf(w, "220|")
		log.Printf("DATABASE CONNECTION FAILED %s IP %s", err, r.RemoteAddr)
	} else {
		/* STATUS CODES
		   200: Found
		   210: Not found
		   220: Internal error
		*/
		for rows.Next() {
			ok = true
			err = rows.Scan(&t.timemarkID, &t.author, &t.authorURL, &t.time, &t.content, &t.votes, &t.date)
		}
		if rows.Err() != nil {
			log.Printf("QUERY ERROR %s IP %s", err, r.RemoteAddr)
			fmt.Fprintf(w, "220|")
		}
		// Check ok to see if there were any results
		if !ok {
			fmt.Fprintf(w, "210|")
		}
		s, err := json.Marshal(t)
		if err != nil {
			log.Printf("JSON ERROR %s IP %s", err, r.RemoteAddr)
			fmt.Fprintf(w, "220|")
		} else {
			fmt.Fprintf(w, "200|"+string(s))
		}
	}
}

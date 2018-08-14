package yourtime

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Search writes to the ResponseWriter a JSON-encoded array
// of Timemark objects matching the given URL, offset and limit
func (y YourTime) Search(w http.ResponseWriter, r *http.Request) {
	// Allow CORS
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var (
		offset string
		limit  string
		t      Timemark
		p      []Timemark
	)
	videoID := r.URL.Query().Get("v")
	// Get and set offset to default value if not specified
	if offset = r.URL.Query().Get("offset"); offset == "" {
		offset = tmrkDefOffset
	}
	// Get and set limit to default value if not specified
	if limit = r.URL.Query().Get("limit"); limit == "" {
		limit = tmrkDefLimit
	}
	rows, err := y.DB.Query("SELECT id, author, authorURL, timemark, content, votes, date FROM timemarks WHERE videoid = $1 ORDER BY votes OFFSET $2 LIMIT $3", videoID, offset, limit)
	defer rows.Close()
	if err != nil {
		fmt.Fprintf(w, sCError)
		log.Printf("DATABASE CONNECTION FAILED %s IP %s", err, r.RemoteAddr)
		return
	}
	for rows.Next() {
		err = rows.Scan(&t.TimemarkID, &t.Author, &t.AuthorURL, &t.Timemark, &t.Content, &t.Votes, &t.Date)
		p = append(p, t)
	}
	if err != nil {
		fmt.Fprintf(w, sCError)
		log.Printf("QUERY ERROR %s IP %s", err, r.RemoteAddr)
		return
	}
	// Check p's length to see if there were any results
	if len(p) == 0 {
		fmt.Fprintf(w, sCNotFound)
		return
	}
	s, err := json.Marshal(p)
	if err != nil {
		fmt.Fprintf(w, sCError)
		log.Printf("JSON ERROR %s IP %s", err, r.RemoteAddr)
		return
	}
	fmt.Fprintf(w, sCFound+string(s))
}

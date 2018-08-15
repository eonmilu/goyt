package goyt

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// Search writes to the ResponseWriter a JSON-encoded array
// of Timemark objects matching the given URL, offset and limit
func (y YourTime) Search(w http.ResponseWriter, r *http.Request) {
	corsWriter{w}.enableCORS()

	params := parameters{
		videoID: r.URL.Query().Get("v"),
		offset:  offset(r.URL.Query().Get("offset")),
		limit:   limit(r.URL.Query().Get("limit")),
	}

	params.checkParametersValue()
	p, err := timemarksDB{y.DB}.getTimemarks(params)
	if err != nil {
		fmt.Fprintf(w, sCError)
		return
	}

	// Check p's length to see if there were any result
	if len(p) == 0 {
		fmt.Fprintf(w, sCNotFound)
		return
	}

	s, err := json.Marshal(p)
	if err != nil {
		fmt.Fprintf(w, sCError)
		return
	}
	fmt.Fprintf(w, sCFound+string(s))
}

type corsWriter struct {
	http.ResponseWriter
}

func (w corsWriter) enableCORS() {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

type timemarksDB struct {
	*sql.DB
}

func (tdb timemarksDB) getTimemarks(params parameters) ([]Timemark, error) {
	var (
		t Timemark
		p []Timemark
	)
	rows, err := tdb.Query("SELECT id, author, authorURL, timemark, content, votes, date FROM timemarks WHERE videoid = $1 ORDER BY votes OFFSET $2 LIMIT $3", params.videoID, params.offset, params.limit)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&t.TimemarkID, &t.Author, &t.AuthorURL, &t.Timemark, &t.Content, &t.Votes, &t.Date)
		p = append(p, t)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

type parameters struct {
	videoID string
	offset  offset
	limit   limit
}

type timemarkRange interface {
	setToDefault()
	isEmpty() bool
}

func (p parameters) checkParametersValue() {
	if p.offset.isEmpty() {
		p.offset.setToDefault()
	}
	if p.limit.isEmpty() {
		p.limit.setToDefault()
	}
}

type offset string

func (o offset) setToDefault() {
	o = tmrkDefOffset
}

func (o offset) isEmpty() bool {
	return o == ""
}

type limit string

func (l limit) setToDefault() {
	l = tmrkDefLimit
}

func (l limit) isEmpty() bool {
	return l == ""
}

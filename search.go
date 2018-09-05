package goyt

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Search writes to the ResponseWriter a JSON-encoded array
// of Timemark objects matching the given URL, offset and limit
func (y YourTime) Search(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)

	params := parameters{
		videoID: r.URL.Query().Get("v"),
		offset:  offset(r.URL.Query().Get("offset")),
		limit:   limit(r.URL.Query().Get("limit")),
	}
	params.checkParameters()

	p, err := y.getTimemarks(params)
	if err != nil {
		fmt.Fprintf(w, sCError)
		fmt.Printf("%s", err)
		return
	}

	// Check p's length to see if there was any result
	if len(p) == 0 {
		fmt.Fprintf(w, sCNotFound)
		return
	}

	s, err := json.Marshal(p)
	if err != nil {
		fmt.Fprintf(w, sCError)
		fmt.Printf("%s", err)
		return
	}
	fmt.Fprintf(w, sCFound+string(s))
}

func (y YourTime) getTimemarks(params parameters) ([]Timemark, error) {
	var (
		t Timemark
		p []Timemark
	)
	rows, err := y.DB.Query("SELECT id, timemark, content, votes, author, approved, timestamp FROM timemarks WHERE videoid = $1 ORDER BY votes DESC OFFSET $2 LIMIT $3", params.videoID, params.offset, params.limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t = Timemark{}
		err = rows.Scan(&t.ID, &t.Timemark, &t.Content, &t.Votes, &t.Author, &t.Approved, &t.Timestamp)
		p = append(p, t)
		if err != nil {
			return nil, err
		}
	}
	if rows.Err() != nil {
		return nil, rows.Err()
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

func (p *parameters) checkParameters() {
	if p.offset.isEmpty() {
		p.offset.setToDefault()
	}
	if p.limit.isEmpty() {
		p.limit.setToDefault()
	}
}

type offset string

func (o *offset) setToDefault() {
	*o = tmrkDefOffset
}

func (o offset) isEmpty() bool {
	return o == ""
}

type limit string

func (l *limit) setToDefault() {
	*l = tmrkDefLimit
}

func (l limit) isEmpty() bool {
	return l == ""
}

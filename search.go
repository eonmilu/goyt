package goyt

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

var fallbackAuthor = Author{-1, "Undefined", ""}

// Search writes to the ResponseWriter a JSON-encoded array
// of Timemark objects matching the given URL, offset and limit
func (y YourTime) Search(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)
	var sp searchResponse
	var err error

	params := parameters{
		videoID: r.URL.Query().Get("v"),
		offset:  offset(r.URL.Query().Get("offset")),
		limit:   limit(r.URL.Query().Get("limit")),
	}
	params.checkParameters()

	sp.Timemarks, err = y.getTimemarks(params)
	if err != nil {
		fmt.Fprintf(w, sCError)
		fmt.Printf("%s", err)
		return
	}

	// Check sp.Timemarks's length to see if there was any result
	if len(sp.Timemarks) == 0 {
		fmt.Fprintf(w, sCNotFound)
		return
	}

	sp.Authors = make([]Author, len(sp.Timemarks))

	for i := 0; i < len(sp.Timemarks); i++ {
		author, err := y.getTimemarkAuthor(sp.Timemarks[i].Author)
		if err != nil {
			author = fallbackAuthor
		}
		sp.Authors[i] = author
	}

	s, err := json.Marshal(sp)
<<<<<<< HEAD

=======
	fmt.Printf("%s", s)
>>>>>>> 49f628d66b619cb33cf585465291fd5b3a0c114a
	if err != nil {
		fmt.Fprintf(w, sCError)
		fmt.Printf("%s", err)
		return
	}
	fmt.Fprintf(w, sCFound+string(s))
}

func (y YourTime) getTimemarkAuthor(id int64) (Author, error) {
	nullableAuthor := nullAuthor{}
	author := fallbackAuthor
	row := y.DB.QueryRow("SELECT username, url FROM users WHERE id=$1", id)

	err := row.Scan(&nullableAuthor.Username, &nullableAuthor.URL)
	if err != nil {
		return author, err
	}

	if !nullableAuthor.Username.Valid || !nullableAuthor.URL.Valid {
		return fallbackAuthor, err
	}

	author = Author{
		id,
		nullableAuthor.Username.String,
		nullableAuthor.URL.String,
	}
	return author, err
}

type nullAuthor struct {
	ID       sql.NullInt64
	Username sql.NullString
	URL      sql.NullString
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

type searchResponse struct {
	Timemarks []Timemark `json:"timemarks"`
	Authors   []Author   `json:"authors"`
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

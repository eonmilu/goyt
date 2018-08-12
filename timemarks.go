package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Timemark is the data structure for a timemark
type Timemark struct {
	TimemarkID int64  `json:"timemarkID"`
	Author     string `json:"author"`
	AuthorURL  string `json:"authorURL"`
	Timemark   int64  `json:"timemark"`
	Content    string `json:"content"`
	Votes      int64  `json:"votes"`
	Date       int64  `json:"date"`
}

// GoogleValidate is the unmarshaled data structure GoogleValidateTokenURL returns
// We will only declare the types we are going to use (aud, sub)
type GoogleValidate struct {
	Aud string `json:"aud"`
	Sub string `json:"sub"`
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
	// GoogleValidateTokenURL is the URL which Google provides to validate a token
	GoogleValidateTokenURL = "https://www.googleapis.com/oauth2/v3/tokeninfo?"
	// YourTimeGoogleClientID is the URL the request is coming from
	YourTimeGoogleClientID = "817145568720-9p70ci9se6tpvn4qh9vbldh16gssfs3v.apps.googleusercontent.com"
)

const (
	// tmrkDefOffset : the default offset for a timemark lookup
	tmrkDefOffset = "0"
	// tmrkDefLimit : the default limit for a timemark lookup
	tmrkDefLimit = "10"
)

// SearchYourTimeAPI writes to the ResponseWriter a JSON-encoded array
// of Timemark objects matching the given URL, offset and limit
func SearchYourTimeAPI(w http.ResponseWriter, r *http.Request) {
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
	rows, err := DB.Query("SELECT id, author, authorURL, timemark, content, votes, date FROM timemarks WHERE videoid = $1 ORDER BY votes OFFSET $2 LIMIT $3", videoID, offset, limit)
	defer rows.Close()
	if err != nil {
		fmt.Fprintf(w, SCError)
		log.Printf("DATABASE CONNECTION FAILED %s IP %s", err, r.RemoteAddr)
		return
	}
	for rows.Next() {
		err = rows.Scan(&t.TimemarkID, &t.Author, &t.AuthorURL, &t.Timemark, &t.Content, &t.Votes, &t.Date)
		p = append(p, t)
	}
	if err != nil {
		fmt.Fprintf(w, SCError)
		log.Printf("QUERY ERROR %s IP %s", err, r.RemoteAddr)
		return
	}
	// Check p's length to see if there were any results
	if len(p) == 0 {
		fmt.Fprintf(w, SCNotFound)
		return
	}
	s, err := json.Marshal(p)
	if err != nil {
		fmt.Fprintf(w, SCError)
		log.Printf("JSON ERROR %s IP %s", err, r.RemoteAddr)
		return
	}
	fmt.Fprintf(w, SCFound+string(s))
}

// InsertYourTimeAPI TODO:
func InsertYourTimeAPI(w http.ResponseWriter, r *http.Request) {

}

// TokenAuthYourTimeAPI gets the token from the HTTPS POST, validates it
// and overrides the past token from the database
func TokenAuthYourTimeAPI(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	r.ParseForm()
	token := r.Form["idtoken"][0]
	payload := url.Values{}
	payload.Add("id_token", token)
	// Validate the token
	req, err := http.NewRequest("GET", GoogleValidateTokenURL+payload.Encode(), nil)
	if err != nil {
		log.Printf("TOKEN ERROR: %s", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("TOKEN ERROR: %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	gv := &GoogleValidate{}
	err = json.Unmarshal(body, &gv)
	if err != nil {
		log.Printf("TOKEN ERROR: %s", err)
		return
	}
	// If both Aud is from the Your Time client ID and Sub contains anything then the login was succesful
	if gv.Aud == YourTimeGoogleClientID && gv.Sub != "" {
		// TODO: insert gv.Sub to the users table
	}
}

package goyt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Auth is the unmarshaled data structure AuthTokenURL returns
// We will only declare the types we are going to use (aud, sub)
type Auth struct {
	Aud string `json:"aud"`
	Sub string `json:"sub"`
}

// Auth gets the token from the HTTPS POST, validates it
// and overrides the past token from the database
func (y YourTime) Auth(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)

	if token.isLegit(y) {
		// TODO: insert gv.Sub to the users table
		fmt.Fprintf(w, string(token))
	} else {
		fmt.Fprintf(w, sCError)
	}
}

type token string

func (t token) isLegit(y YourTime) bool {
	legit, err := t.validate(y)
	if err != nil {
		log.Printf("%s", err)
		return false
	}
	return legit
}

func getToken(r *http.Request) token {
	r.ParseForm()
	return token(r.Form["idtoken"][0])
}

func (t token) validate(y YourTime) (bool, error) {
	var (
		client  *http.Client
		payload url.Values
	)
	payload.Add("id_token", string(t))

	req, err := http.NewRequest("GET", y.AuthTokenURL+payload.Encode(), nil)
	if err != nil {
		log.Printf("%s", err)
		return false, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%s", err)
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	gv := &Auth{}
	err = json.Unmarshal(body, &gv)
	if err != nil {
		log.Printf("%s", err)
		return false, err
	}
	// If both Aud is from the Your Time client ID and Sub contains anything then the login was succesful
	return gv.Aud == y.GoogleClientID && gv.Sub != "", nil
}

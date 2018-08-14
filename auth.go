package yourtime

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
	client := &http.Client{}
	r.ParseForm()
	token := r.Form["idtoken"][0]
	payload := url.Values{}
	payload.Add("id_token", token)
	// Validate the token
	req, err := http.NewRequest("GET", y.AuthTokenURL+payload.Encode(), nil)
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
	gv := &Auth{}
	err = json.Unmarshal(body, &gv)
	if err != nil {
		log.Printf("TOKEN ERROR: %s", err)
		return
	}
	// If both Aud is from the Your Time client ID and Sub contains anything then the login was succesful
	if gv.Aud == y.GoogleClientID && gv.Sub != "" {
		// TODO: insert gv.Sub to the users table
		// Write token
		fmt.Fprintf(w, gv.Sub)
	}
}

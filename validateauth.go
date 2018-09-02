package goyt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// ValidateAuth gets the token from the HTTPS POST, validates it
// and overrides the past token from the database
func (y YourTime) ValidateAuth(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)

	user := User{}
	r.ParseForm()
	token := getToken(r)

	isLegit, err := token.GetIfLegit(y, &user)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	if !isLegit {
		fmt.Fprintf(w, sCBadLogin)
		return
	}

	user.token = string(token)

	err = y.handleExistingUser(user)
	if err != nil {
		log.Printf("%s", err)
		fmt.Fprintf(w, sCError)
		return
	}
	cookie := http.Cookie{
		Name:    "yourtime-token-server",
		Path:    "/",
		Value:   string(token),
		Expires: time.Now().Add(32 * 365 * 24 * time.Hour),
		Secure:  true,
	}
	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, sCOK)
}

func (t timemarksDB) userExistsByIdentifier(identifier string) (bool, error) {
	result := false
	row := t.QueryRow("SELECT exists(SELECT 1 FROM users WHERE identifier=$1)", identifier)
	err := row.Scan(&result)
	return result, err
}

func (y YourTime) handleNewUser(user User) error {
	_, err := y.DB.Exec("INSERT INTO users (token, identifier) VALUES ($1, $2)", user.token, user.Identifier)
	return err
}

func (y YourTime) handleExistingUser(user User) error {
	_, err := y.DB.Exec("UPDATE users SET token=$1 WHERE identifier=$2", user.token, user.Identifier)
	return err
}

type token string

func getToken(r *http.Request) token {
	if len(r.Form["idtoken"]) > 0 {
		return token(r.Form["idtoken"][0])
	}
	return ""
}

func (t token) GetIfLegit(y YourTime, user *User) (bool, error) {
	err := y.getUserData(t, user)
	if err != nil {
		log.Printf("%s", err)
		return false, err
	}
	legit := y.isValidResponse(*user)
	return legit, nil
}

func (y YourTime) getUserData(t token, user *User) error {
	payload := url.Values{}
	payload.Add("id_token", string(t))

	req, err := http.NewRequest("GET", y.AuthTokenURL+payload.Encode(), nil)
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%s", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, user)
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	return nil
}

func (y YourTime) isValidResponse(u User) bool {
	// If both Aud is from the Your Time client ID and Sub contains anything then the login was succesful
	return u.Aud == y.GoogleClientID && u.Sub != ""
}

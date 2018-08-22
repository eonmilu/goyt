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

	isExistingUser, err := timemarksDB{y.DB}.userExists(user)
	if err != nil {
		fmt.Fprintf(w, sCError)
		log.Printf("%s", err)
		return
	}

	if isExistingUser {
		err = y.handleExistingUser(user)
		if err != nil {
			log.Printf("%s", err)
			fmt.Fprintf(w, sCError)
			return
		}
	} else {
		err = y.handleNewUser(user)
		log.Printf("Creating")
		if err != nil {
			log.Printf("%s", err)
			fmt.Fprintf(w, sCError)
			return
		}
	}
	fmt.Fprintf(w, string(token))
}

func (t timemarksDB) userExists(user User) (bool, error) {
	result := false
	rows, err := t.Query("SELECT exists(SELECT 1 FROM users WHERE email=$1)", user.Email)
	if err != nil {
		return false, err
	}
	rows.Next()
	err = rows.Scan(&result)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (y YourTime) handleNewUser(user User) error {
	_, err := y.DB.Exec("INSERT INTO users (token, email) VALUES ($1, $2)", user.token, user.Email)
	return err
}

func (y YourTime) handleExistingUser(user User) error {
	_, err := y.DB.Exec("UPDATE users SET token=$1 WHERE email=$2", user.token, user.Email)
	return err
}

type token string

func getToken(r *http.Request) token {
	return token(r.Form["idtoken"][0])
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

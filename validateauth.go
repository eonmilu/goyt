package goyt

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ValidateAuth gets the token from the HTTPS POST, validates it
// and overrides the past token from the database
func (y YourTime) ValidateAuth(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)

	user := User{}

	r.ParseForm()
	channelID := getFormParameter(r, "channelid")
	secretCode := getFormParameter(r, "secretcode")
	if channelID == "" || secretCode == "" {
		fmt.Fprintf(w, sCError)
		return
	}

	isValid, err := y.validateChannel(channelID, secretCode)

	if err != nil {
		log.Printf("%s", err)
		fmt.Fprintf(w, sCError)
		return
	}
	if !isValid {
		fmt.Fprintf(w, sCBadLogin)
		return
	}

	userExists, err := y.userExistsByIdentifier(user.Identifier)

	if userExists {
		err = y.handleExistingUser(user)
	} else {
		err = y.handleNewUser(user)
	}

	if err != nil {
		log.Printf("%s", err)
		fmt.Fprintf(w, sCError)
		return
	}
	token := "not implemented" // TODO: use token login

	cookie := http.Cookie{
		Name:    "yourtime-token-server",
		Path:    "/",
		Value:   token,
		Expires: time.Now().Add(32 * 365 * 24 * time.Hour),
		Secure:  true,
	}
	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, sCOK)
}

func (y YourTime) validateChannel(channelID, secretCode string) (bool, error) {
	channel, err := getChannel(channelID)
	if err != nil {
		return false, err
	}

	// If the secret code is in the channel's description, the request is valid
	return strings.Contains(channel.Description, secretCode), nil
}

func getChannel(id string) (User, error) {
	channel := User{}

	req, err := http.NewRequest("GET", "https://www.youtube.com/channel/"+id, nil)
	if err != nil {
		log.Printf("%s", err)
		return channel, err
	}

	// This User-Agent header decreases the request data from 1MB to ~81KB
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) (compatible; YourTime/1.0; +https://oxygenrain.com/yourtime/crawler.html)")

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%s", err)
		return channel, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	re := regexp.MustCompile(`<!--(?P<match>(.*))-->`)
	rawData := re.FindStringSubmatch(string(body))[1]
	var data youTubeChannelResponse
	err = json.Unmarshal([]byte(rawData), &data)
	if err != nil {
		return channel, err
	}
	temp := data.Metadata.ChannelMetadataRenderer
	channel = User{
		-1,
		temp.ExternalID,
		temp.Title,
		temp.ChannelURL,
		temp.Avatar.Thumbnails[0].URL,
		temp.Description,
	}

	return channel, errors.New("Not implemented")
}

func (y YourTime) userExistsByIdentifier(identifier string) (bool, error) {
	result := false
	row := y.DB.QueryRow("SELECT exists(SELECT 1 FROM users WHERE identifier=$1)", identifier)
	err := row.Scan(&result)
	return result, err
}

func (y YourTime) handleNewUser(user User, token string) error {
	_, err := y.DB.Exec("INSERT INTO users (token, identifier, username, url, picture) VALUES ($1, $2, $3, $4, $5)",
		token, user.Identifier, user.Username, user.URL, user.Picture)
	return err
}

func (y YourTime) handleExistingUser(user User, token string) error {
	_, err := y.DB.Exec("UPDATE users SET token=$1, username=$2, url=$3, picture=$4 WHERE identifier=$5",
		token, user.Username, user.URL, user.Picture, user.Identifier)
	return err
}

type token string // TODO: delete and replace

func getFormParameter(r *http.Request, param string) string {
	if len(r.Form[param]) > 0 {
		return r.Form[param][0]
	}
	return ""
}

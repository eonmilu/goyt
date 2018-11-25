package goyt

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// CreateVerifToken gets a channel ID, generates a random string of length 32,
// and adds it to a DB table, leaving 2 hours before expiration
func (y YourTime) CreateVerifToken(w http.ResponseWriter, r *http.Request) {
	channelID := getFormParameter(r, "channelid")
	if channelID == "" {
		fmt.Fprintf(w, sCError)
		return
	}
	userExists, err := y.userSecretExists(channelID)
	if err != nil {
		fmt.Printf("%s", err)
		fmt.Fprintf(w, sCError)
		return
	}

	secret := randStringBytes(randStringLength)
	if userExists {
		_, err := y.DB.Exec("UPDATE tokens SET secret = $1 WHERE channelid = $2", secret, channelID)
		if err != nil {
			fmt.Printf("%s", err)
			fmt.Fprintf(w, sCError)
			return
		}
	} else {
		_, err := y.DB.Exec("INSERT INTO tokens VALUES ($1, $2)", channelID, secret)
		if err != nil {
			fmt.Printf("%s", err)
			fmt.Fprintf(w, sCError)
			return
		}
	}

	fmt.Fprintf(w, secret)
}

func (y YourTime) userSecretExists(channelid string) (bool, error) {
	result := false
	row := y.DB.QueryRow("SELECT exists(SELECT 1 FROM tokens WHERE channelid=$1)", channelid)
	err := row.Scan(&result)
	return result, err
}

func (y YourTime) getVerifSecretFromDB(channelid string) (string, error) {
	var secret string
	row := y.DB.QueryRow("SELECT secret FROM tokens WHERE channelid=$1", channelid)
	err := row.Scan(&secret)

	return secret, err
}

// https://stackoverflow.com/a/31832326/8774937
const randStringLength = 32
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func randStringBytes(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// GetRandString writes to r a random string of length 32
func GetRandString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, randStringBytes(randStringLength))
}

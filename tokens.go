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
	}
	// TODO: create token table
	// TODO: add ids and token to table
	fmt.Fprintf(w, randStringBytes(randStringLength))
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

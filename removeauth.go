package goyt

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// RemoveAuth reads a token passed by HTTPS POST and changes the DB's entry
// to an empty string
func (y YourTime) RemoveAuth(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)

	token := getTokenFromCookies(r)

	_, err := y.DB.Exec("UPDATE users SET token=$1 WHERE token=$2", "", string(token))
	if err != nil {
		log.Printf("%s", err)
		fmt.Fprintf(w, sCError)
		return
	}

	cookie := http.Cookie{
		Name:    "yourtime-token-server",
		Value:   "",
		Expires: time.Unix(0, 0),
		Secure:  true,
	}
	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, sCOK)
}

func getTokenFromCookies(r *http.Request) token {
	cookies := r.Header.Get("Cookie")
	re := regexp.MustCompile(`(?m)yourtime-token-server=.*[^\]|;]`)
	tkn := strings.Split(re.FindAllString(cookies, 1)[0], "=")[1]

	return token(tkn)
}

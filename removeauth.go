package goyt

import (
	"fmt"
	"log"
	"net/http"
)

// RemoveAuth reads a token passed by HTTPS POST and changes the DB's entry
// to an empty string
func (y YourTime) RemoveAuth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token := getToken(r)

	_, err := y.DB.Exec("UPDATE users SET token=$1 WHERE token=$2", "", token)
	if err != nil {
		log.Printf("%s", err)
		fmt.Fprintf(w, sCError)
		return
	}
}

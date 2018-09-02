package goyt

import (
	"fmt"
	"net/http"
)

// EnableCORS edits the response headers to allow CORS from any source
func EnableCORS(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

// Middleware tries to identify the user. If it fails, it creates a new one
func (y YourTime) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userExists, err := timemarksDB{y.DB}.userExistsByToken(getTokenFromCookies(r))
		if err != nil {
			fmt.Fprintf(w, sCError)
			fmt.Printf("%s", err)
			return
		}
		if userExists {
			return
		}

		userExists, err = timemarksDB{y.DB}.userExistsByIdentifier(r.RemoteAddr)
		if err != nil {
			fmt.Fprintf(w, sCError)
			fmt.Printf("%s", err)
			return
		}
		if userExists {
			return
		}

		// If there is no record of the user, create one
		err = y.handleNewUser(User{
			Identifier: r.RemoteAddr,
		})
		if err != nil {
			fmt.Fprintf(w, sCError)
			fmt.Printf("%s", err)
			return
		}

		h.ServeHTTP(w, r)
	})
}

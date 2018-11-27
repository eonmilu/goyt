package goyt

import (
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// EnableCORS edits the response headers to allow CORS from any source
func EnableCORS(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Since a wildcard Access-Control-Allow-Origin header is not allowed
	// for jQuery AJAX requests with cookies, allow only some domains.
	w.Header().Set("Access-Control-Allow-Origin", "https://www.youtube.com")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// CreateUsers tries to identify the user. If it fails, it creates a new one
func (y YourTime) CreateUsers(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tkn := getTokenFromCookies(r)
		if tkn != "" {
			userExists, err := y.userExistsByToken(tkn)
			if err != nil {
				fmt.Fprintf(w, sCError)
				fmt.Printf("%s", err)
				return
			}
			if userExists {
				// Pass to the next function
				h(w, r)
				return
			}
		}

		trueAddr := strings.Split(r.RemoteAddr, ":")[0]
		userExists, err := y.userExistsByIdentifier(trueAddr)
		if err != nil {
			fmt.Fprintf(w, sCError)
			fmt.Printf("%s", err)
			return
		}
		if userExists {
			// Pass to the next function
			h(w, r)

			return
		}

		// If there is no record of the user, create one
		err = y.handleNewUser(User{Identifier: trueAddr})

		if err != nil {
			fmt.Fprintf(w, sCError)
			fmt.Printf("%s", err)
			return
		}

		// Pass to the next function
		h(w, r)
	})
}

func parseAndValidateToken(tokenString, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})
	return token, err
}

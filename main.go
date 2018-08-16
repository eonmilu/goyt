package goyt

import (
	"database/sql"

	// Init the postgres' drivers
	_ "github.com/lib/pq"
)

const (
	sCFound    = "200"
	sCNotFound = "210"
	sCError    = "220"
)

const (
	// tmrkDefOffset : the default offset for a timemark lookup
	tmrkDefOffset = "0"
	// tmrkDefLimit : the default limit for a timemark lookup
	tmrkDefLimit = "10"
)

// YourTime is a struct containing the methods and variables
type YourTime struct {
	// AuthTokenURL is the URL which Google provides to validate a token
	AuthTokenURL string
	// GoogleClientID is the URL the request is coming from
	GoogleClientID string
	// DB is the database where the Timemarks and users are stored
	DB *sql.DB
}

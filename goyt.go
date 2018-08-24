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
	sCBadLogin = "230"
	sCOK       = "240"
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

// Timemark is the data structure for a timemark
type Timemark struct {
	ID       int64  `json:"timemarkID"`
	Author   int64  `json:"author"`
	Timemark int64  `json:"timemark"`
	Content  string `json:"content"`
	Votes    int64  `json:"votes"`
	Date     int64  `json:"date"`
	Approved bool   `json:"approved"`
}

// Author is the type used to read info on an author from the database
type Author struct {
	ID       int64
	username string
	url      string
}

// User contains information on a user
type User struct {
	id    int32
	token string
	// Auth is the unmarshaled data structure AuthTokenURL returns
	// We will only declare the types we are going to use (aud, sub, email)
	Aud   string `json:"aud"`
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

type timemarksDB struct {
	*sql.DB
}

package goyt

import (
	"database/sql"
	"fmt"
	"time"

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
	// DB is the database where the Timemarks and users are stored
	DB *sql.DB
	// JWTSecret is the secret that will be used to sign and validate JWT
	JWTSecret string
}

// Timemark is the data structure for a timemark
type Timemark struct {
	ID        int64   `json:"id"`
	VideoID   string  `json:"videoid"`
	IP        string  `json:"ip"`
	Timemark  int64   `json:"timemark"`
	Content   string  `json:"content"`
	Votes     int64   `json:"votes"`
	Author    int64   `json:"author"`
	Approved  bool    `json:"approved"`
	Timestamp string  `json:"timestamp"`
	Reports   []int64 `json:"reports"`
}

// Author is the type used to read info on an author from the database
type Author struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	URL      string `json:"url"`
	Valid    bool   `json:"valid"`
}

// User contains information on a user or channel
type User struct {
	id int32
	// Auth is the unmarshaled data structure AuthTokenURL returns
	Identifier  string `json:"identifier"`
	Username    string `json:"username"`
	URL         string `json:"url"` // TODO: ask user for youtube id
	Picture     string `json:"picture"`
	Description string `json:"description"`
}

type youTubeChannelResponse struct {
	Metadata struct {
		ChannelMetadataRenderer struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			ExternalID  string `json:"externalId"`
			Avatar      struct {
				Thumbnails []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"avatar"`
			ChannelURL   string `json:"channelUrl"`
			IsFamilySafe bool   `json:"isFamilySafe"`
		} `json:"channelMetadataRenderer"`
	} `json:"metadata"`
}

func (y YourTime) removeOldVerifTokens() {
	for {
		_, err := y.DB.Exec("DELETE FROM tokens WHERE age(now(), ts) > INTERVAL '1 day';")
		if err != nil {
			fmt.Printf("%s", err)
		}

		time.Sleep(time.Second)
	}
}

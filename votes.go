package goyt

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// Votes reads the parameters supplied by HTTPS POST (id, action)
// and then modifies the timemarks' votes given the action
func (y YourTime) Votes(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	params, err := getVoteParameters(r)
	if err != nil {
		fmt.Fprintf(w, sCError)
		fmt.Printf("%s", err)
		return
	}

}

func getVoteParameters(r *http.Request) (vote, error) {
	v := vote{}

	id, err := getTimemarkIDFromPost(r)
	if err != nil {
		return v, err
	}

	action, err := getActionFromPost(r)
	if err != nil {
		return v, err
	}

	v = vote{
		ID:     id,
		Action: action,
	}

	return v, nil
}

const (
	actionUpvote   = "upvote"
	actionUnset    = "unset"
	actionDownvote = "downvote"
)

type vote struct {
	ID     int64
	Action string
}

func getTimemarkIDFromPost(r *http.Request) (int64, error) {
	rawID := r.Form["id"]
	if len(rawID) > 0 {
		id, err := strconv.ParseInt(rawID[0], 10, 64)
		return int64(id), err
	}
	return 0, errors.New("There was no id parameter")
}

func getActionFromPost(r *http.Request) (string, error) {
	if len(r.Form["action"]) > 0 {
		return r.Form["action"][0], nil
	}
	return "", errors.New("There was no action parameter")
}

package goyt

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	actionUpvote   = "upvoted"
	actionUnset    = "unset"
	actionDownvote = "downvoted"
)

// Votes reads the parameters supplied by HTTPS POST (id, action)
// and then modifies the timemarks' votes given the action
func (y YourTime) Votes(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)
	r.ParseForm()

	params, err := y.getVoteParameters(r)
	if err != nil {
		fmt.Fprintf(w, sCError)
		fmt.Printf("%s", err)
		return
	}

	err = y.handleVoteAction(params)
	if err != nil {
		fmt.Fprintf(w, sCError)
		fmt.Printf("%s", err)
		return
	}
	fmt.Fprintf(w, sCOK)
}

func (y YourTime) handleVoteAction(v vote) error {
	switch v.Action {
	case actionUpvote:
		hasUpvoted, err := y.hasUpvoted(v)
		if err != nil {
			return err
		}
		if !hasUpvoted {
			err := y.upvote(v)
			return err
		}
		break
	case actionDownvote:
		hasDownvoted, err := y.hasDownvoted(v)
		if err != nil {
			return err
		}
		if !hasDownvoted {
			err := y.downvote(v)
			return err
		}
		break
	case actionUnset:
		err := y.unsetVote(v)
		return err
	}
	return nil
}

func (y YourTime) upvote(v vote) error {
	return y.vote(v, "upvotes")
}

func (y YourTime) downvote(v vote) error {
	return y.vote(v, "downvotes")
}

func (y YourTime) vote(v vote, voteCollection string) error {
	var (
		voteNumber int8
	)

	switch voteCollection {
	case "upvotes":
		voteNumber = 1
	case "downvotes":
		voteNumber = -1
	default:
		return errors.New("Invalid collection string")
	}

	// unset any votes
	y.unsetVote(v)

	// TODO: dangerous unsanitized integer in statement.
	// This is because lib/pq does not support arrays in statements
	stmt := fmt.Sprintf("UPDATE users SET %s= %s || '{%d}' where identifier=$1", voteCollection, voteCollection, v.ID)

	// Set the action in the user's profile
	_, err := y.DB.Exec(stmt, v.Identifier)
	if err != nil {
		return err
	}

	// Change the timemark's votes
	stmt = fmt.Sprintf("UPDATE timemarks SET votes= votes + %d where id=$1", voteNumber)
	_, err = y.DB.Exec(stmt, v.ID)

	return err
}

func (y YourTime) unsetVote(v vote) error {
	// Remove the upvote in the user's profile
	unsetUpvote, err := y.hasUpvoted(v)
	if err != nil {
		return err
	}
	if unsetUpvote {
		_, err := y.DB.Exec("UPDATE users SET upvotes=array_remove(upvotes, $1)", v.ID)

		if err != nil {
			return err
		}
		// Change the timemark's votes
		_, err = y.DB.Exec("UPDATE timemarks SET votes= votes - 1 where id=$1", v.ID)

	}

	// Remove the downvote in the user's profile
	unsetDownvote, err := y.hasDownvoted(v)
	if err != nil {
		return err
	}
	if unsetDownvote {
		_, err := y.DB.Exec("UPDATE users SET downvotes=array_remove(downvotes, $1)", v.ID)

		if err != nil {
			return err
		}
		// Change the timemark's votes
		_, err = y.DB.Exec("UPDATE timemarks SET votes= votes + 1 where id=$1", v.ID)
	}
	return nil
}

func (y YourTime) hasUpvoted(v vote) (bool, error) {
	return y.hasVoted(v, "upvotes")
}

func (y YourTime) hasDownvoted(v vote) (bool, error) {
	return y.hasVoted(v, "downvotes")
}

func (y YourTime) hasVoted(v vote, voteCollection string) (bool, error) {
	// TODO: dangerous unsanitized integer in statement.
	// This is because lib/pq does not support arrays in statements
	stmt := fmt.Sprintf("SELECT '{%d}' && (select %s from users where identifier=$1)", v.ID, voteCollection)

	row := y.DB.QueryRow(stmt, v.Identifier)

	var isVoted sql.NullBool
	err := row.Scan(&isVoted)

	if isVoted.Valid {
		return isVoted.Bool, err
	}
	return false, err
}

func (y YourTime) getVoteParameters(r *http.Request) (vote, error) {
	v := vote{}

	id, err := getTimemarkIDFromPost(r)
	if err != nil {
		return v, err
	}

	action, err := getActionFromPost(r)
	if err != nil {
		return v, err
	}

	identifier, err := y.getUserIdentifier(r)
	if err != nil {
		return v, err
	}

	v = vote{
		ID:         id,
		Action:     action,
		Identifier: identifier,
	}

	return v, nil
}

func (y YourTime) getUserIdentifier(r *http.Request) (string, error) {
	pureAddr := strings.Split(r.RemoteAddr, ":")[0]
	tkn := getTokenFromCookies(r)
	if tkn == "" {
		return pureAddr, nil
	}

	email, err := y.getEmailFromToken(tkn)
	if err != nil {
		return pureAddr, err
	}
	return email, nil
}

func (y YourTime) getEmailFromToken(tkn token) (string, error) {
	email := ""
	row := y.DB.QueryRow("SELECT identifier FROM users WHERE token=$1", string(tkn))
	err := row.Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}

type vote struct {
	ID         int64
	Action     string
	Identifier string
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

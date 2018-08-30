package goyt

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	actionUpvote   = "upvoted"
	actionUnset    = "unset"
	actionDownvote = "downvoted"
)

// Votes reads the parameters supplied by HTTPS POST (id, action)
// and then modifies the timemarks' votes given the action
func (y YourTime) Votes(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	params, err := y.getVoteParameters(r)
	log.Printf("%v", params)
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
		log.Println(hasUpvoted)
		if !hasUpvoted {
			err := y.upvote(v)
			if err != nil {
				return err
			}
		}
		break
	case actionDownvote:
		hasDownvoted, err := y.hasDownvoted(v)
		if err != nil {
			return err
		}
		if !hasDownvoted {
			err := y.downvote(v)
			if err != nil {
				return err
			}
		}
		break
	case actionUnset:
		err := y.unsetVote(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (y YourTime) upvote(v vote) error {
	// Set the upvote in the user's profile
	_, err := y.DB.Exec("UPDATE users SET upvotes= upvotes || '{"+strconv.FormatInt(v.ID, 10)+"}' where identifier=$1", v.Identifier)
	log.Println("1")
	if err != nil {
		return err
	}
	// Remove possible downvote
	_, err = y.DB.Exec("UPDATE users SET downvotes= array_remove(downvotes, $1) where identifier=$2", v.ID, v.Identifier)
	log.Println("2")

	if err != nil {
		return err
	}
	// Change the timemark's votes
	_, err = y.DB.Exec("UPDATE timemarks SET votes= votes + 1 where id=$1", v.ID)
	log.Println("3")

	return err
}

func (y YourTime) downvote(v vote) error {
	// Set the downvote in the user's profile
	_, err := y.DB.Exec("UPDATE users SET downvotes= downvotes || '{"+strconv.FormatInt(v.ID, 10)+"}' where identifier=$1", v.Identifier)
	log.Println("4")

	if err != nil {
		return err
	}
	// Remove possible upvote
	_, err = y.DB.Exec("UPDATE users SET upvotes= array_remove(upvotes, $1) where identifier=$2", v.ID, v.Identifier)
	log.Println("5")

	if err != nil {
		return err
	}
	// Change the timemark's votes
	_, err = y.DB.Exec("UPDATE timemarks SET votes= votes - 1 where id=$1", v.ID)
	log.Println("6")

	return err
}

func (y YourTime) unsetVote(v vote) error {
	// Remove the upvote in the user's profile
	unsetUpvote, err := y.hasUpvoted(v)
	if err != nil {
		return err
	}
	// Remove the downvote in the user's profile
	unsetDownvote, err := y.hasDownvoted(v)
	if err != nil {
		return err
	}
	if unsetUpvote {
		_, err := y.DB.Exec("UPDATE users SET upvotes=array_remove(upvotes, $1)", v.ID)
		log.Println("7")

		if err != nil {
			return err
		}
		// Change the timemark's votes
		_, err = y.DB.Exec("UPDATE timemarks SET votes= votes - 1 where id=$1", v.ID)
		log.Println("8")

	} else if unsetDownvote {
		_, err := y.DB.Exec("UPDATE users SET downvotes=array_remove(downvotes, $1)", v.ID)
		log.Println("9")

		if err != nil {
			return err
		}
		// Change the timemark's votes
		_, err = y.DB.Exec("UPDATE timemarks SET votes= votes + 1 where id=$1", v.ID)
		log.Println("10")

	}
	return nil
}

func (y YourTime) hasUpvoted(v vote) (bool, error) {
	// TODO: dangerous unsanitized integer in statement.
	// This is because lib/pq does not support arrays in statements
	row := y.DB.QueryRow("SELECT '{"+strconv.FormatInt(v.ID, 10)+"}'= ANY(select upvotes from users where identifier=$1)", v.Identifier)
	log.Println("11")

	log.Println("011")

	isUpvoted := false
	err := row.Scan(&isUpvoted)
	log.Printf("%t", isUpvoted)
	return isUpvoted, err
}

func (y YourTime) hasDownvoted(v vote) (bool, error) {
	row := y.DB.QueryRow("SELECT '{"+strconv.FormatInt(v.ID, 10)+"}' = ANY(select downvotes from users where identifier=$1)", v.Identifier)
	log.Println("12")

	log.Println("012")
	isDownvoted := false
	err := row.Scan(&isDownvoted)
	return isDownvoted, err
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
	tkn := getTokenFromCookies(r)
	if tkn == "" {
		return r.RemoteAddr, nil
	}

	email, err := y.getEmailFromToken(tkn)
	if err != nil {
		return r.RemoteAddr, err
	}
	return email, nil
}

func (y YourTime) getEmailFromToken(tkn token) (string, error) {
	email := ""
	row := y.DB.QueryRow("SELECT email FROM users WHERE token=$1", string(tkn))
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

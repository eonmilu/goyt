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
		if hasUpvoted {
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
		if hasDownvoted {
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
	_, err := y.DB.Exec("UPDATE users SET upvotes= upvotes || '{$1}' where identifier=$2", v.ID, v.Identifier)
	if err != nil {
		return err
	}
	// Remove possible downvote
	_, err = y.DB.Exec("UPDATE users SET downvotes= array_remove(downvotes, $1) where identifier=$2", v.ID, v.Identifier)
	if err != nil {
		return err
	}
	// Change the timemark's votes
	_, err = y.DB.Exec("UPDATE timemarks SET votes= votes + 1 where id=$1", v.ID)
	return err
}

func (y YourTime) downvote(v vote) error {
	// Set the downvote in the user's profile
	_, err := y.DB.Exec("UPDATE users SET downvotes= downvotes || '{$1}' where identifier=$2", v.ID, v.Identifier)
	if err != nil {
		return err
	}
	// Remove possible upvote
	_, err = y.DB.Exec("UPDATE users SET upvotes= array_remove(upvotes, $1) where identifier=$2", v.ID, v.Identifier)
	if err != nil {
		return err
	}
	// Change the timemark's votes
	_, err = y.DB.Exec("UPDATE timemarks SET votes= votes - 1 where id=$1", v.ID)
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
		if err != nil {
			return err
		}
		// Change the timemark's votes
		_, err = y.DB.Exec("UPDATE timemarks SET votes= votes - 1 where id=$1", v.ID)
	} else if unsetDownvote {
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
	rows, err := y.DB.Query("SELECT '{$1}' = ANY(select upvotes from users where identifier=$2)", v.ID, v.Identifier)
	if err != nil {
		return false, err
	}
	isUpvoted := false
	rows.Next()
	err = rows.Scan(&isUpvoted)
	if err != nil {
		return false, err
	}
	if rows.Err() != nil {
		return false, rows.Err()
	}
	return isUpvoted, nil
}

func (y YourTime) hasDownvoted(v vote) (bool, error) {
	rows, err := y.DB.Query("SELECT '{$1}' = ANY(select downvotes from users where identifier=$2)", v.ID, v.Identifier)
	if err != nil {
		return false, err
	}
	isDownvoted := false
	rows.Next()
	err = rows.Scan(&isDownvoted)
	if err != nil {
		return false, err
	}
	if rows.Err() != nil {
		return false, rows.Err()
	}
	return isDownvoted, nil
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

	identifier := getUserIdentifier(r)

	v = vote{
		ID:         id,
		Action:     action,
		Identifier: identifier,
	}

	return v, nil
}

func getUserIdentifier(r *http.Request) string {
	identifier := string(getTokenFromCookies(r))
	if identifier == "" {
		identifier = r.RemoteAddr
		return identifier
	}
	return identifier
}

const (
	actionUpvote   = "upvote"
	actionUnset    = "unset"
	actionDownvote = "downvote"
)

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

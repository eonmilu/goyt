package goyt

import (
	"net/http"
)

// Votes reads the parameters supplied by HTTPS POST (id, action)
// and then modifies the timemarks' votes given the action
func (y YourTime) Votes(w http.ResponseWriter, r *http.Request) {

}

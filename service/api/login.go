package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// return the user profile given its id, the user profile is composed by the user’s photos (in reverse chronological order), how many photos have been uploaded, and the user’s followers and following
func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var username string
	var err error
	err = json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.CheckUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if count == 0 {
		err = rt.db.DoLogin(username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var id int
		err, id = rt.db.FindUserId(username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Printf("%d", id)
		return
	}
}

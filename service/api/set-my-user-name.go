package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// update a user's username given their id and new username, returns all User's informations
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var user User
	var id int
	id, err := strconv.Atoi(ps.ByName("idUser"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.IdUser = id

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.CheckUsername(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if count == 0 {
		err = rt.db.SetUsername(user.IdUser, user.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var id int
		var usrName string
		var bio string
		err, id, usrName, bio = rt.db.SelectUser(user.IdUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(UserToJson(id, usrName, bio))

	} else {
		w.WriteHeader(409) //http.Conflict
		return
	}
}
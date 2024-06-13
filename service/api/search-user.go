package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

/*
liking an existing photo using data provided in the body of request and the photo id and the owner of the photo id,
return the full image object
*/
func (rt *_router) searchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User
	var token int
	var err error

	token, err = strconv.Atoi(extractBearer(r.Header.Get("Authorization")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 401 - you must be logged in, not banned
	if isNotLogged(token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user.Username = ps.ByName("username")

	var tmpId int

	err, tmpId = rt.db.FindUserId(user.Username)
	if err != nil {
		// 404 - username not found
		if err.Error() == "sql: no rows in result set" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			err = json.NewEncoder(w).Encode(user.Username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.IdUser = tmpId

	// 403 - you can't search yourself
	if user.IdUser == token {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	var count int

	err, count = rt.db.CheckBan(user.IdUser, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you have been banned from the user searched
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, count = rt.db.CheckBan(token, user.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you have banned this user
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	// 200 - username found
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

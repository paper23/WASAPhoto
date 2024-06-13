package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

/*
following an existing user using data provided in the body of the request and
the user id in the path
*/
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var idUserToFollow int
	var err error
	var token int

	token, err = strconv.Atoi(extractBearer(r.Header.Get("Authorization")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 401 - you must be logged in, not followed
	if isNotLogged(token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	idUserToFollow, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.FindUserById(idUserToFollow)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 404 - user to follow not found, not followed
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(idUserToFollow)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	// 403 - you can't follow yourself, not followed
	if token == idUserToFollow {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(idUserToFollow)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, count = rt.db.CheckBan(token, idUserToFollow)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you have banned the user, not followed
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(idUserToFollow)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, count = rt.db.CheckBan(idUserToFollow, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you have been banned from the user, not followed
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(idUserToFollow)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err = rt.db.FollowUser(token, idUserToFollow)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 200 - user succesfully followed
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return

}

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

	var user User
	var userToFollow User
	var err error
	var token int

	token, err = strconv.Atoi(extractBearer(r.Header.Get("Authorization")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//401 - you must be logged in, not followed
	if isNotLogged(token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user.IdUser, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.FindUserById(user.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//404 - user (follower) not found, not followed
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(user.IdUser)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&userToFollow)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, count = rt.db.FindUserById(userToFollow.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//404 - user (to follow) not found, not followed
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(userToFollow.IdUser)
		return
	}

	//403 - you cannot follow an user for another user, not followed
	if token != user.IdUser {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(user.IdUser)
		json.NewEncoder(w).Encode(userToFollow.IdUser)
		json.NewEncoder(w).Encode(token)
		return
	}

	//403 - you can't follow yourself, not followed
	if user.IdUser == userToFollow.IdUser {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(user.IdUser)
		json.NewEncoder(w).Encode(userToFollow.IdUser)
		json.NewEncoder(w).Encode(token)
		return
	}

	count = 0
	var count2 int
	err, count = rt.db.CheckBan(user.IdUser, userToFollow.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, count2 = rt.db.CheckBan(userToFollow.IdUser, user.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//403 - you have been banned from the user and/or you have banned the user, not followed
	if count+count2 > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(user.IdUser)
		json.NewEncoder(w).Encode(userToFollow.IdUser)
		json.NewEncoder(w).Encode(token)
		return
	}

	err = rt.db.FollowUser(user.IdUser, userToFollow.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, userToFollow.IdUser, userToFollow.Username, userToFollow.Biography = rt.db.SelectUser(userToFollow.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//200 - user succesfully followed
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userToFollow)

}

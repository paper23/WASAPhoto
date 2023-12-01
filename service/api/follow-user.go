package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// following an existing user using data provided in the body of the request and the user id in the path
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User
	var userToFollow User
	var err error

	user.IdUser, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&userToFollow)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.FindUserById(userToFollow.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if count <= 0 {
		var follow DoubleIdUser
		follow.IdUser = user.IdUser
		follow.IdUser2 = userToFollow.IdUser
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(follow)
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

	if count+count2 > 0 {
		var follow DoubleIdUser
		follow.IdUser = user.IdUser
		follow.IdUser2 = userToFollow.IdUser
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(follow)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(userToFollow)

}

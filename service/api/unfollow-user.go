package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// following an existing user using data provided in the body of the request and the user id in the path
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User
	var userToUnfollow User
	var err error

	user.IdUser, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userToUnfollow.IdUser, err = strconv.Atoi(ps.ByName("idUserToUnfollow"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.IdUser == userToUnfollow.IdUser {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var count int
	err, count = rt.db.FindUserById(userToUnfollow.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if count <= 0 {
		var follow DoubleIdUser
		follow.IdUser = user.IdUser
		follow.IdUser2 = userToUnfollow.IdUser
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(follow)
		return
	}

	count = 0
	err, count = rt.db.CheckFollowing(user.IdUser, userToUnfollow.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if count <= 0 {
		var follow DoubleIdUser
		follow.IdUser = user.IdUser
		follow.IdUser2 = userToUnfollow.IdUser
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(follow)
		return
	}

	err = rt.db.UnfollowUser(user.IdUser, userToUnfollow.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, userToUnfollow.IdUser, userToUnfollow.Username, userToUnfollow.Biography = rt.db.SelectUser(userToUnfollow.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(userToUnfollow)

}

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// unfollow an existing and followed user given its id
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User
	var userToUnfollow User
	var err error
	var token int

	token, err = strconv.Atoi(extractBearer(r.Header.Get("Authorization")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 401 - you must be logged in, not unfollowed
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
	// 404 - user (follower) not found, not unfollowed
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(user.IdUser)
		return
	}

	userToUnfollow.IdUser, err = strconv.Atoi(ps.ByName("idUserToUnfollow"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, count = rt.db.FindUserById(userToUnfollow.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 404 - user (to unfollow) not found, not unfollowed
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(userToUnfollow.IdUser)
		return
	}

	// 403 - you cannot unfollow an user for another user, not unfollowed
	if token != user.IdUser {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(user.IdUser)
		json.NewEncoder(w).Encode(token)
		return
	}

	// 403 - you can't unfollow yourself, not unfollowed
	if user.IdUser == userToUnfollow.IdUser {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(user.IdUser)
		json.NewEncoder(w).Encode(userToUnfollow.IdUser)
		json.NewEncoder(w).Encode(token)
		return
	}

	count = 0
	err, count = rt.db.CheckFollowing(user.IdUser, userToUnfollow.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you can't unfollow a user you don't yet follow, not unfollowed
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(user.IdUser)
		json.NewEncoder(w).Encode(userToUnfollow.IdUser)
		json.NewEncoder(w).Encode(token)
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

	// 200 - user succesfully unfollowed
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userToUnfollow)

}

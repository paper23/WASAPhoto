package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// unban an existing (and already banned) user given its id
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User
	var userToSban User
	var err error
	var token int

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

	user.IdUser, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.FindUserById(userToSban.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 404 - user (sbanner) not found, not sbanned
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(user.IdUser)
		return
	}

	userToSban.IdUser, err = strconv.Atoi(ps.ByName("idUserBanned"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, count = rt.db.FindUserById(userToSban.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 404 - user (to sban) not found, not sbanned
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(userToSban.IdUser)
		return
	}

	// 403 - you can't sban an user for another user, not sbanned
	if token != user.IdUser {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(user.IdUser)
		json.NewEncoder(w).Encode(userToSban.IdUser)
		json.NewEncoder(w).Encode(token)
		return
	}

	// 403 - you can't sban yourself, not sbanned
	if user.IdUser == userToSban.IdUser {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(user.IdUser)
		json.NewEncoder(w).Encode(userToSban.IdUser)
		json.NewEncoder(w).Encode(token)
		return
	}

	err, count = rt.db.CheckBan(user.IdUser, userToSban.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 404 - user not banned, non sbanned
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(user.IdUser)
		json.NewEncoder(w).Encode(userToSban.IdUser)
		return
	}

	err = rt.db.SbanUser(user.IdUser, userToSban.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, userToSban.IdUser, userToSban.Username, userToSban.Biography = rt.db.SelectUser(userToSban.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 200 - user sbanned succesfully
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userToSban)

}

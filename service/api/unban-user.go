package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// following an existing user using data provided in the body of the request and the user id in the path
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User
	var userToSban User
	var err error

	user.IdUser, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userToSban.IdUser, err = strconv.Atoi(ps.ByName("idUserBanned"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.IdUser == userToSban.IdUser {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var count int
	err, count = rt.db.FindUserById(userToSban.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if count <= 0 {
		var sban DoubleIdUser
		sban.IdUser = user.IdUser
		sban.IdUser2 = userToSban.IdUser
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(sban)
		return
	}

	count = 0
	err, count = rt.db.CheckBan(user.IdUser, userToSban.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if count <= 0 {
		var ban DoubleIdUser
		ban.IdUser = user.IdUser
		ban.IdUser2 = userToSban.IdUser
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ban)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(userToSban)

}

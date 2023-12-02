package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// following an existing user using data provided in the body of the request and the user id in the path
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var user User
	var userToBan User
	var err error

	user.IdUser, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&userToBan)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.IdUser == userToBan.IdUser {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var count int
	err, count = rt.db.FindUserById(userToBan.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if count <= 0 {
		var ban DoubleIdUser
		ban.IdUser = user.IdUser
		ban.IdUser2 = userToBan.IdUser
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ban)
		return
	}

	count = 0
	err, count = rt.db.CheckBan(user.IdUser, userToBan.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if count > 0 {
		var ban DoubleIdUser
		ban.IdUser = user.IdUser
		ban.IdUser2 = userToBan.IdUser
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ban)
		return
	}
	//QUANDO SI BANNA QUALCUNO BISOGNA TOGLIERE ANCHE FOLLOWS (FATTO) LIKE E COMMENTI!
	err = rt.db.BanUser(user.IdUser, userToBan.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, userToBan.IdUser, userToBan.Username, userToBan.Biography = rt.db.SelectUser(userToBan.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(userToBan)

}

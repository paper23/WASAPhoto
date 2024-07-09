package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// returns the list of users banned by the user who made the request
func (rt *_router) getBannedList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var err error
	var token int
	var idUser int
	var bannedList []database.User

	token, err = strconv.Atoi(extractBearer(r.Header.Get("Authorization")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 401 - you must be logged in, photo not uploaded
	if isNotLogged(token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	idUser, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you can't see another user's banned list
	if idUser != token {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(idUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	var count int
	err, count = rt.db.FindUserById(idUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 404 - user not found
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(idUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, bannedList = rt.db.GetBanned(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if bannedList == nil {
		bannedList = []database.User{}
		bannedList = append(bannedList, database.User{IdUser: 0, Username: ""})
	}

	// 200 - get banned list
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(bannedList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// unfollow an existing and followed user given its id
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var idUserToUnfollow int
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

	idUserToUnfollow, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int

	err, count = rt.db.FindUserById(idUserToUnfollow)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 404 - user to unfollow not found, not unfollowed
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(idUserToUnfollow)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	// 403 - you cannot unfollow yourself, not unfollowed
	if token == idUserToUnfollow {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(idUserToUnfollow)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, count = rt.db.CheckBan(token, idUserToUnfollow)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you can't unfollow an user you have banned
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(idUserToUnfollow)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, count = rt.db.CheckBan(idUserToUnfollow, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you can't unfollow an user who has banned you
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(idUserToUnfollow)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, count = rt.db.CheckFollowing(token, idUserToUnfollow)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you can't unfollow an user you don't follow yet, not unfollowed
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(idUserToUnfollow)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err = rt.db.UnfollowUser(token, idUserToUnfollow)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 200 - user succesfully unfollowed
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return

}

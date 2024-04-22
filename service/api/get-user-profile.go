package api

import (
	//"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

/*
return the user profile given its id, the user profile is composed
by the user’s photos (in reverse chronological order), how many photos
have been uploaded, and the user’s followers and following
*/
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var profile UserProfile
	var token int
	var err error

	token, err = strconv.Atoi(extractBearer(r.Header.Get("Authorization")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 401 - you must be logged in
	if isNotLogged(token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	profile.IdUser, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, profile.FollowCount = rt.db.CountFollowing(profile.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, profile.FollowerCount = rt.db.CountFollowers(profile.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, profile.Images = rt.db.GetUserImages(profile.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ho recuperato dal DB tutte le informazioni che mi servono
	// ora devo fare tutti i possibili errori come riportano le API

}

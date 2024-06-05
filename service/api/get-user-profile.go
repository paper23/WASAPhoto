package api

import (
	"encoding/json"
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

	profile.User.IdUser, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.FindUserById(profile.User.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 404 - user not found
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(profile.User.IdUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, count = rt.db.CheckBan(token, profile.User.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you cannot see the profile of a user who has banned you
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(profile.User.IdUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, count = rt.db.CheckBan(profile.User.IdUser, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you cannot see the profile of a user you banned
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(profile.User.IdUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, profile.FollowCount = rt.db.CountFollowing(profile.User.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, profile.FollowerCount = rt.db.CountFollowers(profile.User.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, profile.Images = rt.db.GetUserImagesId(profile.User.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, profile.ImgFiles = rt.db.GetUserImagesFile(profile.Images)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	print(profile.ImgFiles)

	err, profile.User.Biography = rt.db.FindUserBio(profile.User.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, profile.User.Username = rt.db.FindUsername(profile.User.IdUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 200 - profile succesfully showed
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(profile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

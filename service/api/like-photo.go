package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

/*
liking an existing photo using data provided in the body of request and the photo id and the owner of the photo id,
return the full image object
*/
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var like Liking
	var err error

	like.IdLiker, err = strconv.Atoi(extractBearer(r.Header.Get("Authorization")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 401 - you must be logged in, not banned
	if isNotLogged(like.IdLiker) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	like.IdOwner, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	like.IdImage, err = strconv.Atoi(ps.ByName("idImage"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.FindUserById(like.IdLiker)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 404 - user (liker) not found, not liked
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(like.IdLiker)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, count = rt.db.FindUserById(like.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 404 - user (owner) not found, not liked
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(like.IdOwner)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, count = rt.db.FindImage(like.IdImage, like.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 404 - image not found, not liked
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(like.IdImage)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(like.IdOwner)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, count = rt.db.CheckBan(like.IdOwner, like.IdLiker)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you can't like a photo of a user who has banned you, not liked
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(like)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, count = rt.db.CheckBan(like.IdLiker, like.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you can't like a photo of a user you have banned, not liked
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(like)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err, count = rt.db.CheckLike(like.IdLiker, like.IdImage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you can't like a photo you have already liked, not liked
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(like)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err = rt.db.LikePhoto(like.IdLiker, like.IdImage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 200 - photo succesfully liked
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(like)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

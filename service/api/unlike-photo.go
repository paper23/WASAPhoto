package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// remove a like given its id (the user liker id), the photo id and the owner of the photo id
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var like Liking
	var err error
	var token int

	token, err = strconv.Atoi(extractBearer(r.Header.Get("Authorization")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 401 - you must be logged in, not unliked
	if isNotLogged(token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	like.IdOwner, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.FindUserById(like.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 404 - user (owner) not found, not unliked
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(like.IdOwner)
		return
	}

	like.IdImage, err = strconv.Atoi(ps.ByName("idImage"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, count = rt.db.FindImage(like.IdImage, like.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 404 - photo not found, not unliked
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(like.IdImage)
		return
	}

	like.IdLiker, err = strconv.Atoi(ps.ByName("idLiker"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, count = rt.db.FindUserById(like.IdLiker)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 404 - user (liker) not found, not unliked
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(like.IdLiker)
		return
	}

	// 403 - you cannot unlike a photo for another user, not unliked
	if token != like.IdLiker {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(like.IdLiker)
		json.NewEncoder(w).Encode(token)
		return
	}

	err, count = rt.db.CheckBan(like.IdOwner, like.IdLiker)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 403 - you have been banned, not unliked
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(like)
		return
	}

	err, count = rt.db.CheckBan(like.IdLiker, like.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 403 - you have banned the user, not unliked
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(like)
		return
	}

	err, count = rt.db.CheckLike(like.IdLiker, like.IdImage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 403 - you can't unlike a photo you haven't liked yet, not unliked
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(like)
		return
	}

	err = rt.db.UnlikePhoto(like.IdLiker, like.IdImage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 200 - photo succesfully unliked
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(like)

}

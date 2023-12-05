package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// comment on a photo given its id, the id of the owner of the photo, the id of the user who comments on the photo and the text of the comment
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var like Liking
	var err error
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

	like.IdLiker, err = strconv.Atoi(ps.ByName("idLiker"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.FindImage(like.IdImage, like.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(like)
		return
	}

	err, count = rt.db.FindUserById(like.IdLiker)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(like)
		return
	}

	err, count = rt.db.CheckBan(like.IdOwner, like.IdLiker)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(like)

}

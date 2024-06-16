package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

/*
// get images of followed users in chronological order (from the most recent to the oldest)
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var token int
	var err error

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

	var idPath int

	idPath, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you can't see another user's stream
	if idPath != token {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var count int
	err, count = rt.db.FindUserById(idPath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 404 - user not found
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(idPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	var stream Stream

	err, stream.IdImage, stream.IdUser, stream.Username = rt.db.GetStream(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 200 - get my stream
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(stream)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
*/

// get images of followed users in chronological order (from the most recent to the oldest)
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	token, err := strconv.Atoi(extractBearer(r.Header.Get("Authorization")))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 401 - you must be logged in, not banned
	if isNotLogged(token) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var idPath int

	idPath, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you can't see another user's stream
	if idPath != token {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var count int
	err, count = rt.db.FindUserById(idPath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 404 - user not found
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(idPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	var stream []database.Stream2

	err, stream = rt.db.GetStream2(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for index, value := range stream {
		err, stream[index].Image.LikesCount = rt.db.CountLikes(value.Image.IdImage)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err, stream[index].Image.CommentsCount = rt.db.CountComments(value.Image.IdImage)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err, stream[index].Image.LikeStatus = rt.db.CheckLikeStatus(token, value.Image.IdImage)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err, stream[index].Image.Comments = rt.db.GetComments(value.Image.IdImage)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}

	// 200 - get my stream
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(stream)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

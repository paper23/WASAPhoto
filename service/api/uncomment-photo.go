package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// remove an existing comment given its id, the id of the owner of the photo, the id of the user who comments on the photo
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var comm Commenting
	var err error
	comm.IdOwner, err = strconv.Atoi(ps.ByName("idUser"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comm.IdImage, err = strconv.Atoi(ps.ByName("idImage"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comm.IdComment, err = strconv.Atoi(ps.ByName("idComment"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.FindImage(comm.IdImage, comm.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(comm)
		return
	}

	err, count = rt.db.FindComment(comm.IdComment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(comm)
		return
	}

	err, comm.Text = rt.db.SelectCommentText(comm.IdComment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	/*
		var check bool
		err, check = rt.db.CheckOwnership(comm.IdComment, comm.IdUserWriter)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !check {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(comm)
			return
		}
	*/

	err = rt.db.UncommentPhoto(comm.IdComment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(comm)

}

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// remove an existing comment given its id, the photo id and the owner of the photo id
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var comm Commenting
	var err error
	var token int

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

	comm.IdOwner, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.FindUserById(comm.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 404 - user not found, comment not deleted
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(comm.IdOwner)
		return
	}

	comm.IdImage, err = strconv.Atoi(ps.ByName("idImage"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, count = rt.db.FindImage(comm.IdImage, comm.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 404 - image not found, comment not deleted
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(comm.IdOwner)
		json.NewEncoder(w).Encode(comm.IdImage)
		return
	}

	comm.IdComment, err = strconv.Atoi(ps.ByName("idComment"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, count = rt.db.FindComment(comm.IdComment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 404 - comment not found, comment not deleted
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(comm.IdComment)
		return
	}

	err, comm.Text = rt.db.SelectCommentText(comm.IdComment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.UncommentPhoto(comm.IdComment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 200 - comment succesfully deleted
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comm)

}

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// comment on a photo given its id, the id of the owner of the photo, the id of the user who comments on the photo and the text of the comment
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	err = json.NewDecoder(r.Body).Decode(&comm)
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

	err, count = rt.db.FindUserById(comm.IdUserWriter)
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

	err, count = rt.db.CheckBan(comm.IdOwner, comm.IdUserWriter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(comm)
		return
	}

	err, count = rt.db.CheckBan(comm.IdUserWriter, comm.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if count > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(comm)
		return
	}

	err, comm.IdComment = rt.db.CommentPhoto(comm.IdUserWriter, comm.IdImage, comm.Text)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(comm)

}

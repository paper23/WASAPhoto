package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// remove a photo given its id and the id of the owner of the photo
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var img Image
	var err error
	img.IdOwner, err = strconv.Atoi(ps.ByName("idUser"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.FindUserById(img.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//404 - user not found, photo not deleted
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(img.IdOwner)
		return
	}

	img.IdImage, err = strconv.Atoi(ps.ByName("idImage"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, count = rt.db.FindImage(img.IdImage, img.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//404 - photo not found, photo not deleted
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(img.IdImage)
		return
	}

	err, count = rt.db.CheckPhotoOwnership(img.IdImage, img.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//401 - unauthorized, only the owner can delete a photo, photo not deleted
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(img)
		return
	}

	err, img.Url = rt.db.SelectImgUrl(img.IdImage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err, img.DateTime = rt.db.SelectImgDate(img.IdImage)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteImage(img.IdImage, img.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//200 - photo succesfully deleted
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(img)

}

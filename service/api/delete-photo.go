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

	img.IdImage, err = strconv.Atoi(ps.ByName("idImage"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var count int
	err, count = rt.db.FindImage(img.IdImage, img.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(img)

}
package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

/*
upload a new photo using data provided in the body of the request and
the id of the owner of the photo, return the full image object with the ID,
URL and date/time
*/
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var img Image
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

	img.IdOwner, err = strconv.Atoi(ps.ByName("idUser"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 403 - you cannot upload a photo for another user, photo not uploaded
	if token != img.IdOwner {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(img.IdOwner)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	var count int
	err, count = rt.db.FindUserById(img.IdOwner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 404 - user not found, photo not uploaded
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(img.IdOwner)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	err = json.NewDecoder(r.Body).Decode(&img)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	layout := "2006-01-02T15:04:05Z07:00"
	img.DateTime = time.Now().Format(layout)

	err, img.IdImage = rt.db.InsertPhoto(img.IdOwner, img.DateTime, img.Url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 201 - photo succesfully uploaded
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(img)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

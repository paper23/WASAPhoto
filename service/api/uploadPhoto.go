package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"fmt"

	"github.com/julienschmidt/httprouter"
)

// upload a new photo using data provided in the body of the request and the id of the owner of the photo, return the full image object with the ID, URL and date/time
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var img Image
	var id int
	id, err := strconv.Atoi(ps.ByName("idUser"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("1")
		return
	}

	img.IdOwner = id
	err = json.NewDecoder(r.Body).Decode(&img)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("2", err)
		return
	}

	layout := "2006-01-02T15:04:05Z07:00"
	img.DateTime = time.Now().Format(layout)
	
	err, img.IdImage = rt.db.InsertPhoto(img.IdOwner, img.DateTime, img.Url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("3")
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(img)


}
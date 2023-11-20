package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// return the user profile given its id, the user profile is composed by the user’s photos (in reverse chronological order), how many photos have been uploaded, and the user’s followers and following
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	idUser, err := strconv.Atoi(ps.ByName("idUser"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if idUser < 0 || idUser >= len(Users) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	newUsername := r.URL.Query().Get("newUsername")
	if len(newUsername) < 3 || len(newUsername) > 16 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	Users[idUser].Username = newUsername

	json.NewEncoder(w).Encode(Users[idUser])
}

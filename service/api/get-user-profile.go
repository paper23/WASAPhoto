package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type User struct {
	idUser    int
	Username  string
	Biography string
	Follow    []Follow
}

var Users = []User{
	User{
		idUser:    0,
		Username:  "lorenzo23",
		Biography: "Hi, I'm 21, I live in Rome",
		Follow: []Follow{
			{IdUserFollowed: 1},
			{IdUserFollowed: 2},
		},
	},
}

// return the user profile given its id, the user profile is composed by the user’s photos (in reverse chronological order), how many photos have been uploaded, and the user’s followers and following
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	json.NewEncoder(w).Encode(Users[idUser])
}

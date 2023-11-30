package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes

	//login
	rt.router.POST("/session", rt.login)

	//user
	rt.router.PUT("/users/:idUser", rt.setMyUserName)

	//image
	rt.router.POST("/users/:idUser/images/", rt.uploadPhoto)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}

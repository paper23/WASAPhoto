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
	rt.router.POST("/users/:idUser/follows/", rt.followUser)
	rt.router.POST("/users/:idUser/bans/", rt.banUser)
	rt.router.DELETE("/users/:idUser/bans/:idUserBanned", rt.unbanUser)

	//image
	rt.router.POST("/users/:idUser/images/", rt.uploadPhoto)
	rt.router.DELETE("/users/:idUser/images/:idImage", rt.deletePhoto)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}

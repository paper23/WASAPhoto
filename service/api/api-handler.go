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
	rt.router.DELETE("/users/:idUser/follows/:idUserToUnfollow", rt.unfollowUser)
	rt.router.POST("/users/:idUser/bans/", rt.banUser)
	rt.router.DELETE("/users/:idUser/bans/:idUserBanned", rt.unbanUser)

	//image
	rt.router.POST("/users/:idUser/images/", rt.uploadPhoto)
	rt.router.DELETE("/users/:idUser/images/:idImage", rt.deletePhoto)
	rt.router.POST("/users/:idUser/images/:idImage/comments/", rt.commentPhoto)
	rt.router.DELETE("/users/:idUser/images/:idImage/comments/:idComment", rt.uncommentPhoto)
	rt.router.POST("/users/:idUser/images/:idImage/likes/", rt.likePhoto)
	rt.router.DELETE("/users/:idUser/images/:idImage/likes/:idLiker", rt.unlikePhoto)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}

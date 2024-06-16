package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes

	// login
	rt.router.POST("/session", rt.doLogin) // FATTO

	// user
	rt.router.PUT("/users/:idUser", rt.setMyUserName)                   // USATO, DA CONTROLLARE CON API
	rt.router.POST("/users/:idUser/follows/", rt.followUser)            // FATTO
	rt.router.DELETE("/users/:idUser/follows/", rt.unfollowUser)        // FATTO
	rt.router.POST("/users/:idUser/bans/", rt.banUser)                  // FATTO
	rt.router.DELETE("/users/:idUser/bans/:idUserBanned", rt.unbanUser) // FATTO
	rt.router.GET("/users/:idUser/bans/", rt.getBannedList)             // FATTO
	rt.router.GET("/users/:idUser", rt.getUserProfile)                  // USATO, DA CONTROLLARE CON API
	rt.router.GET("/search/:username", rt.searchUser)

	// image
	rt.router.POST("/users/:idUser/images/", rt.uploadPhoto)                                  // USATO, DA CONTROLLARE CON API
	rt.router.DELETE("/users/:idUser/images/:idImage", rt.deletePhoto)                        // FATTO
	rt.router.POST("/users/:idUser/images/:idImage/comments/", rt.commentPhoto)               // FATTO
	rt.router.DELETE("/users/:idUser/images/:idImage/comments/:idComment", rt.uncommentPhoto) // FATTO
	rt.router.POST("/users/:idUser/images/:idImage/likes/", rt.likePhoto)                     // FATTO
	rt.router.DELETE("/users/:idUser/images/:idImage/likes/:idLiker", rt.unlikePhoto)         // FATTO

	// stream
	rt.router.GET("/users/:idUser/stream", rt.getMyStream) // USATO, DA CONTROLLARE CON API

	// Special routes
	rt.router.GET("/liveness", rt.liveness) // OK

	return rt.router
}

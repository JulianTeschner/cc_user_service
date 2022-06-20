package router

import (
	"log"

	"github.com/JulianTeschner/cc_user_service/user"
	"github.com/gin-gonic/gin"

	_ "github.com/gwatts/gin-adapter"
)

// New returns a new router
func New() *gin.Engine {
	log.Println("Setting up router")
	gin.ForceConsoleColor()

	r := gin.Default()

	// Wrap the http handler with gin adapter
	userGroup := r.Group("/user")
	userGroup.GET("/:name", user.GetUserHandler)
	userGroup.POST("", user.PostUserHandler)
	userGroup.DELETE("/:name", user.DeleteUserHandler)
	return r
}

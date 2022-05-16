package user

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// GetUser is the handler for the GET api/user/* route
func GetUserHandler(c *gin.Context) {
	name := c.Param("name")
	log.Println("GetUser: ", name)

	var err error
	user, err := GetUser(os.Getenv("MONGO_INITDB_ROOT_DATABASE"), "users", "username", name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &user)
}

// PostUser is the handler for the POST api/user/* route
func PostUserHandler(c *gin.Context) {
	var user User
	c.Request.ParseForm()
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	_, err = PostUser(os.Getenv("MONGO_INITDB_ROOT_DATABASE"), "users", &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &user)
}

func DeleteUserHandler(c *gin.Context) {
	name := c.Param("name")
	log.Println("DeleteUser: ", name)

	result, err := DeleteUser(os.Getenv("MONGO_INITDB_ROOT_DATABASE"), "users", "username", name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

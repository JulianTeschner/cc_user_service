package user

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

var r *gin.Engine

func setupHandlersTest() func() {
	r = gin.Default()

	userGroup := r.Group("/user")
	{
		userGroup.GET("/:name", GetUserHandler)
		userGroup.POST("", PostUserHandler)
		userGroup.DELETE("/:name", DeleteUserHandler)
	}

	return func() {
		log.Println("teardown suite")
	}
}

func TestGetUserHandler(t *testing.T) {
	teardown := addDummyUserToDb()
	defer teardown()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/fish", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserHandlerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/NotFound", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostUserHandler(t *testing.T) {
	w := httptest.NewRecorder()
	dummyUser := createDummyUser()
	data, err := json.Marshal(&dummyUser)
	if err != nil {
		log.Fatal(err)
	}
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	Client.Database("test").Collection("users").DeleteOne(context.Background(), bson.M{"name": dummyUser.Username})
}

func TestPostUserHandlerMarshallError(t *testing.T) {
	w := httptest.NewRecorder()
	data, err := json.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostUserHandlerNoConnection(t *testing.T) {
	w := httptest.NewRecorder()
	dummyUser := createDummyUser()
	data, err := json.Marshal(&dummyUser)
	if err != nil {
		log.Fatal(err)
	}
	Client.Disconnect(ctx)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	r.ServeHTTP(w, req)
	NewClient()
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteUserHandler(t *testing.T) {
	addDummyUserToDb()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/user/fish", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteUserHandlerNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	Client.Disconnect(context.Background())
	req, _ := http.NewRequest("DELETE", "/user/NotFound", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	NewClient()
}

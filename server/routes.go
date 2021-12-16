package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// API: curl -X GET http://localhost:1337/
func HelloWorld(c *gin.Context) {

	hello := Message{200, "Hello", "Hello World!"}

	c.JSON(http.StatusOK, hello)
}

func GetPeople(c *gin.Context) {

}

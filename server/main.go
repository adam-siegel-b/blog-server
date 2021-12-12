package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Code  int
	Title string
	Value string
}

func main() {
	router := gin.Default()
	router.GET("/", HelloWorld)
	router.Run(":1337")
}

// API: curl -X GET http://localhost:1337/
func HelloWorld(c *gin.Context) {

	hello := Message{200, "Hello", "Hello World!"}

	c.JSON(http.StatusOK, hello)
}

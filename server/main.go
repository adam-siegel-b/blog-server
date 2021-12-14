package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", HelloWorld)

	router.GET("/posts", GetPosts)
	router.POST("/posts", NewPost)
	router.PUT("/posts", UpdatePost)
	router.DELETE("/posts", DeletePost)
	router.Run(":1337")
}

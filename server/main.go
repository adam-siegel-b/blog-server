package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/hello", HelloWorld)
	router.Static("/static", "./../static")

	router.Run(":1337")
}

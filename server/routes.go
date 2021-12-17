package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// API: curl -X GET http://localhost:1337/
func HelloWorld(c *gin.Context) {

	hello := Message{http.StatusOK, "Hello", "Hello World!"}

	c.JSON(http.StatusOK, hello)
}

func pwd(c *gin.Context) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	out := []string{}
	for _, f := range files {
		out = append(out, f.Name())
	}
	msg := Message{http.StatusOK, "PWD", "[" + strings.Join(out, ",") + "]"}
	c.JSON(http.StatusOK, msg)
}

func env(c *gin.Context) {
	sp := os.Getenv("STATIC_ROOT")
	msg := Message{http.StatusOK, "STATIC_ROOT", sp}
	c.JSON(http.StatusOK, msg)
}

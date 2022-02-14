package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func LogIn(c *gin.Context) {
	var s Slalomer
	var s1 Slalomer
	c.ShouldBindJSON(&s)

	conn, err := connect2DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		msg := Message{http.StatusInternalServerError, "unable to connect to db", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	defer conn.Close(context.Background())

	err = conn.QueryRow(context.Background(), LOGIN_CHECK, s.Email, s.Name).Scan(&s1.ID, &s1.Name, &s1.Email, &s1.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "login unable to query: %v\n", err)
		msg := Message{http.StatusNotFound, "query error: %v", err.Error()}
		c.JSON(http.StatusNotFound, msg)
		return
	}

	if CheckPasswordHash(s.Password, s1.Password) {
		Authenticate(c, s1.ID)
		c.JSON(http.StatusOK, s1)
		return
	}
	msg := Message{http.StatusNotFound, "Not Authorized", "Not Authorized"}
	c.JSON(http.StatusNotFound, msg)
}

func LogOut(c *gin.Context) {
	DeAuthenticate(c)
	c.JSON(http.StatusOK, Message{http.StatusOK, "logout", "bye"})
}

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
	c.BindJSON(&s)

	conn, err := connect2DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		msg := Message{http.StatusInternalServerError, "unable to connect to db", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	defer conn.Close(context.Background())

	s.Password, err = HashPassword(s.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to hash: %v\n", err)
		msg := Message{http.StatusForbidden, "unable to hash password", err.Error()}
		c.JSON(http.StatusForbidden, msg)
		return
	}

	err = conn.QueryRow(context.Background(), LOGIN_CHECK, s.Email, s.Password).Scan(&s1.ID, &s1.Name, &s1.Email, &s1.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "login unable to query: %v\n", err)
		msg := Message{http.StatusNotFound, "query error: %v", err.Error()}
		c.JSON(http.StatusNotFound, msg)
		return
	}

	Authenticate(c, s1.ID)
	c.JSON(http.StatusOK, s1)
}

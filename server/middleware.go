package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

var activeSession = map[string]string{}

func AuthReq(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("token")
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if val, ok := activeSession[fmt.Sprintf("%v", user)]; ok {
		if val != hashActiveUser(session.Get("user")) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()

	}
	DeAuth(c)
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
}

func DeAuth(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	// session.Options(sessions.Options{Path: "/", MaxAge: -1})
	session.Save()
}

func Authenticate(c *gin.Context, id string) {
	sesh := sessions.Default(c)
	userval := hashActiveUser(id)
	sesh.Set("user", userval)
	sid := uuid.New()
	sesh.Set("token", sid.String())
	activeSession[sid.String()] = id
	if err := sesh.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
}

func DeAuthenticate(c *gin.Context) {
	sesh := sessions.Default(c)
	t := sesh.Get("token")
	if t != nil {
		delete(activeSession, t.(string))
	}
	sesh.Clear()
	if err := sesh.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
}

func hashActiveUser(val interface{}) string {
	h := sha3.New512()
	h.Write([]byte(fmt.Sprintf("%v", val)))
	sum := h.Sum([]byte("H@XxY"))
	return string(sum)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	fmt.Fprintf(os.Stderr, "pass: %s hash: %s\n", password, string(bytes))
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

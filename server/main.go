package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func main() {
	dbURI := fmt.Sprintf("postgres://%s:%s@database:5432/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	conn, err := pgx.Connect(context.Background(), dbURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	err = conn.QueryRow(context.Background(), "select current_database();").Scan(&name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	router := gin.Default()
	store := sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 60 * 60 * 24})
	router.Use(sessions.Sessions("mysession", store))
	// testing functions
	router.GET("/hello", HelloWorld)
	router.GET("/pwd", pwd)
	router.GET("/env", env)

	// Auth
	router.POST("/user", CreateUser)
	router.PUT("/user", UpdateUser)
	router.GET("/user", GetUser)

	router.GET("/users", GetAllUsers)

	router.PUT("/login", LogIn)

	sp := os.Getenv("STATIC_ROOT")
	if len(sp) == 0 {
		sp = "./../static"
	}
	router.Static("/static", sp)

	router.Run(":1337")
}

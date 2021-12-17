package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func CreateUser(c *gin.Context) {
	// get the values and turn them into a slalomer object
	var s Slalomer
	c.BindJSON(&s)

	// create a UUIDv4
	uuuid := uuid.New()
	s.ID = uuuid.String()

	// Validate email
	if !validEmail(s.Email) {
		fmt.Fprintf(os.Stderr, "bad email: %v\n", s.Email)
		msg := Message{http.StatusBadRequest, "invalid email", s.Email}
		c.JSON(http.StatusForbidden, msg)
		return
	}
	// validate the user is from slalom
	// if !validSlalomEmail(s.Email) {
	//  fmt.Fprintf(os.Stderr, "not a slalom email: %v\n", s.Email)
	// 	msg := Message{http.StatusForbidden, "not a slalom email", s.Email}
	// 	c.JSON(http.StatusForbidden, msg)
	// 	return
	// }

	// remove SQL injection from Name
	s.Name = stripSketchyChars(s.Name)

	// Hash password (security first)
	soHash, err := HashPassword(s.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to hash: %v\n", err)
		msg := Message{http.StatusForbidden, "unable to hash password", err.Error()}
		c.JSON(http.StatusForbidden, msg)
		return
	}
	s.Password = soHash

	// Db Connection
	dbURI := fmt.Sprintf("postgres://%s:%s@database:5432/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	conn, err := pgx.Connect(context.Background(), dbURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		msg := Message{http.StatusInternalServerError, "unable to connect to db", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "INSERT INTO slalomer (id,name,email, password) VALUES ($1,$2,$3,$4)", s.ID, s.Name, s.Email, s.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Insert record failed: %v\n", err)
		msg := Message{http.StatusInternalServerError, "Insert failed", err.Error()}
		DeAuth(c)
		c.JSON(http.StatusInternalServerError, msg)
		return
	}

	//createSession
	Authenticate(c, s.ID)

	c.JSON(http.StatusOK, fmt.Sprintf("%+v", s))
}

func UpdateUser(c *gin.Context) {
	var s Slalomer
	var s1 Slalomer
	var l1 Location
	// Db Connection
	dbURI := fmt.Sprintf("postgres://%s:%s@database:5432/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	conn, err := pgx.Connect(context.Background(), dbURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		msg := Message{http.StatusInternalServerError, "unable to connect to db", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	defer conn.Close(context.Background())

	// get ID from sessionhash
	sesh := sessions.Default(c)
	seshKey := fmt.Sprintf("%v", sesh.Get("token"))
	id := activeSession[seshKey]
	err = conn.QueryRow(context.Background(), "Select id, name, email, password FROM slalomer WHERE id=$1", id).Scan(&s.ID, &s.Name, &s.Email, &s.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to query: %v\n", err)
		msg := Message{http.StatusNotFound, "query error: %v", err.Error()}
		c.JSON(http.StatusNotFound, msg)
		return
	}
	// get the values and turn them into a slalomer object
	c.BindJSON(&s1)
	c.BindJSON(&l1)

	// Validate email
	if !validEmail(s1.Email) {
		fmt.Fprintf(os.Stderr, "bad email: %v\n", s1.Email)
		msg := Message{http.StatusBadRequest, "invalid email", s1.Email}
		c.JSON(http.StatusForbidden, msg)
		return
	}
	s.Email = s1.Email
	// validate the user is from slalom
	// if !validSlalomEmail(s.Email) {
	//  fmt.Fprintf(os.Stderr, "not a slalom email: %v\n", s.Email)
	// 	msg := Message{http.StatusForbidden, "not a slalom email", s.Email}
	// 	c.JSON(http.StatusForbidden, msg)
	// 	return
	// }

	// remove SQL injection from Name
	s.Name = stripSketchyChars(s1.Name)

	// Hash password (security first)
	soHash, err := HashPassword(s1.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to hash: %v\n", err)
		msg := Message{http.StatusForbidden, "unable to hash password", err.Error()}
		c.JSON(http.StatusForbidden, msg)
		return
	}
	s.Password = soHash

	_, err = conn.Exec(context.Background(), "INSERT INTO slalomer (id,name,email, password) VALUES ($1,$2,$3,$4)", s.ID, s.Name, s.Email, s.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Insert record failed: %v\n", err)
		msg := Message{http.StatusInternalServerError, "Insert failed", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("%+v", s))
}

func GetUser(c *gin.Context) {
	var s Slalomer

	// Db Connection
	dbURI := fmt.Sprintf("postgres://%s:%s@database:5432/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	conn, err := pgx.Connect(context.Background(), dbURI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		msg := Message{http.StatusInternalServerError, "unable to connect to db", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	defer conn.Close(context.Background())

	// get ID from sessionhash
	sesh := sessions.Default(c)
	seshKey := fmt.Sprintf("%v", sesh.Get("token"))
	id := activeSession[seshKey]
	fmt.Fprintf(os.Stderr, "idval: %v\n", id)
	err = conn.QueryRow(context.Background(), "Select id, name, email, password FROM slalomer WHERE id=$1", id).Scan(&s.ID, &s.Name, &s.Email, &s.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to query: %v\n", err)
		msg := Message{http.StatusNotFound, "query error: %v", err.Error()}
		c.JSON(http.StatusNotFound, msg)
		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("%+v", s))
}

func GetAllUsers(c *gin.Context) {
	// rows, err := Conn.Query(context.Background(), "SELECT * FROM languages")
	// if err != nil {
	//     log.Fatal(err)
	// }
	// defer rows.Close()

	// var rowSlice []Row
	// for rows.Next() {
	//     var r Row
	//     err := rows.Scan(&r.Id, &r.Language, &r.Name)
	//     if err != nil {
	//         log.Fatal(err)
	//     }
	//    rowSlice = append(rowSlice, r)
	// }
	// if err := rows.Err(); err != nil {
	//     log.Fatal(err)
	// }

	// fmt.Println(rowSlice)
}

func LogIn(c *gin.Context) {
	var s Slalomer
	var s1 Slalomer
	c.BindJSON(&s)

	dbURI := fmt.Sprintf("postgres://%s:%s@database:5432/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	conn, err := pgx.Connect(context.Background(), dbURI)
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

	err = conn.QueryRow(context.Background(), "Select id, name, email, password FROM slalomer WHERE email=$1 and password=$2", s.Email, s.Password).Scan(&s1.ID, &s1.Name, &s1.Email, &s1.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to query: %v\n", err)
		msg := Message{http.StatusNotFound, "query error: %v", err.Error()}
		c.JSON(http.StatusNotFound, msg)
		return
	}

	Authenticate(c, s1.ID)
	c.JSON(http.StatusOK, fmt.Sprintf("%+v", s1))
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/pgtype"
)

func CreateUser(c *gin.Context) {
	// get the values and turn them into a slalomer object
	var s Slalomer
	c.BindJSON(&s)

	// create a UUIDv4
	uuuid := uuid.New()
	s.ID = uuuid.String()

	//validate the user is from slalom
	if !validSlalomEmail(s.Email) {
		fmt.Fprintf(os.Stderr, "not a slalom email: %v\n", s.Email)
		msg := Message{http.StatusForbidden, "not a slalom email", s.Email}
		c.JSON(http.StatusForbidden, msg)
		return
	}

	// remove SQL injection from Name
	s.Name = stripSketchyChars(s.Name)

	fmt.Fprintf(os.Stderr, "s: %v\n", s)
	// Hash password (security first)
	soHash, err := HashPassword(s.Password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to hash: %v\n", err)
		msg := Message{http.StatusForbidden, "unable to hash password", err.Error()}
		c.JSON(http.StatusForbidden, msg)
		return
	}
	s.Password = soHash
	fmt.Fprintf(os.Stderr, "s: %v\n", s)

	// Db Connection
	conn, err := connect2DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		msg := Message{http.StatusInternalServerError, "unable to connect to db", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), CREATE_SLALOMER, s.ID, s.Name, s.Email, s.Password, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Insert record failed: %v\n", err)
		msg := Message{http.StatusInternalServerError, "Insert failed", err.Error()}
		DeAuth(c)
		c.JSON(http.StatusInternalServerError, msg)
		return
	}

	//createSession
	Authenticate(c, s.ID)

	c.JSON(http.StatusOK, s)
}

func UpdateUser(c *gin.Context) {
	var s Slalomer
	var s1 Slalomer
	var loc sql.NullString

	// Db Connection
	conn, err := connect2DB()
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
	err = conn.QueryRow(context.Background(), READ_SLALOMER, id).Scan(&s.ID, &s.Name, &s.Email, &s.Password, &loc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "update unable to query: %v\n", err)
		msg := Message{http.StatusNotFound, "query error", err.Error()}
		c.JSON(http.StatusNotFound, msg)
		return
	}

	if loc.Valid {
		s.Location.ID = loc.String
		err = getLocation(c, &s.Location, conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "location error: %v\n", err)
			msg := Message{http.StatusInternalServerError, "location error", err.Error()}
			c.JSON(http.StatusInternalServerError, msg)
			return
		}
	}
	// get the values and turn them into a slalomer object
	err = c.BindJSON(&s1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parseError: %+v\n", err)
		return
	}
	if IsValidUUID(s1.ID) {
		s.ID = s1.ID
	}
	if len(s1.Name) > 0 {
		// remove SQL injection from Name
		s.Name = stripSketchyChars(s1.Name)
	}

	fmt.Fprintf(os.Stderr, "location: %+v\n", s.Location)
	if len(s1.Location.ID) > 0 {
		s.Location.ID = s1.Location.ID
	}
	if len(s1.Location.Name) > 0 {
		s.Location.Name = stripSketchyChars(s1.Location.Name)
	}
	if s1.Location.Lat != 0 {
		s.Location.Lat = s1.Location.Lat
	}
	if s1.Location.Lon != 0 {
		s.Location.Lon = s1.Location.Lon
	}

	// validate the user is from slalom
	if len(s1.Email) < 1 || !validSlalomEmail(s1.Email) {
		fmt.Fprintf(os.Stderr, "not a slalom email: %v\n", s.Email)
		msg := Message{http.StatusForbidden, "not a slalom email", s.Email}
		c.JSON(http.StatusForbidden, msg)
		return
	}

	s.Email = s1.Email

	// Hash password (security first)
	if len(s1.Password) > 0 {
		soHash, err := HashPassword(s1.Password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to hash: %v\n", err)
			msg := Message{http.StatusForbidden, "unable to hash password", err.Error()}
			c.JSON(http.StatusForbidden, msg)
			return
		}
		s.Password = soHash
	}

	fmt.Fprintf(os.Stderr, "location: %+v\n", s.Location)
	err = upsertLocation(c, &s.Location, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Update record failed: %v\n", err)
		msg := Message{http.StatusInternalServerError, "Insert location failed", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}

	_, err = conn.Exec(context.Background(), UPDATE_SLALOMER, s.ID, s.Name, s.Email, s.Password, s.Location.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Update record failed: %v\n", err)
		msg := Message{http.StatusInternalServerError, "Update failed", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	fmt.Fprintf(os.Stderr, "C: %+v\n", c)
	c.JSON(http.StatusOK, s)
}

func GetUser(c *gin.Context) {
	var s Slalomer
	var loc sql.NullString

	// Db Connection
	conn, err := connect2DB()
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
	err = conn.QueryRow(context.Background(), READ_SLALOMER, id).Scan(&s.ID, &s.Name, &s.Email, &s.Password, &loc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get unable to query: %v\n", err)
		msg := Message{http.StatusNotFound, "query error", err.Error()}
		c.JSON(http.StatusNotFound, msg)
		return
	}
	if loc.Valid {
		s.Location.ID = loc.String
		fmt.Fprintf(os.Stderr, "loc: %+v\n", s.Location)
		err = getLocation(c, &s.Location, conn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "location error: %v\n", err)
			msg := Message{http.StatusInternalServerError, "location error", err.Error()}
			c.JSON(http.StatusInternalServerError, msg)
			return
		}
	}

	c.JSON(http.StatusOK, s)
}

func DeleteUser(c *gin.Context) {
	// get the values and turn them into a slalomer object
	var s Slalomer
	c.BindJSON(&s)

	conn, err := connect2DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		msg := Message{http.StatusInternalServerError, "unable to connect to db", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), DELETE_SLALOMER, s.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Delete record failed: %v\n", err)
		msg := Message{http.StatusInternalServerError, "Delete failed", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}

	c.JSON(http.StatusOK, s)
}

func GetAllUsers(c *gin.Context) {
	conn, err := connect2DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		msg := Message{http.StatusInternalServerError, "unable to connect to db", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), GET_ALL_SLALOMERS, 100, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get rows: %v\n", err)
		msg := Message{http.StatusInternalServerError, "unable to get rows", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}
	defer rows.Close()

	var Alls Slalomers
	for rows.Next() {
		var r Slalomer
		var l Location
		var loc sql.NullString
		var pho sql.NullString
		var lID sql.NullString
		var name sql.NullString
		var latlon pgtype.Point
		err := rows.Scan(&r.ID, &r.Name, &r.Email, &r.Password, &pho, &loc, &lID, &name, &latlon)
		if err != nil {
			fmt.Fprintf(os.Stderr, "row scan error: %v\n", err)
			msg := Message{http.StatusInternalServerError, "row scan error", err.Error()}
			c.JSON(http.StatusInternalServerError, msg)
			return
		}
		if loc.Valid {
			l.ID = loc.String
			l.Name = name.String
			l.Lat = latlon.P.Y
			l.Lon = latlon.P.X
			fmt.Fprintf(os.Stderr, "lID name latlon: %+v %+v %+v\n", lID.String, name.String, latlon)
		}
		r.Password = ""
		r.Location = l
		Alls.Users = append(Alls.Users, r)
	}
	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "rows err: %v\n", err)
		msg := Message{http.StatusInternalServerError, "rows err", err.Error()}
		c.JSON(http.StatusInternalServerError, msg)
		return
	}

	c.JSON(http.StatusOK, Alls)
}

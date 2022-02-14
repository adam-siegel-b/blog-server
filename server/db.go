package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func connect2DB() (*pgx.Conn, error) {
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_SERVER_NAME"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))
	return pgx.Connect(context.Background(), dbURI)
}

var CREATE_SLALOMER = "INSERT INTO slalomer (id,name,email, password,location) VALUES ($1,$2,$3,$4,$5)"
var READ_SLALOMER = "SELECT id, name, email, password, location FROM slalomer WHERE id=$1"
var UPDATE_SLALOMER = "UPDATE slalomer SET name = $2, email = $3, password = $4, location = $5 WHERE id = $1 RETURNING id, name, email, location"
var DELETE_SLALOMER = "DELETE FROM slalomer WHERE id=$1"

var GET_ALL_SLALOMERS = "SELECT * FROM slalomer s INNER JOIN location l ON (l.id = s.location) LIMIT $1 OFFSET $2"
var LOGIN_CHECK = "SELECT id, name, email, password FROM slalomer WHERE email=$1 or name=$2"

var CREATE_LOCATION = "INSERT INTO location (id,name,latlon) VALUES ($1,$2,$3)"
var READ_LOCATION = "SELECT id, name, latlon FROM Location WHERE id=$1"
var UPDATE_LOCATION = "UPDATE location SET name = $2, latlon= $3 WHERE id=$1"
var DELETE_LOCATION = "DELETE FROM location WHERE id=$1"

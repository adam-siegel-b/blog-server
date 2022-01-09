package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v4"
)

func getLocation(c *gin.Context, l *Location, cn *pgx.Conn) error {
	var latlon pgtype.Point
	var n sql.NullString
	err := cn.QueryRow(context.Background(), READ_LOCATION, l.ID).Scan(&l.ID, &n, &latlon)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to query: %v\n", err)
		return err
	}

	if n.Valid {
		l.Name = n.String
	}

	l.Lat = latlon.P.X
	l.Lon = latlon.P.Y
	return nil
}

func upsertLocation(c *gin.Context, l *Location, cn *pgx.Conn) error {
	var action string
	if len(l.ID) > 0 {
		err := getLocation(c, l, cn)
		if err != nil {
			return err
		}
		action = UPDATE_LOCATION
	} else {
		uuuid := uuid.New()
		l.ID = uuuid.String()
		action = CREATE_LOCATION
	}

	var latlon pgtype.Point
	latlon.Status = pgtype.Present
	latlon.P.X = l.Lat
	latlon.P.Y = l.Lon
	qll, err := latlon.Value()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse Lat Lon: %v\n", err)
		return err
	}

	_, err = cn.Exec(context.Background(), action, l.ID, l.Name, qll)
	if err != nil {
		fmt.Fprintf(os.Stderr, "latlon: %+v\n l: %+v\n", latlon, l)
		fmt.Fprintf(os.Stderr, "Insert Location record failed: %v\n", err)
		return err
	}

	return nil
}

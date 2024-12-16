package main

import (
	"database/sql"
	"dtla/internal/util"
	"errors"
	"net/http"
	"strconv"
	"time"
)

// Fill out Auth.ID from request cookie and Auth.{User,SID} from database
func getUserData(r *http.Request, db *sql.DB, auth *util.Auth) error {
	var err error

	idCookie, err := r.Cookie("id")
	if err != nil {
		return err
	}

	auth.ID, err = strconv.Atoi(idCookie.Value)
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT user, sid, sstart, sage FROM users WHERE id IS ?", auth.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return errors.New("No rows")
	}

	var sstart int64
	var sage int
	err = rows.Scan(&auth.User, &auth.SID, &sstart, &sage)
	if err != nil {
		return err
	}

	auth.SStart = time.Unix(sstart, 0)
	SEnd := auth.SStart.Add(time.Second * time.Duration(sage))
	if time.Now().After(SEnd) {
		return util.ErrSessionExpired
	}

	return nil
}

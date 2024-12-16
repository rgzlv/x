package util

import (
	"database/sql"
	"net/http"
	"time"
)

type Auth struct {
	// From request cookie
	ID int

	// From database
	User string

	// From database
	SID []byte

	// Session creation time
	// Stored as seconds since UNIX epoch in database
	SStart time.Time

	// Session maximum age until it expires
	// Stored as seconds in database
	SAge time.Duration

	// ASDefault - the sid cookie wasn't found
	// ASError - a different error occured
	// ASOk when the sid cookie hash comparison succeeds with the one in the database
	Status uint

	// Repeated here so the identifiers can be used in templates
	// Initialized in makeHandler()
	ASDefault int
	ASError   int
	ASOk      int

	// Should be set when Status == ASError
	Error string
}

const (
	ASDefault = iota
	ASError
	ASOk
)

func ValidCookie(r *http.Request, db *sql.DB) (*Auth, bool, error) {
	var err error

	cookie, err := r.Cookie("uid")
	if err != nil {
		return nil, false, nil
	}

	_ = cookie
	return nil, false, nil
}

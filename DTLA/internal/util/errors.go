package util

import "errors"

// Returned if session start (seconds since UNIX epoch) + session age (seconds) > time.Now()
var ErrSessionExpired error = errors.New("Sessija ir beigusies")

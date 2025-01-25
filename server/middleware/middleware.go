package middleware

import "webserver/database"

type Middleware struct {
	Database database.DatabaseAPI
}

// NewMiddleware is a constructor for the Middleware struct
func NewMiddleware(db database.DatabaseAPI) *Middleware {
	return &Middleware{
		Database: db,
	}
}

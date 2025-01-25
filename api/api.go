package api

import (
	"webserver/database"
)

type API struct {
	Database database.DatabaseAPI
}

func NewAPI(db database.DatabaseAPI) *API {
	return &API{
		Database: db,
	}
}

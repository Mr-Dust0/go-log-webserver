package main

import (
	"webserver/api"
	"webserver/database"
	"webserver/server"
	"webserver/utils"
)

func main() {
	db := database.NewDatabaseFromEnv()
	if utils.IsTestEnv() {
		db.InsertTestData()
	}

	api := api.NewAPI(db)
	
	server := server.NewServer(db, api)
	server.Start()
}

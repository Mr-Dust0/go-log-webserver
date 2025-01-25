package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"webserver/api"
	"webserver/database"
	"webserver/server/middleware"
	"webserver/utils"
)

type Server struct {
	Engine     *gin.Engine
	API        *api.API
	Middleware *middleware.Middleware
	Database   database.DatabaseAPI
}

func NewServer(db database.DatabaseAPI, api *api.API) *Server {
	return &Server{
		Middleware: middleware.NewMiddleware(db),
		API:        api,
	}
}

// Start initiates the HTTP server
func (server *Server) Start() {
	// Declare the server
	server.Engine = gin.Default()

	server.AttachRoutes()

	if utils.IsTestEnv() {
		err := server.Engine.Run("127.0.0.1:8080")
		if err != nil {
			fmt.Printf("Server error: %v", err)
		}
		return
	}

	env := utils.GetDotEnvVariables()
	err := server.Engine.RunTLS(":443", env["CERT"], env["CERT_KEY"])
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	_ "encoding/json"
	"webserver/controllers"
	"webserver/initializers"
	"webserver/middleware"

	"github.com/gin-gonic/gin"
)

func initialize() {
	initializers.LoadEnvs()
	initializers.InitDatabase()
	// Create sample data this should only be used when testing the application
	initializers.InsertTestData()
}

func main() {

	// Initlaize database and load .env file
	initialize()
	// Create router
	r := gin.Default()
	r.StaticFile("/favicon.ico", "./favicon.ico")
	r.Static("/static", "./static")
	// Tell the router when to find the html files from
	r.LoadHTMLGlob("templates/**/*.html")
	// Run CheckAuth before GetHomePageHandler to make sure that the user is authenicated before being able to see the logs
	// The index is used to get the header and footer
	r.GET("/", middleware.CheckAuth, middleware.GetUsedLoggedIn, controllers.GetIndex)
	// Populate log table for index page
	r.GET("/logs", middleware.CheckAuth, middleware.GetUsedLoggedIn, controllers.GetHomePageHandler)
	// Get suggestions for what hostname to search for
	r.GET("/hostname-suggestions", controllers.HomeSuggestions)
	// Populate the log table but only with data that meets search parameters
	r.POST("/logs", middleware.CheckAuth, middleware.GetUsedLoggedIn, controllers.PostSearch)
	// Used by qr-code reader to add logs to the database
	r.POST("/log", controllers.InsertLog)
	// Used by qr-code reader to update time closed
	r.PUT("/fileclosed", controllers.UpdateTimeClosed)
	r.GET("/login", middleware.GetUsedLoggedIn, controllers.GetLoginPage)
	r.POST("/login", middleware.GetUsedLoggedIn, controllers.Login)
	r.GET("/reset", middleware.GetUsedLoggedIn, controllers.GetResetPage)
	r.GET("/username", middleware.GetUser)
	r.POST("/reset", controllers.ChangePassword)
	// If running in production use this to use TLS/https instead of using http and allow any on the network to reach the application
	//err := r.RunTLS(":443", initializers.EnvFile["CERT"], initializers.EnvFile["CERT_KEY"])
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	// Run on localhost to make sure no one else on the network can hit the application
	r.Run("127.0.0.1:8080")
}

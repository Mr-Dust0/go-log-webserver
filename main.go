package main

import (
	_ "encoding/json"
	"net/http"
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
	r.GET("/", middleware.CheckAuth, middleware.GetUsedLoggedIn, controllers.GetIndex)
	r.GET("/logs", middleware.CheckAuth, middleware.GetUsedLoggedIn, controllers.GetHomePageHandler)
	r.GET("/2", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index2.html", gin.H{})
	})

	r.GET("/hostname-suggestions", controllers.HomeSuggestions)
	// r.GET("/hostname-suggestions", func(ctx *gin.Context) {
	// 	var logs []models.LogEntry
	// 	hostnames := make([]string, 0)
	// 	hostname := ctx.Query("hostname")
	// 	hostname = "%" + hostname + "%" //
	// 	initializers.DB.Where("host_name LIKE?", hostname).Find(&logs)
	// 	for _, log := range logs {
	// 		hostnames = append(hostnames, log.HostName)
	// 		fmt.Println(log.HostName)
	// 	}
	// 	ctx.HTML(http.StatusOK, "suggestions.html", gin.H{"hostnames": hostnames})
	// })
	r.POST("/", middleware.CheckAuth, middleware.GetUsedLoggedIn, controllers.PostHomePageHandler)
	r.POST("/logtable", middleware.CheckAuth, middleware.GetUsedLoggedIn, controllers.PostSearch)
	r.POST("/log", controllers.InsertLog)
	r.PUT("/fileclosed", controllers.UpdateTimeClosed)
	r.GET("/login", middleware.GetUsedLoggedIn, controllers.GetLoginPage)
	r.POST("/login", middleware.GetUsedLoggedIn, controllers.Login)
	r.GET("/reset", middleware.GetUsedLoggedIn, controllers.GetResetPage)
	r.POST("/reset", middleware.GetUsedLoggedIn, controllers.ChangePassword)
	// If running in production use this to use TLS/https instead of using http and allow any on the network to reach the application
	//err := r.RunTLS(":443", initializers.EnvFile["CERT"], initializers.EnvFile["CERT_KEY"])
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	// Run on localhost to make sure no one else on the network can hit the application
	r.Run("127.0.0.1:8080")
}

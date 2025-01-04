package main

import (
	_ "encoding/json"
	"webserver/controllers"
	"webserver/initializers"
	middlewares "webserver/middleware"

	"github.com/gin-gonic/gin"
)

func initialize() {
	initializers.LoadEnvs()
	initializers.InitDatabase()
	// Create sample data this should only be used when testing the application
	initializers.InsertTestData()
}

func main() {

	initialize()
	r := gin.Default()
	r.StaticFile("/favicon.ico", "./favicon.ico")
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/**/*.html")
	//r.LoadHTMLGlob("templates/partials/*.html")
	r.GET("/", middlewares.CheckAuth, controllers.GetHomePageHandler)
	r.POST("/", middlewares.CheckAuth, controllers.PostHomePageHandler)
	r.POST("/log", controllers.InsertLog)
	r.PUT("/fileclosed", controllers.UpdateTimeClosed)
	r.GET("/login", controllers.GetLoginPage)
	r.POST("/login", controllers.Login)
	r.GET("/reset", controllers.GetResetPage)
	r.POST("/reset", controllers.ChangePassword)
	// If running in production use this to use TLS/https instead of using http
	//err := r.RunTLS(":443", initializers.EnvFile["CERT"], initializers.EnvFile["CERT_KEY"])
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	r.Run("127.0.0.1:8080")
}

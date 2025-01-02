package main

import (
	_ "encoding/json"
	"log"
	"net/http"
	"webserver/controllers"
	"webserver/initializers"
	middlewares "webserver/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	initializers.InitDatabase()
	initializers.LoadEnvs()
	r := gin.Default()
	r.StaticFile("/favicon.ico", "./favicon.ico")
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*.html")
	r.GET("/", middlewares.CheckAuth, controllers.GetHomePageHandler)
	r.POST("/", middlewares.CheckAuth, controllers.PostHomePageHandler)
	r.POST("/log", controllers.InsertLog)
	r.PUT("/fileclosed", controllers.UpdateTimeClosed)
	r.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.POST("/login", controllers.Login)
	err := r.RunTLS(":443", initializers.EnvFile["CERT"], initializers.EnvFile["CERT_KEY"])
	if err != nil {
		log.Fatal(err)
	}
}

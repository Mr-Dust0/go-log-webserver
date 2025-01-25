package server

// GetRouter returns a list of routes available and being listened too on the server
func (server *Server) AttachRoutes() {
	server.loadHTMLFiles()

	server.attachStaticRoutes()
	server.attachUserRoutes()
	server.attachLogRoutes()
	server.attachIndexRoutes()
	server.attachSuggestionRoutes()
}

func (server *Server) attachStaticRoutes() {
	server.Engine.StaticFile("/favicon.ico", "./favicon.ico")
	server.Engine.Static("/static", "./templates/static")
}

func (server *Server) attachUserRoutes() {
	server.Engine.GET("/login", server.Middleware.GetUsedLoggedIn, server.API.GetLoginPage)
	server.Engine.POST("/login", server.Middleware.GetUsedLoggedIn, server.API.RequestLogin)
	server.Engine.GET("/reset", server.Middleware.GetUsedLoggedIn, server.API.GetChangePasswordPage)
	server.Engine.POST("/reset", server.API.RequestChangePassword)
	server.Engine.GET("/username", server.Middleware.GetUser)
}

func (server *Server) attachLogRoutes() {
	server.Engine.GET("/logs", server.Middleware.CheckAuth, server.Middleware.GetUsedLoggedIn, server.API.GetLogTablePage)
	server.Engine.POST("/logs", server.Middleware.CheckAuth, server.Middleware.GetUsedLoggedIn, server.API.SearchLogs)
	server.Engine.POST("/log", server.API.InsertLog)
	server.Engine.PUT("/fileclosed", server.API.CloseLog)
	server.Engine.GET("/openfile", server.API.GetLogTablePage)
}

func (server *Server) attachIndexRoutes() {
	server.Engine.GET("/", server.Middleware.CheckAuth, server.Middleware.GetUsedLoggedIn, server.API.GetDashboardPage)
}

func (server *Server) attachSuggestionRoutes() {
	server.Engine.GET("/hostname-suggestions", server.API.GetSuggestions)
}

func (server *Server) loadHTMLFiles() {
	server.Engine.LoadHTMLGlob("templates/**/*.html")
}

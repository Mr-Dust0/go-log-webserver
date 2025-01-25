package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) GetDashboardPage(ctx *gin.Context) {
	// Display the index page which is used by html to load the logs
	username, _ := ctx.Get("userName")

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"date":     "Showing all logs for all days",
		"userName": username.(string),
	})
}

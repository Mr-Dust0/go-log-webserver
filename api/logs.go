package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"webserver/database"
	"webserver/types"
	"webserver/utils"
)

func (api *API) GetLogTablePage(ctx *gin.Context) {
	checkbox := ctx.Query("showopenonly")

	// Find all the logs in the database
	logs, err := api.Database.ListLogs(database.ListLogsParams{
		OnlyOpen: checkbox == "on",
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	// Load only the table which is used by htmx to be displayed on the index page
	ctx.HTML(http.StatusOK, "logtable.html", gin.H{
		"Logs": logs.Format(),
		"date": "Showing all logs for all days",
	})
}

func (api *API) GetSuggestions(ctx *gin.Context) {
	hostname := ctx.Query("hostname")

	// Find all the logs in the database that have the currently typed in query anywhere in the hostname
	logs, err := api.Database.ListLogs(database.ListLogsParams{
		HostnameSimilar: hostname,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	hostnames := make([]string, 0)

	for _, log := range logs {
		// Append the hostname to an array of hostnames that match the entered term
		hostnames = append(hostnames, log.HostName)
	}

	// Get rid of duplicates so only show in once as an data entry and not many times
	hostnames = utils.RemoveDuplicatesInSlice(hostnames)

	// Load the suggestions template which creates the datalist which is used by htmx
	ctx.HTML(http.StatusOK, "suggestions.html", gin.H{"hostnames": hostnames})
}

func (api *API) SearchLogs(ctx *gin.Context) {
	// Gets data from user input
	date := ctx.PostForm("date")
	hostname := ctx.PostForm("hostname")

	if date == "" && hostname == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Date or HostName is required"})
		return
	}

	// This is to be displayed to the user depending on what data was passed into the form
	datemessage := ""
	switch {
	case date != "" && hostname != "":
		datemessage = "Showing results for the day: " + date + " and the hostname " + hostname
	case date != "":
		datemessage = "Showing results for the day: " + date
	case hostname != "":
		datemessage = "Showing results for the hostname: " + hostname
	}

	params := database.ListLogsParams{}
	if date != "" {
		params.Date = date
	}

	if hostname != "" {
		params.Hostname = hostname
	}

	// Find all the logs in the database that have the currently typed in query anywhere in the hostname
	logs, err := api.Database.ListLogs(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	// Load only the table which is used by htmx to be shown on the index page
	ctx.HTML(http.StatusOK, "logtable.html", gin.H{
		"Logs": logs.Format(),
		"date": datemessage,
	})
}


func (api *API) CloseLog(ctx *gin.Context) {
	var request struct {
		HostName string `json:"hostname"`
	}
	// Get json data from qr-code reqder and store it in request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	logs, err := api.Database.ListLogs(database.ListLogsParams{
		OnlyOpen: true,
		Date: time.Now().Format("2006-01-02"),
		Hostname: request.HostName,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return	
	}

	if len(logs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No open logs found for this hostname"})
		return
	}

	// Update TimeStampClosed for matching hostname and where TimeStampClosed is still null
	for _, log := range logs {
		err := api.Database.UpdateLogCloseTime(&log, time.Now())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update log"})
			return
		}
	}

	// Send sucess message back to qr-code reader
	ctx.JSON(http.StatusOK, gin.H{"status": "File(s) closed successfully"})
}

func (api *API) InsertLog(ctx *gin.Context) {

	var logEntry types.LogEntry
	// Get Json data from qr code reader request and store the data into logEntry
	if err := ctx.ShouldBindJSON(&logEntry); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	logEntry.TimeStamp = time.Now()

	logs, err := api.Database.ListLogs(database.ListLogsParams{
		OnlyOpen: true,
		Date: time.Now().Format("2006-01-02"),
		Hostname: logEntry.HostName,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	for _, log := range logs {
		err := api.Database.UpdateLogCloseTime(&log, time.Now())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update log"})
			return
		}
	}

	// Create an new entry in the database for the log
	if err := api.Database.InsertLog(&logEntry); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert log"})
		return
	}

	// Return sucess message to the qr-code reader
	ctx.JSON(http.StatusOK, gin.H{"status": "Log added successfully"})
}

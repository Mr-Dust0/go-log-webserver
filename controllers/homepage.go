package controllers

import (
	"fmt"
	"net/http"
	"time"
	"webserver/initializers"
	"webserver/models"

	"github.com/gin-gonic/gin"
)

func GetIndex(ctx *gin.Context) {
	// Display the index page which is used by htmx to load the logs
	userName, _ := ctx.Get("userName")
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"date":     "Showing all logs for all days",
		"userName": userName.(string),
	})
}

func GetHomePageHandler(ctx *gin.Context) {
	var logs []models.LogEntry
	var datemessage string
	var err error
	// Get all the paramaters from the get parameters
	hostname := ctx.Query("hostname")
	checkbox := ctx.Query("openonly")
	date := ctx.Query("date")

	// Start with a base query that we can append where statements to depending on what parameters are passed on
	query := initializers.DB.Model(&models.LogEntry{})
	// Add date where statement if date was passed on
	if date != "" {
		query = query.Where("DATE(time_stamp) = ?", date)

		datemessage = "Showing results for the day: " + date
	} else {
		datemessage = "Showing all logs for all days"
	}

	if hostname != "" {
		// Add date to query if hostname is not empty
		query = query.Where("host_name = ?", hostname)
		// Check if date is not empty to show the correct message to the user
		if date != "" {

			datemessage = "Showing results for the day: " + date + " and the hostname: " + hostname
		} else {
			datemessage = "Showing results for the hostname: " + hostname
		}
	}
	if checkbox == "on" {
		// Check if the date has an empty time stamp closed which means the file is still open
		query = query.Where("time_stamp_closed = ?", "0001-01-01 00:00:00+00:00")
		datemessage = "Showing open logs"
	}

	// Execute the query and fetch the logs
	err = query.Find(&logs).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	formattedLogs := formatLogs(logs)

	// Load only the table which is used by htmx to be displayed on the index page
	ctx.HTML(http.StatusOK, "logtable.html", gin.H{
		"Logs": formattedLogs,
		"date": datemessage,
	})
}

func formatDuration(d time.Duration) string {
	// Convert the duration to minutes and seconds
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60

	// Format the result as "Xm Ys"
	return fmt.Sprintf("%dm %ds", minutes, seconds)
}

func formatLogs(logs []models.LogEntry) []models.FormatedLog {
	// Create list of FormattedLogs with the capacity the same size of the logs which should reduce the amount of heap allocations that need to happen and speed up performance an tiny bit
	formatLogs := make([]models.FormatedLog, 0, len(logs))
	for _, logEntry := range logs {
		// Add an class to make the row red or green depending if the file is still open
		rowClass := "open"
		timestampClosed := "File has not closed"
		timeOpenMessage := "File has not closed yet"
		// Check if the file is closed
		if logEntry.TimeStampClosed.IsZero() == false {
			// Get time differnce between opening time and closing time
			timeOpen := logEntry.TimeStampClosed.Sub(logEntry.TimeStamp)
			// Check to see if the file was open for more than 2 miniutes because it was probaly an mistake if the file was open for less than 2 miniutes
			if timeOpen < time.Minute*2 {
				continue
			}

			// Update timeOpenMessage to the duration the file was open and format it to be in the format 3m 45s
			timeOpenMessage = formatDuration(logEntry.TimeStampClosed.Sub(logEntry.TimeStamp))
			// Change class to closed to make row red
			rowClass = "closed"
			timestampClosed = logEntry.TimeStampClosed.Format("2006-01-02 15:04:05")
		}

		// Create an FormattedLog which will be looped through in the index.html templated.
		log := models.FormatedLog{RowClass: rowClass,
			TimeStampFormatted:       string(logEntry.TimeStamp.Format("2006-01-02 15:04:05")),
			HostName:                 logEntry.HostName,
			FileName:                 logEntry.FileName,
			TimeStampClosedFormatted: timestampClosed,
			TimeFileWasOpened:        timeOpenMessage}
		// Append the created log to the list
		formatLogs = append(formatLogs, log)
	}
	return formatLogs
}
func HomeSuggestions(ctx *gin.Context) {
	var logs []models.LogEntry
	hostnames := make([]string, 0)
	hostname := ctx.Query("hostname")
	// Allow hostname to match any part of the hostname stored in the database
	hostname = "%" + hostname + "%" //
	// Find the logs that have the currently typed in query anywhere in the hostname
	initializers.DB.Where("host_name LIKE?", hostname).Find(&logs)
	for _, log := range logs {
		// Append the hostname to an array of hostnames that match the entered term
		hostnames = append(hostnames, log.HostName)
	}
	// Get rid of duplicates so only show in once as an dataentry and not many times
	hostnames = removeDuplicates(hostnames)
	// Load the suggestions template which creates the datalist which is used by htmx
	ctx.HTML(http.StatusOK, "suggestions.html", gin.H{"hostnames": hostnames})
}
func removeDuplicates(s []string) []string {
	bucket := make(map[string]bool)
	var result []string
	for _, str := range s {
		// Check if the key already exists in the map if and it isnt add the value to an list that is returned at the end of the function and do nothing if the key is already present
		if _, ok := bucket[str]; !ok {
			bucket[str] = true
			result = append(result, str)
		}
	}
	return result
}

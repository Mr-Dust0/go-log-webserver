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
	var logs []models.LogEntry
	// Find all the logs in the database
	err := initializers.DB.Find(&logs).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}
	// Get the html for the logs to be displayed in html
	userName, _ := ctx.Get("userName")
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"date":     "Showing all logs for all days",
		"userName": userName.(string),
	})
}

func GetHomePageHandler(ctx *gin.Context) {
	var logs []models.LogEntry
	// Find all the logs in the database
	err := initializers.DB.Find(&logs).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	// Get the html for the logs to be displayed in html
	formattedLogs := formatLogs(logs)
	ctx.HTML(http.StatusOK, "logtable.html", gin.H{
		"Logs": formattedLogs,
		"date": "Showing all logs for all days",
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
	hostname = "%" + hostname + "%" //
	initializers.DB.Where("host_name LIKE?", hostname).Find(&logs)
	for _, log := range logs {
		hostnames = append(hostnames, log.HostName)
	}
	hostnames = removeDuplicates(hostnames)
	fmt.Println(hostnames)
	fmt.Println("Hello")
	ctx.HTML(http.StatusOK, "suggestions.html", gin.H{"hostnames": hostnames})
}
func removeDuplicates(s []string) []string {
	bucket := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := bucket[str]; !ok {
			bucket[str] = true
			result = append(result, str)
		}
	}
	return result
}

func PostSearch(ctx *gin.Context) {
	// Gets data from user input
	date := ctx.PostForm("date")
	hostname := ctx.PostForm("hostname")
	var err error
	var logs []models.LogEntry
	// This is to be displayed to the user depending on what data was passed into the form
	var datemessage string
	if date == "" && hostname == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Date or HostName is required"})
		return
	}
	if date != "" {
		// Both hostname and data were entered
		if hostname != "" {
			// Get logs that match the search cateria of both hostname and date
			err = initializers.DB.Where("DATE(time_stamp) = ? AND host_name = ?", date, hostname).Find(&logs).Error
			// Update datemessage to show what search results lead to the output
			datemessage = "Showing results for the day: " + date + " and the hostname " + hostname
		} else {
			// Match enteries that match the date entered
			err = initializers.DB.Where("DATE(time_stamp) = ?", date).Find(&logs).Error
			datemessage = "Showing results for the day: " + date
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
			return
		}
	} else {
		// Match enteries for hostname entered
		err = initializers.DB.Where("host_name = ?", hostname).Find(&logs).Error
		datemessage = "Showing results for the hostname: " + hostname
	}
	fmt.Println(datemessage)
	formattedLogs := formatLogs(logs)
	ctx.HTML(http.StatusOK, "logtable.html", gin.H{
		"Logs": formattedLogs,
		"date": datemessage,
	})
}

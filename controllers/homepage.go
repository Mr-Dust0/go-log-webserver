package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
	"webserver/initializers"
	"webserver/models"

	"github.com/gin-gonic/gin"
)

type FormatedLog struct {
	RowClass                 string
	TimeStampFormatted       string
	HostName                 string
	FileName                 string
	TimeStampClosedFormatted string
	TimeFileWasOpened        string
}

func formatLogs(logs []models.LogEntry) string {
	html := ""
	for _, logEntry := range logs {
		// Add an class to make the row red or green depending if the file is still open
		rowClass := "class='open'"
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
			rowClass = "class='closed'"
			timestampClosed = logEntry.TimeStampClosed.Format("2006-01-02 15:04:05")
		}

		// Html for the table row that will be used as an variable in index.html
		html += fmt.Sprintf(`
        <tr %s>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
	    <td>%s</td>
        </tr>`, rowClass, logEntry.TimeStamp.Format("2006-01-02 15:04:05"), logEntry.HostName, logEntry.FileName, timestampClosed, timeOpenMessage)
	}
	return html
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
	html := formatLogs(logs)
	userName, _ := ctx.Get("userName")
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"data":     template.HTML(html),
		"date":     "Showing all logs for all days",
		"userName": "Welcome " + userName.(string),
	})
}
func PostHomePageHandler(ctx *gin.Context) {
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
	html := formatLogs(logs)
	userName, _ := ctx.Get("userName")
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"data":     template.HTML(html),
		"date":     datemessage,
		"userName": "Welcome " + userName.(string),
	})
}
func formatDuration(d time.Duration) string {
	// Convert the duration to minutes and seconds
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60

	// Format the result as "Xm Ys"
	return fmt.Sprintf("%dm %ds", minutes, seconds)
}

func GetHomePageHandler2(ctx *gin.Context) {
	var logs []models.LogEntry
	// Find all the logs in the database
	err := initializers.DB.Find(&logs).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	// Get the html for the logs to be displayed in html
	userName, _ := ctx.Get("userName")
	formatedLogs := formatLogs2(logs)
	fmt.Println(formatedLogs[0].FileName)
	ctx.HTML(http.StatusOK, "index2.html", gin.H{
		"Logs":     formatedLogs,
		"date":     "Showing all logs for all days",
		"userName": "Welcome " + userName.(string),
	})
}
func formatLogs2(logs []models.LogEntry) []FormatedLog {
	formatLogs := make([]FormatedLog, 0, len(logs))
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

		// Html for the table row that will be used as an variable in index.html
		log := FormatedLog{RowClass: rowClass,
			TimeStampFormatted:       string(logEntry.TimeStamp.Format("2006-01-02 15:04:05")),
			HostName:                 logEntry.HostName,
			FileName:                 logEntry.FileName,
			TimeStampClosedFormatted: timestampClosed,
			TimeFileWasOpened:        timeOpenMessage}
		formatLogs = append(formatLogs, log)
	}
	return formatLogs
}

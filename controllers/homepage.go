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

func formatLogs(logs []models.LogEntry) string {
	html := ""
	for _, logEntry := range logs {
		rowClass := "class='open'"
		timestampClosed := "File has not closed"
		timeOpenMessage := "File has not closed yet"
		fmt.Println(logEntry.TimeStampClosed)
		fmt.Println(logEntry.TimeStampClosed.IsZero())
		if logEntry.TimeStampClosed.IsZero() == false {
			timeOpen := logEntry.TimeStampClosed.Sub(logEntry.TimeStamp)
			fmt.Println(timeOpen)
			if timeOpen < time.Minute*2 {
				fmt.Println("The file was not open long enough for this to be logged")
				continue
			}

			timeOpenMessage = formatDuration(logEntry.TimeStampClosed.Sub(logEntry.TimeStamp))
			rowClass = "class='closed'"
			timestampClosed = logEntry.TimeStampClosed.Format("2006-01-02 15:04:05")
		}

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
	err := initializers.DB.Find(&logs).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	html := formatLogs(logs)
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"data": template.HTML(html),
		"date": "Showing all logs for all days",
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
		if hostname != "" {
			// Both hostname and data were entered
			err = initializers.DB.Where("DATE(time_stamp) = ? AND host_name = ?", date, hostname).Find(&logs).Error
			datemessage = "Showing results for the day: " + date + " and the hostname " + hostname
		} else {
			// Only date entered
			err = initializers.DB.Where("DATE(time_stamp) = ?", date).Find(&logs).Error
			datemessage = "Showing results for the day: " + date
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
			return
		}
	} else {
		// Only hostname entered
		err = initializers.DB.Where("host_name = ?", hostname).Find(&logs).Error
		datemessage = "Showing results for the hostname: " + hostname
	}
	html := formatLogs(logs)
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"data": template.HTML(html),
		"date": datemessage,
	})
}
func formatDuration(d time.Duration) string {
	// Convert the duration to minutes and seconds
	minutes := int(d.Minutes())      // Get the integer part of minutes
	seconds := int(d.Seconds()) % 60 // Get the integer part of seconds

	// Format the result as "Xm Ys"
	return fmt.Sprintf("%dm %ds", minutes, seconds)
}

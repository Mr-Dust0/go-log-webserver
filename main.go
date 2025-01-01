package main

import (
	_ "encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// LogEntry represents the database structure
type LogEntry struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TimeStamp       time.Time `json:"TimeStamp"`
	TimeStampClosed time.Time `json:"TimeStampClosed"` // Nullable field
	HostName        string    `json:"hostname"`
	FileName        string    `json:"filename"`
}

// Initialize the database connection
func initDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("/home/john/Downloads/dad-project/test/build/pdf_tracker/tracker.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = DB.AutoMigrate(&LogEntry{})
	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}
}

// Insert a log entry
func insertLog(ctx *gin.Context) {
	var logEntry LogEntry
	if err := ctx.ShouldBindJSON(&logEntry); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	logEntry.TimeStamp = time.Now()
	if err := DB.Create(&logEntry).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert log"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Log added successfully"})
}

// Update the `timeclosed` field
func updateTimeclosed(ctx *gin.Context) {
	var request struct {
		HostName   string    `json:"hostname"`
		TimeClosed time.Time `json:"TimeStampClosed"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	parsedDate := time.Now().Format("2006-01-02")
	result := DB.Model(&LogEntry{}).
		Where("host_name = ? AND DATE(time_stamp) = ? AND time_stamp_closed = '0001-01-01 00:00:00+00:00'", request.HostName, parsedDate).
		Update("time_stamp_closed", time.Now().String())
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update log"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "File closed successfully"})
}

// Fetch logs for a specific date
func logDateHandler(ctx *gin.Context) {
	date := ctx.PostForm("date")
	hostname := ctx.PostForm("hostname")
	fmt.Println(hostname)
	fmt.Println(date)
	var err error
	var logs []LogEntry
	var datemessage string
	if date == "" && hostname == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Date or HostName is required"})
		return
	}
	if date != "" {

		if hostname != "" {
			fmt.Println(hostname)
			fmt.Println(date)
			err = DB.Where("DATE(time_stamp) = ? AND host_name = ?", date, hostname).Find(&logs).Error
			datemessage = "Showing results for the day: " + date + " and the hostname " + hostname

		} else {
			err = DB.Where("DATE(time_stamp) = ?", date).Find(&logs).Error
			datemessage = "Showing results for the day: " + date
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
			return
		}
	} else {
		err = DB.Where("host_name = ?", hostname).Find(&logs).Error
		datemessage = "Showing results for the hostname: " + hostname
	}

	html := formatLogs(logs)
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"data": template.HTML(html),
		"date": datemessage,
	})
}

// Fetch all logs
func homePage(ctx *gin.Context) {
	var logs []LogEntry
	err := DB.Find(&logs).Error
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

// Format logs as HTML table rows
func formatLogs(logs []LogEntry) string {
	html := ""
	for _, logEntry := range logs {
		rowClass := "class='open'"
		timestampClosed := "File has not closed"
		timeOpenMessage := "File is still open"
		if logEntry.TimeStampClosed.IsZero() == false {
			timeOpen := logEntry.TimeStampClosed.Sub(logEntry.TimeStamp)
			if timeOpen < time.Minute*2 {
				fmt.Println("The file was not open long enough for this to be logged")
				continue
			}

			timeOpenMessage = logEntry.TimeStampClosed.Sub(logEntry.TimeStamp).String()
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

func main() {
	initDatabase()
	r := gin.Default()
	r.StaticFile("/favicon.ico", "./favicon.ico")
	r.LoadHTMLGlob("templates/*")
	r.GET("/", homePage)
	r.POST("/", logDateHandler)
	r.POST("/log", insertLog)
	r.PUT("/fileclosed", updateTimeclosed)
	err := r.RunTLS(":443", "/home/john/Downloads/go/web/hello_world/cert.pem", "/home/john/Downloads/go/web/hello_world/cert-key.pem")
	if err != nil {
		log.Fatal(err)
	}
}

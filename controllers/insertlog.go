package controllers

import (
	"fmt"
	"net/http"
	"time"
	"webserver/initializers"
	"webserver/models"

	"github.com/gin-gonic/gin"
)

func InsertLog(ctx *gin.Context) {
	var logEntry models.LogEntry
	// Get Json data from qr code reader request and store the data into logEntry
	if err := ctx.ShouldBindJSON(&logEntry); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	logEntry.TimeStamp = time.Now()
	// Check to see if the hostname already has an file open and if so update the file to be closed this cloud happen if the pi was shut down without sending an closed request for the file it had open or the pi not having internet when the file was closed so need to check this.
	err := initializers.DB.Where("time_stamp_closed = ? AND host_name = ?", "0001-01-01 00:00:00+00:00", logEntry.HostName).Update("time_stamp_closed", logEntry.TimeStamp)
	if err != nil {
		fmt.Println("Couldnt check whever that hostname already has an file open")
	}
	// Generate Timestamp the file was open here because the pi time can be off
	// Create an new entry in the database for the log
	if err := initializers.DB.Create(&logEntry).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert log"})
		return
	}

	// Return sucess message to the qr-code reader
	ctx.JSON(http.StatusOK, gin.H{"status": "Log added successfully"})
}

package controllers

import (
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
	// Generate Timestamp the file was open here because the pi time can be off
	logEntry.TimeStamp = time.Now()
	// Create an new entry in the database for the log
	if err := initializers.DB.Create(&logEntry).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert log"})
		return
	}

	// Return sucess message to the qr-code reader
	ctx.JSON(http.StatusOK, gin.H{"status": "Log added successfully"})
}

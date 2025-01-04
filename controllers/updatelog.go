package controllers

import (
	"net/http"
	"time"
	"webserver/initializers"
	"webserver/models"

	"github.com/gin-gonic/gin"
)

func UpdateTimeClosed(ctx *gin.Context) {
	var request struct {
		HostName string `json:"hostname"`
	}
	// Get json data from qr-code reqder and store it in request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Better to do the timestamp on the server side beacuse the time on the raspberry pi can be incorrect since it does not have an battery powered clock so time will be off until ntp fixes it.
	parsedDate := time.Now().Format("2006-01-02")
	// Update TImeStampClosed for matching hostname and where TimeStampClosed is still null
	result := initializers.DB.Model(&models.LogEntry{}).
		Where("host_name = ? AND DATE(time_stamp) = ? AND time_stamp_closed = '0001-01-01 00:00:00+00:00'", request.HostName, parsedDate).
		Update("time_stamp_closed", time.Now())
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update log"})
		return
	}
	// Send sucess message back to qr-code reader
	ctx.JSON(http.StatusOK, gin.H{"status": "File closed successfully"})
}

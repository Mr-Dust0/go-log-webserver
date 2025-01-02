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
	if err := ctx.ShouldBindJSON(&logEntry); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	logEntry.TimeStamp = time.Now()
	if err := initializers.DB.Create(&logEntry).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert log"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Log added successfully"})
}

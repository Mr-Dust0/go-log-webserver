package types

import (
	"time"

	"webserver/utils"
)

// Create LogEntry Structre and give metadata on the corresponding json keys and make ID primary key
type LogEntry struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TimeStamp       time.Time `json:"TimeStamp"`
	TimeStampClosed time.Time `json:"TimeStampClosed"` // Nullable field
	HostName        string    `json:"hostname"`
	FileName        string    `json:"filename"`
}

func (l *LogEntry) Format() FormattedLog {
	formattedLog := FormattedLog{
		RowClass:                 "open",
		TimeStampFormatted:       string(l.TimeStamp.Format("2006-01-02 15:04:05")),
		HostName:                 l.HostName,
		FileName:                 l.FileName,
		TimeStampClosedFormatted: "File has not closed",
		TimeFileWasOpened:        "File has not closed yet",
	}

	// Check if the file is closed
	if !l.TimeStampClosed.IsZero() {
		// Get time difference between opening time and closing time
		timeOpen := l.TimeStampClosed.Sub(l.TimeStamp)

		// Check to see if the file was open for more than 2 miniutes because it was probaly an mistake if the file was open for less than 2 miniutes
		if timeOpen >= time.Minute*2 {
			// Update timeOpenMessage to the duration the file was open and format it to be in the format 3m 45s
			formattedLog.TimeFileWasOpened = utils.FormatDuration(l.TimeStampClosed.Sub(l.TimeStamp))
			// Change class to closed to make row red
			formattedLog.RowClass = "closed"
			formattedLog.TimeStampClosedFormatted = l.TimeStampClosed.Format("2006-01-02 15:04:05")
		}
	}

	return formattedLog
}

type LogEntrySlice []LogEntry

func (ls LogEntrySlice) Format() []FormattedLog {

	// Create list of FormattedLogs with the capacity the same size of the logs which should reduce the amount of heap allocations that need to happen and speed up performance an tiny bit
	formatLogs := make([]FormattedLog, 0, len(ls))

	for _, logEntry := range ls {
		// Append the created log to the list
		formatLogs = append(formatLogs, logEntry.Format())
	}

	return formatLogs
}

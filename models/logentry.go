package models

import "time"

// Create LogEntry Structre and give metadata on the corresponding json keys and make ID primary key
type LogEntry struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	TimeStamp       time.Time `json:"TimeStamp"`
	TimeStampClosed time.Time `json:"TimeStampClosed"` // Nullable field
	HostName        string    `json:"hostname"`
	FileName        string    `json:"filename"`
}

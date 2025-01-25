package database

import (
	"time"

	"webserver/types"
)

// DB is a struct that represents a high-level database package
type DatabaseAPI interface {
	GetUserByID(id float64) (*types.User, error)
	GetUserByUsername(username string) (*types.User, error)
	UpdatePassword(user *types.User, newPassword []byte) error

	ListLogs(params ListLogsParams) (types.LogEntrySlice, error)
	UpdateLogCloseTime(logEntry *types.LogEntry, closedAt time.Time) error
	InsertLog(logEntry *types.LogEntry) error
}

type ListLogsParams struct {
	OnlyOpen        bool
	HostnameSimilar string
	Hostname        string
	Date            string
}

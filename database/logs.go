package database

import (
	"time"

	"webserver/types"
)

func (db *Database) ListLogs(params ListLogsParams) (types.LogEntrySlice, error) {
	var logs types.LogEntrySlice

	trx := db.Instance

	if params.OnlyOpen {
		trx = trx.Where("time_stamp_closed = ?", "0001-01-01 00:00:00+00:00")
	}

	if params.HostnameSimilar != "" {
		trx = trx.Where("host_name LIKE ?", "%"+params.HostnameSimilar+"%")
	}

	if params.Hostname != "" {
		trx = trx.Where("host_name = ?", params.Hostname)
	}

	if params.Date != "" {
		trx = trx.Where("DATE(time_stamp) = ?", params.Date)
	}

	err := trx.Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// Update log only allows updating the time_stamp_closed field
func (db *Database) UpdateLogCloseTime(logEntry *types.LogEntry, closedAt time.Time) error {
	parsedDate := closedAt.Format("2006-01-02")
	
	trx := db.Instance.Model(logEntry).Update("time_stamp_closed", parsedDate)
	if trx.Error != nil {
		return trx.Error
	}

	return nil
}

func (db *Database) InsertLog(logEntry *types.LogEntry) error {
	trx := db.Instance.Create(logEntry)
	if trx.Error != nil {
		return trx.Error
	}

	return nil
}
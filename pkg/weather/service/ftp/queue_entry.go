package ftp

import (
	"github.com/lmacrc/weather/pkg/sql/driver/sqlite"
	"github.com/lmacrc/weather/pkg/weather/service"
)

type QueueEntry struct {
	ID             uint `gorm:"primarykey"`
	CreatedAt      sqlite.Timestamp
	Due            sqlite.Timestamp
	LocalPath      string
	RemoteDir      string
	RemoteFilename string
	ExpiresAt      *sqlite.Timestamp // ExpiresAt specifies the time which the entry should no longer be uploaded and dropped from the queue.
	RemoveLocal    bool              // RemoveLocal indicates the file should be removed from the local filesystem after it has been successfully uploaded.
}

func (q QueueEntry) TableName() string { return "ftp_queue_entries" }

func NewFromFtpRequest(req service.FtpRequest) QueueEntry {
	var expiresAt *sqlite.Timestamp
	if req.ExpiresAt != nil {
		expiresAt = &sqlite.Timestamp{Time: (*req.ExpiresAt).UTC()}
	}
	return QueueEntry{
		LocalPath:      req.LocalPath,
		RemoteDir:      req.RemoteDir,
		RemoteFilename: req.RemoteFilename,
		ExpiresAt:      expiresAt,
		RemoveLocal:    req.RemoveLocal,
	}
}

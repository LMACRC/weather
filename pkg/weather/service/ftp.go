package service

import (
	"time"
)

type FtpRequest struct {
	LocalPath      string     // LocalPath is the full path to the file to be uploaded.
	RemoteDir      string     // RemoteDir is the name of the target directory on the remote system.
	RemoteFilename string     // RemoteFilename is the name of the target filename on the remote system.
	ExpiresAt      *time.Time // ExpiresAt specifies the time which the entry should no longer be uploaded and removed from the queue.
	RemoveLocal    bool       // RemoveLocal the LocalPath after successfully uploading to remote system.
}

type Ftp interface {
	Enqueue(req FtpRequest) error
}

package fs

import (
	"github.com/crawlab-team/goseaweedfs"
	"time"
)

type Manager interface {
	Init() (err error)
	Close() (err error)
	ListDir(remotePath string, isRecursive bool) (files []goseaweedfs.FilerFileInfo, err error)
	UploadFile(localPath, remotePath string, args ...interface{}) (err error)
	UploadDir(localPath, remotePath string, args ...interface{}) (err error)
	DownloadFile(remotePath, localPath string, args ...interface{}) (err error)
	DownloadDir(remotePath, localPath string, args ...interface{}) (err error)
	DeleteFile(remotePath string) (err error)
	DeleteDir(remotePath string) (err error)
	SyncLocalToRemote(localPath, remotePath string, args ...interface{}) (err error)
	SyncRemoteToLocal(remotePath, localPath string, args ...interface{}) (err error)
	GetFile(remotePath string, args ...interface{}) (data []byte, err error)
	GetFileInfo(remotePath string) (file *goseaweedfs.FilerFileInfo, err error)
	UpdateFile(remotePath string, data []byte, args ...interface{}) (err error)
	Exists(remotePath string, args ...interface{}) (ok bool, err error)
	SetFilerUrl(url string)
	SetFilerAuthKey(authKey string)
	SetTimeout(timeout time.Duration)
	SetWorkerNum(num int)
	SetRetryInterval(interval time.Duration)
	SetRetryNum(num int)
	SetMaxQps(qps int)
}

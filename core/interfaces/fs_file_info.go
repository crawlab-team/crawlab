package interfaces

import (
	"os"
	"time"
)

type FsFileInfo interface {
	GetName() string
	GetPath() string
	GetFullPath() string
	GetExtension() string
	GetIsDir() bool
	GetFileSize() int64
	GetModTime() time.Time
	GetMode() os.FileMode
	GetHash() string
	GetChildren() []FsFileInfo
}

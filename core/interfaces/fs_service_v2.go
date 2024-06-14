package interfaces

type FsServiceV2 interface {
	List(path string) (files []FsFileInfo, err error)
	GetFile(path string) (data []byte, err error)
	GetFileInfo(path string) (file FsFileInfo, err error)
	Save(path string, data []byte) (err error)
	CreateDir(path string) (err error)
	Rename(path, newPath string) (err error)
	Delete(path string) (err error)
	Copy(path, newPath string) (err error)
}

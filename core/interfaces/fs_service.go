package interfaces

import (
	cfs "github.com/crawlab-team/crawlab-fs"
	vcs "github.com/crawlab-team/crawlab-vcs"
)

type FsService interface {
	WithConfigPath
	List(path string, opts ...ServiceCrudOption) (files []FsFileInfo, err error)
	GetFile(path string, opts ...ServiceCrudOption) (data []byte, err error)
	GetFileInfo(path string, opts ...ServiceCrudOption) (file FsFileInfo, err error)
	Save(path string, data []byte, opts ...ServiceCrudOption) (err error)
	Rename(path, newPath string, opts ...ServiceCrudOption) (err error)
	Delete(path string, opts ...ServiceCrudOption) (err error)
	Copy(path, newPath string, opts ...ServiceCrudOption) (err error)
	Commit(msg string) (err error)
	SyncToFs(opts ...ServiceCrudOption) (err error)
	SyncToWorkspace() (err error)
	GetFsPath() (path string)
	SetFsPath(path string)
	GetWorkspacePath() (path string)
	SetWorkspacePath(path string)
	GetRepoPath() (path string)
	SetRepoPath(path string)
	GetFs() (fs cfs.Manager)
	GetGitClient() (c *vcs.GitClient)
}

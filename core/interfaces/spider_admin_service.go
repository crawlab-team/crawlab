package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SpiderAdminService interface {
	WithConfigPath
	Start() (err error)
	// Schedule a new task of the spider
	Schedule(id primitive.ObjectID, opts *SpiderRunOptions) (taskIds []primitive.ObjectID, err error)
	// Clone the spider
	Clone(id primitive.ObjectID, opts *SpiderCloneOptions) (err error)
	// Delete the spider
	Delete(id primitive.ObjectID) (err error)
	// SyncGit syncs all git repositories
	SyncGit() (err error)
	// SyncGitOne syncs one git repository
	SyncGitOne(g Git) (err error)
	// Export exports the spider and return zip file path
	Export(id primitive.ObjectID) (filePath string, err error)
}

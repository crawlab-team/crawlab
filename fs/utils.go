package fs

import (
	"fmt"
	"github.com/crawlab-team/goseaweedfs"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

func IsGitFile(file goseaweedfs.FileInfo) (res bool) {
	// skip .git
	res, err := regexp.MatchString("/?\\.git/", file.Path)
	if err != nil {
		return false
	}
	return res
}

func getCollectionAndTtlFromArgs(args ...interface{}) (collection, ttl string) {
	if len(args) > 0 {
		collection = args[0].(string)
	}
	if len(args) > 1 {
		ttl = args[1].(string)
	}
	return
}

func getUrlValuesFromArgs(args ...interface{}) (values url.Values) {
	if len(args) > 0 {
		values = args[0].(url.Values)
	}
	return values
}

func getFilesAndFilesMaps(f *goseaweedfs.Filer, localPath, remotePath string) (localFiles []goseaweedfs.FileInfo, remoteFiles []goseaweedfs.FilerFileInfo, localFilesMap map[string]goseaweedfs.FileInfo, remoteFilesMap map[string]goseaweedfs.FilerFileInfo, err error) {
	// declare maps
	localFilesMap = map[string]goseaweedfs.FileInfo{}
	remoteFilesMap = map[string]goseaweedfs.FilerFileInfo{}

	// cache local files info
	localFiles, err = goseaweedfs.ListLocalFilesRecursive(localPath)
	if err != nil {
		return localFiles, remoteFiles, localFilesMap, remoteFilesMap, err
	}
	for _, file := range localFiles {
		fileRemotePath := fmt.Sprintf("%s%s", remotePath, strings.Replace(file.Path, localPath, "", -1))
		localFilesMap[fileRemotePath] = file

		// directory
		dirRemotePath := filepath.Dir(fileRemotePath)
		_, ok := localFilesMap[dirRemotePath]
		if !ok {
			localFilesMap[dirRemotePath] = goseaweedfs.FileInfo{
				Name: filepath.Base(dirRemotePath),
				Path: dirRemotePath,
			}
		}
	}

	// cache remote files info
	remoteFiles, err = f.ListDirRecursive(remotePath)
	if err != nil {
		if err.Error() != FilerResponseNotFoundErrorMessage {
			return localFiles, remoteFiles, localFilesMap, remoteFilesMap, err
		}
		err = nil
	}
	remoteFiles = getFlattenRemoteFiles(remoteFiles)
	for _, file := range remoteFiles {
		remoteFilesMap[file.FullPath] = file
	}

	return
}

func getFlattenRemoteFiles(files []goseaweedfs.FilerFileInfo) (flattenFiles []goseaweedfs.FilerFileInfo) {
	flattenFiles = []goseaweedfs.FilerFileInfo{}
	for _, file := range files {
		flattenFiles = append(flattenFiles, file)

		if file.IsDir {
			flattenFiles = append(flattenFiles, getFlattenRemoteFiles(file.Children)...)
		}
	}
	return
}

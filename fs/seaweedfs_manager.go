package fs

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/crawlab-team/goseaweedfs"
	"github.com/emirpasic/gods/queues/linkedlistqueue"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type seaweedFsManagerFn func(params seaweedFsManagerParams) (res seaweedFsManagerResults)

type seaweedFsManagerHandle struct {
	params  seaweedFsManagerParams
	fn      seaweedFsManagerFn
	resChan chan seaweedFsManagerResults
}

type seaweedFsManagerParams struct {
	localPath   string
	remotePath  string
	isRecursive bool
	collection  string
	ttl         string
	urlValues   url.Values
	data        []byte
}

type seaweedFsManagerResults struct {
	files []goseaweedfs.FilerFileInfo
	file  *goseaweedfs.FilerFileInfo
	data  []byte
	ok    bool
	err   error
}

type SeaweedFsManager struct {
	// settings variables
	filerUrl      string
	timeout       time.Duration
	authKey       string
	workerNum     int
	retryNum      uint64
	retryInterval time.Duration
	maxQps        int

	// internals
	f      *goseaweedfs.Filer
	q      *linkedlistqueue.Queue
	cr     int
	ch     chan seaweedFsManagerHandle
	closed bool
}

func (m *SeaweedFsManager) Init() (err error) {
	// filer options
	var filerOpts []goseaweedfs.FilerOption

	// auth key
	if m.authKey != "" {
		filerOpts = append(filerOpts, goseaweedfs.WithFilerAuthKey(m.authKey))
	}

	// handle channel
	m.ch = make(chan seaweedFsManagerHandle, m.workerNum)

	// filer instance
	m.f, err = goseaweedfs.NewFiler(m.filerUrl, &http.Client{Timeout: m.timeout}, filerOpts...)
	if err != nil {
		return trace.TraceError(err)
	}

	// start async
	go m.start()

	return nil
}

func (m *SeaweedFsManager) Close() (err error) {
	m.closed = true
	if err := m.f.Close(); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (m *SeaweedFsManager) ListDir(remotePath string, isRecursive bool) (files []goseaweedfs.FilerFileInfo, err error) {
	params := seaweedFsManagerParams{
		remotePath:  remotePath,
		isRecursive: isRecursive,
	}
	res := m.process(params, m.listDir)
	return res.files, res.err
}

func (m *SeaweedFsManager) ListDirRecursive(remotePath string) (files []goseaweedfs.FilerFileInfo, err error) {
	params := seaweedFsManagerParams{
		remotePath: remotePath,
	}
	res := m.process(params, m.listDirRecursive)
	return res.files, res.err
}

func (m *SeaweedFsManager) UploadFile(localPath, remotePath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return trace.TraceError(err)
	}
	collection, ttl := getCollectionAndTtlFromArgs(args...)
	params := seaweedFsManagerParams{
		localPath:  localPath,
		remotePath: remotePath,
		collection: collection,
		ttl:        ttl,
	}
	res := m.process(params, m.uploadFile)
	return res.err
}

func (m *SeaweedFsManager) UploadDir(localPath, remotePath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return trace.TraceError(err)
	}
	collection, ttl := getCollectionAndTtlFromArgs(args...)
	params := seaweedFsManagerParams{
		localPath:  localPath,
		remotePath: remotePath,
		collection: collection,
		ttl:        ttl,
	}
	res := m.process(params, m.uploadDir)
	return res.err
}

func (m *SeaweedFsManager) DownloadFile(remotePath, localPath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return trace.TraceError(err)
	}
	collection, ttl := getCollectionAndTtlFromArgs(args...)
	urlValues := getUrlValuesFromArgs(args...)
	params := seaweedFsManagerParams{
		localPath:  localPath,
		remotePath: remotePath,
		collection: collection,
		ttl:        ttl,
		urlValues:  urlValues,
	}
	res := m.process(params, m.downloadFile)
	return res.err
}

func (m *SeaweedFsManager) DownloadDir(remotePath, localPath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return trace.TraceError(err)
	}
	params := seaweedFsManagerParams{
		remotePath: remotePath,
	}
	res := m.process(params, m.downloadDir)
	return res.err
}

func (m *SeaweedFsManager) DeleteFile(remotePath string) (err error) {
	params := seaweedFsManagerParams{
		remotePath: remotePath,
	}
	res := m.process(params, m.deleteFile)
	return res.err
}

func (m *SeaweedFsManager) DeleteDir(remotePath string) (err error) {
	params := seaweedFsManagerParams{
		remotePath: remotePath,
	}
	res := m.process(params, m.deleteDir)
	return res.err
}

func (m *SeaweedFsManager) SyncLocalToRemote(localPath, remotePath string, args ...interface{}) (err error) {
	localPath, err = filepath.Abs(localPath)
	if err != nil {
		return trace.TraceError(err)
	}
	collection, ttl := getCollectionAndTtlFromArgs(args...)
	params := seaweedFsManagerParams{
		localPath:  localPath,
		remotePath: remotePath,
		collection: collection,
		ttl:        ttl,
	}
	res := m.process(params, m.syncLocalToRemote)
	return res.err
}

func (m *SeaweedFsManager) SyncRemoteToLocal(remotePath, localPath string, args ...interface{}) (err error) {
	collection, ttl := getCollectionAndTtlFromArgs(args...)
	params := seaweedFsManagerParams{
		localPath:  localPath,
		remotePath: remotePath,
		collection: collection,
		ttl:        ttl,
	}
	res := m.process(params, m.syncRemoteToLocal)
	return res.err
}

func (m *SeaweedFsManager) GetFile(remotePath string, args ...interface{}) (data []byte, err error) {
	urlValues := getUrlValuesFromArgs(args...)
	params := seaweedFsManagerParams{
		remotePath: remotePath,
		urlValues:  urlValues,
	}
	res := m.process(params, m.getFile)
	return res.data, res.err
}

func (m *SeaweedFsManager) GetFileInfo(remotePath string) (file *goseaweedfs.FilerFileInfo, err error) {
	params := seaweedFsManagerParams{
		remotePath: remotePath,
	}
	res := m.process(params, m.getFileInfo)
	return res.file, res.err
}

func (m *SeaweedFsManager) UpdateFile(remotePath string, data []byte, args ...interface{}) (err error) {
	collection, ttl := getCollectionAndTtlFromArgs(args...)
	params := seaweedFsManagerParams{
		remotePath: remotePath,
		collection: collection,
		ttl:        ttl,
		data:       data,
	}
	res := m.process(params, m.updateFile)
	return res.err
}

func (m *SeaweedFsManager) Exists(remotePath string, args ...interface{}) (ok bool, err error) {
	_, err = m.GetFile(remotePath, args...)
	if err == nil {
		// exists
		return true, nil
	}
	if strings.Contains(err.Error(), FilerStatusNotFoundErrorMessage) {
		// not exists
		return false, nil
	}
	return ok, trace.TraceError(err)
}

func (m *SeaweedFsManager) SetFilerUrl(url string) {
	m.filerUrl = url
}

func (m *SeaweedFsManager) SetFilerAuthKey(authKey string) {
	m.authKey = authKey
}

func (m *SeaweedFsManager) SetTimeout(timeout time.Duration) {
	m.timeout = timeout
}

func (m *SeaweedFsManager) SetWorkerNum(num int) {
	m.workerNum = num
}

func (m *SeaweedFsManager) SetRetryInterval(interval time.Duration) {
	m.retryInterval = interval
}

func (m *SeaweedFsManager) SetRetryNum(num int) {
	m.retryNum = uint64(num)
}

func (m *SeaweedFsManager) SetMaxQps(qps int) {
	m.maxQps = qps
}

func (m *SeaweedFsManager) newHandle(params seaweedFsManagerParams, fn seaweedFsManagerFn) (handle seaweedFsManagerHandle) {
	return seaweedFsManagerHandle{
		params:  params,
		fn:      fn,
		resChan: make(chan seaweedFsManagerResults),
	}
}

func (m *SeaweedFsManager) start() {
	for {
		if m.closed {
			return
		}
		handle := <-m.ch
		go func() {
			if err := backoff.Retry(func() error {
				res := handle.fn(handle.params)
				if res.err != nil {
					return res.err
				}
				handle.resChan <- res
				return nil
			}, backoff.WithMaxRetries(
				backoff.NewConstantBackOff(m.retryInterval), m.retryNum),
			); err != nil {
				handle.resChan <- m.error(err)
			}
			m.wait()
		}()
	}
}

func (m *SeaweedFsManager) process(params seaweedFsManagerParams, fn seaweedFsManagerFn) (res seaweedFsManagerResults) {
	handle := m.newHandle(params, fn)
	//log.Infof("handle: %v", handle)
	m.ch <- handle
	res = <-handle.resChan
	return
}

func (m *SeaweedFsManager) error(err error) (res seaweedFsManagerResults) {
	if err != nil {
		trace.PrintError(err)
	}
	return seaweedFsManagerResults{err: err}
}

func (m *SeaweedFsManager) wait() {
	ms := float32(1) / float32(m.maxQps) * 1e3
	d := time.Duration(ms) * time.Millisecond
	time.Sleep(d)
}

func (m *SeaweedFsManager) getCollectionAndTtlArgsFromParams(params seaweedFsManagerParams) (args []interface{}) {
	if params.collection != "" {
		args = append(args, params.collection)
	}
	if params.ttl != "" {
		args = append(args, params.ttl)
	}
	return args
}

func (m *SeaweedFsManager) listDir(params seaweedFsManagerParams) (res seaweedFsManagerResults) {
	var err error
	if params.isRecursive {
		res.files, err = m.ListDirRecursive(params.remotePath)
	} else {
		res.files, err = m.f.ListDir(params.remotePath)
	}
	if err != nil {
		return m.error(err)
	}
	return
}

func (m *SeaweedFsManager) listDirRecursive(params seaweedFsManagerParams) (res seaweedFsManagerResults) {
	entries, err := m.f.ListDir(params.remotePath)
	if err != nil {
		return m.error(err)
	}
	for _, file := range entries {
		file = goseaweedfs.GetFileWithExtendedFields(file)
		if file.IsDir {
			file.Children, err = m.ListDirRecursive(file.FullPath)
			if err != nil {
				return m.error(err)
			}
		}
		res.files = append(res.files, file)
	}
	return
}

func (m *SeaweedFsManager) uploadFile(params seaweedFsManagerParams) (res seaweedFsManagerResults) {
	r, err := m.f.UploadFile(params.localPath, params.remotePath, params.collection, params.ttl)
	if err != nil {
		return m.error(err)
	}
	if r.Error != "" {
		err = errors.New(r.Error)
		return m.error(err)
	}
	return
}

func (m *SeaweedFsManager) uploadDir(params seaweedFsManagerParams) (res seaweedFsManagerResults) {
	args := m.getCollectionAndTtlArgsFromParams(params)
	if strings.HasSuffix(params.localPath, "/") {
		params.localPath = params.localPath[:(len(params.localPath) - 1)]
	}
	if !strings.HasPrefix(params.remotePath, "/") {
		params.remotePath = "/" + params.remotePath
	}
	files, err := goseaweedfs.ListFilesRecursive(params.localPath)
	if err != nil {
		return m.error(err)
	}
	for _, info := range files {
		newFilePath := params.remotePath + strings.Replace(info.Path, params.localPath, "", -1)
		if err := m.UploadFile(info.Path, newFilePath, args...); err != nil {
			return m.error(err)
		}
	}
	return
}

func (m *SeaweedFsManager) downloadFile(params seaweedFsManagerParams) (res seaweedFsManagerResults) {
	err := m.f.Download(params.remotePath, params.urlValues, func(reader io.Reader) error {
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			return trace.TraceError(err)
		}
		dirPath := filepath.Dir(params.localPath)
		_, err = os.Stat(dirPath)
		if err != nil {
			// if not exists, create a new directory
			if err := os.MkdirAll(dirPath, DefaultDirMode); err != nil {
				return trace.TraceError(err)
			}
		}
		fileMode := DefaultFileMode
		fileInfo, err := os.Stat(params.localPath)
		if err == nil {
			// if file already exists, save file mode and remove it
			fileMode = fileInfo.Mode()
			if err := os.Remove(params.localPath); err != nil {
				return trace.TraceError(err)
			}
		}
		if err := ioutil.WriteFile(params.localPath, data, fileMode); err != nil {
			return trace.TraceError(err)
		}
		return nil
	})
	if err != nil {
		return m.error(err)
	}
	return
}

func (m *SeaweedFsManager) downloadDir(params seaweedFsManagerParams) (res seaweedFsManagerResults) {
	args := m.getCollectionAndTtlArgsFromParams(params)
	var files []goseaweedfs.FilerFileInfo
	files, res.err = m.ListDir(params.remotePath, true)
	for _, file := range files {
		if file.IsDir {
			if err := m.DownloadDir(file.FullPath, path.Join(params.localPath, file.Name), args...); err != nil {
				return m.error(err)
			}
		} else {
			if err := m.DownloadFile(file.FullPath, path.Join(params.localPath, file.Name), args...); err != nil {
				return m.error(err)
			}
		}
	}
	return
}

func (m *SeaweedFsManager) deleteFile(params seaweedFsManagerParams) (res seaweedFsManagerResults) {
	if err := m.f.DeleteFile(params.remotePath); err != nil {
		return m.error(err)
	}
	return
}

func (m *SeaweedFsManager) deleteDir(params seaweedFsManagerParams) (res seaweedFsManagerResults) {
	if err := m.f.DeleteDir(params.remotePath); err != nil {
		return m.error(err)
	}
	return
}

func (m *SeaweedFsManager) syncLocalToRemote(params seaweedFsManagerParams) (res seaweedFsManagerResults) {
	// args
	args := m.getCollectionAndTtlArgsFromParams(params)

	// raise error if local path does not exist
	if _, err := os.Stat(params.localPath); err != nil {
		return m.error(err)
	}

	// get files and maps
	localFiles, remoteFiles, localFilesMap, remoteFilesMap, err := getFilesAndFilesMaps(m.f, params.localPath, params.remotePath)
	if err != nil {
		return m.error(err)
	}

	// compare remote files with local files and delete files absent in local files
	for _, remoteFile := range remoteFiles {
		// attempt to get corresponding local file
		_, ok := localFilesMap[remoteFile.FullPath]

		if !ok {
			// file does not exist on local, delete
			if remoteFile.IsDir {
				if err := m.DeleteDir(remoteFile.FullPath); err != nil {
					return m.error(err)
				}
			} else {
				if err := m.DeleteFile(remoteFile.FullPath); err != nil {
					return m.error(err)
				}
			}
		}
	}

	// compare local files with remote files and upload files with difference
	for _, localFile := range localFiles {
		// skip .git
		if IsGitFile(localFile) {
			continue
		}

		// corresponding remote file path
		fileRemotePath := fmt.Sprintf("%s%s", params.remotePath, strings.Replace(localFile.Path, params.localPath, "", -1))

		// attempt to get corresponding remote file
		remoteFile, ok := remoteFilesMap[fileRemotePath]

		if !ok {
			// file does not exist on remote, upload
			if err := m.UploadFile(localFile.Path, fileRemotePath, args...); err != nil {
				return m.error(err)
			}
		} else {
			// file exists on remote, upload if md5sum values are different
			if remoteFile.Md5 != localFile.Md5 {
				if err := m.UploadFile(localFile.Path, fileRemotePath, args...); err != nil {
					return m.error(err)
				}
			}
		}
	}

	return
}

func (m *SeaweedFsManager) syncRemoteToLocal(params seaweedFsManagerParams) (res seaweedFsManagerResults) {
	// args
	args := m.getCollectionAndTtlArgsFromParams(params)

	// create directory if local path does not exist
	if _, err := os.Stat(params.localPath); err != nil {
		if err := os.MkdirAll(params.localPath, os.ModePerm); err != nil {
			return m.error(err)
		}
	}

	// get files and maps
	localFiles, remoteFiles, localFilesMap, remoteFilesMap, err := getFilesAndFilesMaps(m.f, params.localPath, params.remotePath)
	if err != nil {
		return m.error(err)
	}

	// compare local files with remote files and delete files absent on remote
	for _, localFile := range localFiles {
		// skip .git
		if IsGitFile(localFile) {
			continue
		}

		// corresponding remote file path
		fileRemotePath := fmt.Sprintf("%s%s", params.remotePath, strings.Replace(localFile.Path, params.localPath, "", -1))

		// attempt to get corresponding remote file
		_, ok := remoteFilesMap[fileRemotePath]

		if !ok {
			// file does not exist on remote, upload
			if err := os.Remove(localFile.Path); err != nil {
				return m.error(err)
			}
		}
	}

	// compare remote files with local files and download if files with difference
	for _, remoteFile := range remoteFiles {
		// directory
		if remoteFile.IsDir {
			localDirRelativePath := strings.Replace(remoteFile.FullPath, params.remotePath, "", 1)
			localDirPath := fmt.Sprintf("%s%s", params.localPath, localDirRelativePath)
			if err := m.SyncRemoteToLocal(remoteFile.FullPath, localDirPath); err != nil {
				return m.error(err)
			}
			continue
		}

		// local file path
		localFileRelativePath := strings.Replace(remoteFile.FullPath, params.remotePath, "", 1)
		localFilePath := fmt.Sprintf("%s%s", params.localPath, localFileRelativePath)

		// attempt to get corresponding local file
		localFile, ok := localFilesMap[remoteFile.FullPath]

		if !ok {
			// file does not exist on local, download
			if err := m.DownloadFile(remoteFile.FullPath, localFilePath); err != nil {
				return m.error(err)
			}
		} else {
			// file exists on remote, download if md5sum values are different
			if remoteFile.Md5 != localFile.Md5 {
				if err := m.DownloadFile(remoteFile.FullPath, localFilePath, args...); err != nil {
					return m.error(err)
				}
			}
		}
	}

	return
}

func (m *SeaweedFsManager) getFile(params seaweedFsManagerParams) (res seaweedFsManagerResults) {
	var buf bytes.Buffer
	res.err = m.f.Download(params.remotePath, params.urlValues, func(reader io.Reader) error {
		_, err := io.Copy(&buf, reader)
		if err != nil {
			return trace.TraceError(err)
		}
		return nil
	})
	res.data = buf.Bytes()
	return
}

func (m *SeaweedFsManager) getFileInfo(params seaweedFsManagerParams) (res seaweedFsManagerResults) {
	arr := strings.Split(params.remotePath, "/")
	dirName := strings.Join(arr[:(len(arr)-1)], "/")
	files, err := m.f.ListDir(dirName)
	if err != nil {
		return m.error(err)
	}
	for _, f := range files {
		if f.FullPath == params.remotePath {
			res.file = &f
			return
		}
	}
	return m.error(ErrorFsNotExists)
}

func (m *SeaweedFsManager) updateFile(params seaweedFsManagerParams) (res seaweedFsManagerResults) {
	tmpRootDir := os.TempDir()
	tmpDirPath := path.Join(tmpRootDir, ".seaweedfs")
	if _, err := os.Stat(tmpDirPath); err != nil {
		if err := os.MkdirAll(tmpDirPath, os.ModePerm); err != nil {
			return m.error(err)
		}
	}
	tmpFilePath := path.Join(tmpDirPath, fmt.Sprintf(".%s", uuid.New().String()))
	if _, err := os.Stat(tmpFilePath); err == nil {
		if err := os.Remove(tmpFilePath); err != nil {
			return m.error(err)
		}
	}
	if err := ioutil.WriteFile(tmpFilePath, params.data, os.ModePerm); err != nil {
		return m.error(err)
	}
	params2 := seaweedFsManagerParams{
		localPath:  tmpFilePath,
		remotePath: params.remotePath,
		collection: params.collection,
		ttl:        params.ttl,
	}
	if res := m.uploadFile(params2); res.err != nil {
		return m.error(res.err)
	}
	if err := os.Remove(tmpFilePath); err != nil {
		return m.error(err)
	}
	return
}

func NewSeaweedFsManager(opts ...Option) (m2 Manager, err error) {
	// manager
	m := &SeaweedFsManager{
		filerUrl:      "http://localhost:8888",
		timeout:       5 * time.Minute,
		workerNum:     1,
		retryInterval: 500 * time.Millisecond,
		retryNum:      3,
		maxQps:        5,
		q:             linkedlistqueue.New(),
	}

	// apply options
	for _, opt := range opts {
		opt(m)
	}

	// initialize
	if err := m.Init(); err != nil {
		return nil, err
	}

	return m, nil
}

var _seaweedFsManager Manager

func GetSeaweedFsManager(opts ...Option) (m2 Manager, err error) {
	if _seaweedFsManager == nil {
		_seaweedFsManager, err = NewSeaweedFsManager(opts...)
		if err != nil {
			return nil, err
		}
	}
	return _seaweedFsManager, nil
}

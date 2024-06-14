package handler

import (
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRunner struct {
	mock.Mock
	Runner
}

func (m *MockRunner) downloadFile(url string, filePath string) error {
	args := m.Called(url, filePath)
	return args.Error(0)
}

func newMockRunner() *MockRunner {
	r := &MockRunner{}
	r.s = &models.Spider{}
	return r
}

func TestSyncFiles_SuccessWithDummyFiles(t *testing.T) {
	// Create a test server that responds with a list of files
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/scan") {
			w.Write([]byte(`{"file1.txt":{"path": "file1.txt", "hash": "hash1"}, "file2.txt":{"path": "file2.txt", "hash": "hash2"}}`))
			return
		}

		if strings.HasSuffix(r.URL.Path, "/download") {
			w.Write([]byte("file content"))
			return
		}
	}))
	defer ts.Close()

	// Create a mock runner
	r := newMockRunner()
	r.On("downloadFile", mock.Anything, mock.Anything).Return(nil)

	// Set the master URL to the test server URL
	viper.Set("api.endpoint", ts.URL)

	localPath := filepath.Join(os.TempDir(), uuid.New().String())
	os.MkdirAll(filepath.Join(localPath, r.s.GetId().Hex()), os.ModePerm)
	defer os.RemoveAll(localPath)
	viper.Set("workspace", localPath)

	// Call the method under test
	err := r.syncFiles()
	assert.NoError(t, err)

	// Assert that the files were downloaded
	assert.FileExists(t, filepath.Join(localPath, r.s.GetId().Hex(), "file1.txt"))
	assert.FileExists(t, filepath.Join(localPath, r.s.GetId().Hex(), "file2.txt"))
}

func TestSyncFiles_DeletesFilesNotOnMaster(t *testing.T) {
	// Create a test server that responds with an empty list of files
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/scan") {
			w.Write([]byte(`{}`))
			return
		}
	}))
	defer ts.Close()

	// Create a mock runner
	r := newMockRunner()
	r.On("downloadFile", mock.Anything, mock.Anything).Return(nil)

	// Set the master URL to the test server URL
	viper.Set("api.endpoint", ts.URL)

	localPath := filepath.Join(os.TempDir(), uuid.New().String())
	os.MkdirAll(filepath.Join(localPath, r.s.GetId().Hex()), os.ModePerm)
	defer os.RemoveAll(localPath)
	viper.Set("workspace", localPath)

	// Create a dummy file that should be deleted
	dummyFilePath := filepath.Join(localPath, r.s.GetId().Hex(), "dummy.txt")
	dummyFile, _ := os.Create(dummyFilePath)
	dummyFile.Close()

	// Call the method under test
	err := r.syncFiles()
	assert.NoError(t, err)

	// Assert that the dummy file was deleted
	assert.NoFileExists(t, dummyFilePath)
}

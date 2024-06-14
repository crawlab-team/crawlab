package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceV2_List(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "fsTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir) // clean up

	testDir := filepath.Join(rootDir, "dir")
	os.Mkdir(testDir, 0755)
	ioutil.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("hello world"), 0644)
	ioutil.WriteFile(filepath.Join(testDir, "file2.txt"), []byte("hello again"), 0644)
	subDir := filepath.Join(testDir, "subdir")
	os.Mkdir(subDir, 0755)
	ioutil.WriteFile(filepath.Join(subDir, "file3.txt"), []byte("subdir file"), 0644)
	os.Mkdir(filepath.Join(testDir, "empty"), 0755) // explicitly testing empty dir inclusion

	svc := NewFsServiceV2(rootDir)

	files, err := svc.List("dir")
	if err != nil {
		t.Errorf("Failed to list files: %v", err)
	}

	// Assert correct number of items
	assert.Len(t, files, 4)
	// Use a map to verify presence and characteristics of files/directories to avoid order issues
	items := make(map[string]bool)
	for _, item := range files {
		items[item.GetName()] = item.GetIsDir()
	}

	_, file1Exists := items["file1.txt"]
	_, file2Exists := items["file2.txt"]
	_, subdirExists := items["subdir"]
	_, emptyExists := items["empty"]

	assert.True(t, file1Exists)
	assert.True(t, file2Exists)
	assert.True(t, subdirExists)
	assert.True(t, emptyExists) // Verify that the empty directory is included

	if subdirExists && len(files[2].GetChildren()) > 0 {
		assert.Equal(t, "file3.txt", files[2].GetChildren()[0].GetName())
	}
}

func TestServiceV2_GetFile(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "fsTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir) // clean up

	expectedContent := []byte("hello world")
	ioutil.WriteFile(filepath.Join(rootDir, "file.txt"), expectedContent, 0644)

	svc := NewFsServiceV2(rootDir)

	content, err := svc.GetFile("file.txt")
	if err != nil {
		t.Errorf("Failed to get file: %v", err)
	}
	assert.Equal(t, expectedContent, content)
}

func TestServiceV2_Delete(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "fsTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir) // clean up

	filePath := filepath.Join(rootDir, "file.txt")
	ioutil.WriteFile(filePath, []byte("hello world"), 0644)

	svc := NewFsServiceV2(rootDir)

	// Delete the file
	err = svc.Delete("file.txt")
	if err != nil {
		t.Errorf("Failed to delete file: %v", err)
	}

	// Verify deletion
	_, err = os.Stat(filePath)
	assert.True(t, os.IsNotExist(err))
}

func TestServiceV2_CreateDir(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "fsTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir) // clean up

	svc := NewFsServiceV2(rootDir)

	// Create a new directory
	err = svc.CreateDir("newDir")
	if err != nil {
		t.Errorf("Failed to create directory: %v", err)
	}

	// Verify the directory was created
	_, err = os.Stat(filepath.Join(rootDir, "newDir"))
	assert.NoError(t, err)
}

func TestServiceV2_Save(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "fsTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir) // clean up

	svc := NewFsServiceV2(rootDir)

	// Save a new file
	err = svc.Save("newFile.txt", []byte("Hello, world!"))
	if err != nil {
		t.Errorf("Failed to save file: %v", err)
	}

	// Verify the file was saved
	data, err := ioutil.ReadFile(filepath.Join(rootDir, "newFile.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", string(data))
}

func TestServiceV2_Rename(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "fsTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir) // clean up

	svc := NewFsServiceV2(rootDir)

	// Create a file to rename
	ioutil.WriteFile(filepath.Join(rootDir, "oldName.txt"), []byte("Hello, world!"), 0644)

	// Rename the file
	err = svc.Rename("oldName.txt", "newName.txt")
	if err != nil {
		t.Errorf("Failed to rename file: %v", err)
	}

	// Verify the file was renamed
	_, err = os.Stat(filepath.Join(rootDir, "newName.txt"))
	assert.NoError(t, err)
}

func TestServiceV2_RenameDir(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "fsTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir) // clean up

	svc := NewFsServiceV2(rootDir)

	// Create a directory to rename
	os.Mkdir(filepath.Join(rootDir, "oldName"), 0755)

	// Rename the directory
	err = svc.Rename("oldName", "newName")
	if err != nil {
		t.Errorf("Failed to rename directory: %v", err)
	}

	// Verify the directory was renamed
	_, err = os.Stat(filepath.Join(rootDir, "newName"))
	assert.NoError(t, err)
}

func TestServiceV2_Copy(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "fsTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir) // clean up

	svc := NewFsServiceV2(rootDir)

	// Create a file to copy
	ioutil.WriteFile(filepath.Join(rootDir, "source.txt"), []byte("Hello, world!"), 0644)

	// Copy the file
	err = svc.Copy("source.txt", "copy.txt")
	if err != nil {
		t.Errorf("Failed to copy file: %v", err)
	}

	// Verify the file was copied
	data, err := ioutil.ReadFile(filepath.Join(rootDir, "copy.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", string(data))
}

func TestServiceV2_CopyDir(t *testing.T) {
	rootDir, err := ioutil.TempDir("", "fsTest")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(rootDir) // clean up

	svc := NewFsServiceV2(rootDir)

	// Create a directory to copy
	os.Mkdir(filepath.Join(rootDir, "sourceDir"), 0755)
	ioutil.WriteFile(filepath.Join(rootDir, "sourceDir", "file.txt"), []byte("Hello, world!"), 0644)

	// Copy the directory
	err = svc.Copy("sourceDir", "copyDir")
	if err != nil {
		t.Errorf("Failed to copy directory: %v", err)
	}

	// Verify the directory was copied
	_, err = os.Stat(filepath.Join(rootDir, "copyDir"))
	assert.NoError(t, err)

	// Verify the file inside the directory was copied
	data, err := ioutil.ReadFile(filepath.Join(rootDir, "copyDir", "file.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", string(data))
}

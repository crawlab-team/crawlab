package controllers

import (
	"errors"
	"fmt"
	log2 "github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/fs"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func GetBaseFileListDir(rootPath string, c *gin.Context) {
	path := c.Query("path")

	fsSvc, err := getBaseFileFsSvc(rootPath)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	files, err := fsSvc.List(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	HandleSuccessWithData(c, files)
}

func GetBaseFileFile(rootPath string, c *gin.Context) {
	path := c.Query("path")

	fsSvc, err := getBaseFileFsSvc(rootPath)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	data, err := fsSvc.GetFile(path)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, string(data))
}

func GetBaseFileFileInfo(rootPath string, c *gin.Context) {
	path := c.Query("path")

	fsSvc, err := getBaseFileFsSvc(rootPath)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	info, err := fsSvc.GetFileInfo(path)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccessWithData(c, info)
}

func PostBaseFileSaveFile(rootPath string, c *gin.Context) {
	fsSvc, err := getBaseFileFsSvc(rootPath)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	if c.GetHeader("Content-Type") == "application/json" {
		var payload struct {
			Path string `json:"path"`
			Data string `json:"data"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			HandleErrorBadRequest(c, err)
			return
		}
		if err := fsSvc.Save(payload.Path, []byte(payload.Data)); err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	} else {
		path, ok := c.GetPostForm("path")
		if !ok {
			HandleErrorBadRequest(c, errors.New("missing required field 'path'"))
			return
		}
		file, err := c.FormFile("file")
		if err != nil {
			HandleErrorBadRequest(c, err)
			return
		}
		f, err := file.Open()
		if err != nil {
			HandleErrorBadRequest(c, err)
			return
		}
		fileData, err := io.ReadAll(f)
		if err != nil {
			HandleErrorBadRequest(c, err)
			return
		}
		if err := fsSvc.Save(path, fileData); err != nil {
			HandleErrorInternalServerError(c, err)
			return
		}
	}

	HandleSuccess(c)
}

func PostBaseFileSaveFiles(rootPath string, c *gin.Context) {
	fsSvc, err := getBaseFileFsSvc(rootPath)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(form.File))
	for path := range form.File {
		go func(path string) {
			file, err := c.FormFile(path)
			if err != nil {
				log2.Warnf("invalid file header: %s", path)
				log2.Error(err.Error())
				wg.Done()
				return
			}
			f, err := file.Open()
			if err != nil {
				log2.Warnf("unable to open file: %s", path)
				log2.Error(err.Error())
				wg.Done()
				return
			}
			fileData, err := io.ReadAll(f)
			if err != nil {
				log2.Warnf("unable to read file: %s", path)
				log2.Error(err.Error())
				wg.Done()
				return
			}
			if err := fsSvc.Save(path, fileData); err != nil {
				log2.Warnf("unable to save file: %s", path)
				log2.Error(err.Error())
				wg.Done()
				return
			}
			wg.Done()
		}(path)
	}
	wg.Wait()

	HandleSuccess(c)
}

func PostBaseFileSaveDir(rootPath string, c *gin.Context) {
	var payload struct {
		Path    string `json:"path"`
		NewPath string `json:"new_path"`
		Data    string `json:"data"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	fsSvc, err := getBaseFileFsSvc(rootPath)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := fsSvc.CreateDir(payload.Path); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func PostBaseFileRenameFile(rootPath string, c *gin.Context) {
	var payload struct {
		Path    string `json:"path"`
		NewPath string `json:"new_path"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	fsSvc, err := getBaseFileFsSvc(rootPath)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := fsSvc.Rename(payload.Path, payload.NewPath); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
}

func DeleteBaseFileFile(rootPath string, c *gin.Context) {
	var payload struct {
		Path string `json:"path"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}
	if payload.Path == "~" {
		payload.Path = "."
	}

	fsSvc, err := getBaseFileFsSvc(rootPath)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := fsSvc.Delete(payload.Path); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	_, err = fsSvc.GetFileInfo(".")
	if err != nil {
		_ = fsSvc.CreateDir("/")
	}

	HandleSuccess(c)
}

func PostBaseFileCopyFile(rootPath string, c *gin.Context) {
	var payload struct {
		Path    string `json:"path"`
		NewPath string `json:"new_path"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	fsSvc, err := getBaseFileFsSvc(rootPath)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := fsSvc.Copy(payload.Path, payload.NewPath); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func PostBaseFileExport(rootPath string, c *gin.Context) {
	fsSvc, err := getBaseFileFsSvc(rootPath)
	if err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// zip file path
	zipFilePath, err := fsSvc.Export()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	// download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFilePath))
	c.File(zipFilePath)
}

func GetBaseFileFsSvc(rootPath string) (svc interfaces.FsServiceV2, err error) {
	return getBaseFileFsSvc(rootPath)
}

func getBaseFileFsSvc(rootPath string) (svc interfaces.FsServiceV2, err error) {
	workspacePath := viper.GetString("workspace")
	fsSvc := fs.NewFsServiceV2(filepath.Join(workspacePath, rootPath))

	return fsSvc, nil
}

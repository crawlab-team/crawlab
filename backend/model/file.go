package model

import (
	"crawlab/utils"
	"github.com/apex/log"
	"os"
)

type File struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
	Size  int64  `json:"size"`
}

func RemoveFile(path string) error {
	if !utils.Exists(path) {
		log.Info("file not found: " + path)
		return nil
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}

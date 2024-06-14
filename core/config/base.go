package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"path/filepath"
)

var HomeDirPath, _ = homedir.Dir()

const configDirName = ".crawlab"

const configName = "config.json"

func GetConfigPath() string {
	if viper.GetString("metadata") != "" {
		MetadataPath := viper.GetString("metadata")
		return filepath.Join(MetadataPath, configName)
	}
	return filepath.Join(HomeDirPath, configDirName, configName)
}

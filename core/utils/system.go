package utils

import "github.com/spf13/viper"

func IsPro() bool {
	return viper.GetString("info.edition") == "global.edition.pro"
}

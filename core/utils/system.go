package utils

import "github.com/spf13/viper"

func IsPro() bool {
	return viper.GetString("edition") == "global.edition.pro"
}

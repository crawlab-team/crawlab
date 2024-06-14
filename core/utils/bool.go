package utils

import "github.com/spf13/viper"

func EnvIsTrue(key string, defaultOk bool) bool {
	isTrueBool := viper.GetBool(key)
	isTrueString := viper.GetString(key)
	if isTrueString == "" {
		return defaultOk
	}
	return isTrueBool || isTrueString == "Y"
}

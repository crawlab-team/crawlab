package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

func IsDebug() bool {
	return viper.GetBool("debug")
}

func LogDebug(msg string) {
	if !IsDebug() {
		return
	}
	fmt.Println(fmt.Sprintf("[DEBUG] %s: %s", time.Now().Format("2006-01-02 15:04:05"), msg))
}

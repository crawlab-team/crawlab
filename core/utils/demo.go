package utils

import (
	"fmt"
	"github.com/crawlab-team/crawlab/core/sys_exec"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/spf13/viper"
)

func GetApiAddress() (res string) {
	apiAddress := viper.GetString("api.address")
	if apiAddress == "" {
		return "http://localhost:8000"
	}
	return apiAddress
}

func IsDemo() (ok bool) {
	return EnvIsTrue("demo", true)
}

func InitializedDemo() (ok bool) {
	col := mongo.GetMongoCol("users")
	n, err := col.Count(nil)
	if err != nil {
		return true
	}
	return n > 0
}

func ImportDemo() (err error) {
	cmdStr := fmt.Sprintf("crawlab-cli login -a %s && crawlab-demo import", GetApiAddress())
	cmd := sys_exec.BuildCmd(cmdStr)
	if err := cmd.Run(); err != nil {
		trace.PrintError(err)
	}
	return nil
}

func ReimportDemo() (err error) {
	cmdStr := fmt.Sprintf("crawlab-cli login -a %s && crawlab-demo reimport", GetApiAddress())
	cmd := sys_exec.BuildCmd(cmdStr)
	if err := cmd.Run(); err != nil {
		trace.PrintError(err)
	}
	return nil
}

func CleanupDemo() (err error) {
	cmdStr := fmt.Sprintf("crawlab-cli login -a %s && crawlab-demo reimport", GetApiAddress())
	cmd := sys_exec.BuildCmd(cmdStr)
	if err := cmd.Run(); err != nil {
		trace.PrintError(err)
	}
	return nil
}

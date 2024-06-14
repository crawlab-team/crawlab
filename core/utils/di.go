package utils

import (
	"github.com/crawlab-team/go-trace"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"os"
)

func VisualizeContainer(c *dig.Container) (err error) {
	if !viper.GetBool("debug.di.visualize") {
		return nil
	}
	if err := dig.Visualize(c, os.Stdout); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

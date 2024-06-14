package cmd

import (
	"github.com/crawlab-team/crawlab/core/apps"
	"os"
	"testing"
)

func TestCmdServer(t *testing.T) {
	_ = os.Setenv("CRAWLAB_PPROF", "true")

	// app
	svr := apps.GetServerV2()

	// start
	apps.Start(svr)
}

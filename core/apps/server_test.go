package apps

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func init() {
	_ = os.Setenv("CRAWLAB_DEMO", "false")
}

func TestServer_Start(t *testing.T) {
	svr := GetServer()

	// start
	go Start(svr)
	time.Sleep(1 * time.Second)

	res, err := req.Get(fmt.Sprintf("http://localhost:%s/system-info", viper.GetString("server.port")))
	require.Nil(t, err)
	resStr, err := res.ToString()
	require.Nil(t, err)
	require.Contains(t, resStr, "success")
}

package apps

import (
	"bufio"
	"fmt"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/sys_exec"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/imroc/req"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Docker struct {
	// parent
	parent ServerApp

	// dependencies
	interfaces.WithConfigPath

	// seaweedfs log
	fsLogFilePath string
	fsLogFile     *os.File
	fsReady       bool
}

func (app *Docker) Init() {
	var err error
	app.fsLogFile, err = os.OpenFile(app.fsLogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(0777))
	if err != nil {
		trace.PrintError(err)
	}

	// replace paths
	if err := app.replacePaths(); err != nil {
		panic(err)
	}

	// start nginx
	go app.startNginx()

	// start seaweedfs
	go app.startSeaweedFs()
}

func (app *Docker) Start() {
	// import demo
	//if utils.IsDemo() && utils.InitializedDemo() {
	//	go app.importDemo()
	//}
}

func (app *Docker) Wait() {
	DefaultWait()
}

func (app *Docker) Stop() {
}

func (app *Docker) GetParent() (parent ServerApp) {
	return app.parent
}

func (app *Docker) SetParent(parent ServerApp) {
	app.parent = parent
}

func (app *Docker) Ready() (ok bool) {
	return app.fsReady &&
		app.parent.GetApi().Ready()
}

func (app *Docker) replacePaths() (err error) {
	// read
	indexHtmlPath := "/app/dist/index.html"
	indexHtmlBytes, err := os.ReadFile(indexHtmlPath)
	if err != nil {
		return trace.TraceError(err)
	}
	indexHtml := string(indexHtmlBytes)

	// replace paths
	baseUrl := viper.GetString("base.url")
	if baseUrl != "" {
		indexHtml = app._replacePath(indexHtml, "js", baseUrl)
		indexHtml = app._replacePath(indexHtml, "css", baseUrl)
		indexHtml = app._replacePath(indexHtml, "<link rel=\"stylesheet\" href=\"", baseUrl)
		indexHtml = app._replacePath(indexHtml, "<link rel=\"stylesheet\" href=\"", baseUrl)
		indexHtml = app._replacePath(indexHtml, "window.VUE_APP_API_BASE_URL = '", baseUrl)
	}

	// replace path of baidu tongji
	initBaiduTongji := viper.GetString("string")
	if initBaiduTongji != "" {
		indexHtml = strings.ReplaceAll(indexHtml, "window.VUE_APP_INIT_BAIDU_TONGJI = ''", fmt.Sprintf("window.VUE_APP_INIT_BAIDU_TONGJI = '%s'", initBaiduTongji))
	}

	// replace path of umeng
	initUmeng := viper.GetString("string")
	if initUmeng != "" {
		indexHtml = strings.ReplaceAll(indexHtml, "window.VUE_APP_INIT_UMENG = ''", fmt.Sprintf("window.VUE_APP_INIT_UMENG = '%s'", initUmeng))
	}

	// write
	if err := os.WriteFile(indexHtmlPath, []byte(indexHtml), os.FileMode(0766)); err != nil {
		return trace.TraceError(err)
	}

	return nil
}

func (app *Docker) _replacePath(text, path, baseUrl string) (res string) {
	text = strings.ReplaceAll(text, path, fmt.Sprintf("%s/%s", baseUrl, path))
	return text
}

func (app *Docker) startNginx() {
	cmd := exec.Command("service", "nginx", "start")
	sys_exec.ConfigureCmdLogging(cmd, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := fmt.Sprintf("[nginx] %s\n", scanner.Text())
			_, _ = os.Stdout.WriteString(line)
		}
	})
	if err := cmd.Run(); err != nil {
		trace.PrintError(err)
	}
}

func (app *Docker) startSeaweedFs() {
	seaweedFsDataPath := "/data/seaweedfs"
	if !utils.Exists(seaweedFsDataPath) {
		_ = os.MkdirAll(seaweedFsDataPath, os.FileMode(0777))
	}
	cmd := exec.Command("weed", "server",
		"-dir", "/data",
		"-master.dir", seaweedFsDataPath,
		"-volume.dir.idx", seaweedFsDataPath,
		"-ip", "localhost",
		"-volume.port", "9999",
		"-volume.minFreeSpace", "1GiB",
		"-filer",
	)
	sys_exec.ConfigureCmdLogging(cmd, func(scanner *bufio.Scanner) {
		for scanner.Scan() {
			line := fmt.Sprintf("[seaweedfs] %s\n", scanner.Text())
			_, _ = app.fsLogFile.WriteString(line)
		}
	})
	go func() {
		if err := cmd.Run(); err != nil {
			trace.PrintError(err)
		}
	}()
	for {
		_, err := req.Get("http://localhost:8888")
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		app.fsReady = true
		return
	}
}

func (app *Docker) importDemo() {
	for {
		if app.Ready() {
			break
		}
		time.Sleep(1 * time.Second)
	}
	_ = utils.ImportDemo()
}

func NewDocker(svr ServerApp) *Docker {
	dck := &Docker{
		parent:        svr,
		fsLogFilePath: "/var/log/weed.log",
	}

	dck.Init()

	return dck
}

var dck *Docker

func GetDocker(svr ServerApp) *Docker {
	if dck != nil {
		return dck
	}
	dck = NewDocker(svr)
	return dck
}

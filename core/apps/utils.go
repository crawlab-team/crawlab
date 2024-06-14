package apps

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/color"
	"github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/container"
	grpcclient "github.com/crawlab-team/crawlab/core/grpc/client"
	grpcserver "github.com/crawlab-team/crawlab/core/grpc/server"
	modelsclient "github.com/crawlab-team/crawlab/core/models/client"
	modelsservice "github.com/crawlab-team/crawlab/core/models/service"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/core/schedule"
	"github.com/crawlab-team/crawlab/core/spider/admin"
	"github.com/crawlab-team/crawlab/core/stats"
	"github.com/crawlab-team/crawlab/core/task/handler"
	"github.com/crawlab-team/crawlab/core/task/scheduler"
	taskstats "github.com/crawlab-team/crawlab/core/task/stats"
	"github.com/crawlab-team/crawlab/core/user"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/trace"
)

func Start(app App) {
	start(app)
}

func start(app App) {
	app.Init()
	go app.Start()
	app.Wait()
	app.Stop()
}

func DefaultWait() {
	utils.DefaultWait()
}

func initModule(name string, fn func() error) (err error) {
	if err := fn(); err != nil {
		log.Error(fmt.Sprintf("init %s error: %s", name, err.Error()))
		_ = trace.TraceError(err)
		panic(err)
	}
	log.Info(fmt.Sprintf("initialized %s successfully", name))
	return nil
}

func initApp(name string, app App) {
	_ = initModule(name, func() error {
		app.Init()
		return nil
	})
}

var injectors = []interface{}{
	modelsservice.GetService,
	modelsclient.NewServiceDelegate,
	modelsclient.NewNodeServiceDelegate,
	modelsclient.NewSpiderServiceDelegate,
	modelsclient.NewTaskServiceDelegate,
	modelsclient.NewTaskStatServiceDelegate,
	modelsclient.NewEnvironmentServiceDelegate,
	grpcclient.NewClient,
	grpcclient.NewPool,
	grpcserver.GetServer,
	grpcserver.NewModelDelegateServer,
	grpcserver.NewModelBaseServiceServer,
	grpcserver.NewNodeServer,
	grpcserver.NewTaskServer,
	grpcserver.NewMessageServer,
	config.NewConfigPathService,
	user.GetUserService,
	schedule.GetScheduleService,
	admin.GetSpiderAdminService,
	stats.GetStatsService,
	nodeconfig.NewNodeConfigService,
	taskstats.GetTaskStatsService,
	color.NewService,
	scheduler.GetTaskSchedulerService,
	handler.GetTaskHandlerService,
}

func injectModules() {
	c := container.GetContainer()
	for _, injector := range injectors {
		if err := c.Provide(injector); err != nil {
			panic(err)
		}
	}
}

package main

import (
	"context"
	"crawlab/config"
	"crawlab/database"
	"crawlab/lib/validate_bridge"
	"crawlab/middlewares"
	"crawlab/model"
	"crawlab/routes"
	"crawlab/services"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"
)

func main() {
	binding.Validator = new(validate_bridge.DefaultValidator)
	app := gin.Default()

	// 初始化配置
	if err := config.InitConfig(""); err != nil {
		log.Error("init config error:" + err.Error())
		panic(err)
	}
	log.Info("initialized config successfully")

	// 初始化日志设置
	logLevel := viper.GetString("log.level")
	if logLevel != "" {
		log.SetLevelFromString(logLevel)
	}
	log.Info("initialized log config successfully")
	if viper.GetString("log.isDeletePeriodically") == "Y" {
		err := services.InitDeleteLogPeriodically()
		if err != nil {
			log.Error("init DeletePeriodically failed")
			panic(err)
		}
		log.Info("initialized periodically cleaning log successfully")
	} else {
		log.Info("periodically cleaning log is switched off")
	}

	// 初始化Mongodb数据库
	if err := database.InitMongo(); err != nil {
		log.Error("init mongodb error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("initialized MongoDB successfully")

	// 初始化Redis数据库
	if err := database.InitRedis(); err != nil {
		log.Error("init redis error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("initialized Redis successfully")

	if model.IsMaster() {
		// 初始化定时任务
		if err := services.InitScheduler(); err != nil {
			log.Error("init scheduler error:" + err.Error())
			debug.PrintStack()
			panic(err)
		}
		log.Info("initialized schedule successfully")

		// 初始化用户服务
		if err := services.InitUserService(); err != nil {
			log.Error("init user service error:" + err.Error())
			debug.PrintStack()
			panic(err)
		}
		log.Info("initialized user service successfully")

		// 初始化依赖服务
		if err := services.InitDepsFetcher(); err != nil {
			log.Error("init dependency fetcher error:" + err.Error())
			debug.PrintStack()
			panic(err)
		}
		log.Info("initialized dependency fetcher successfully")
	}

	// 初始化任务执行器
	if err := services.InitTaskExecutor(); err != nil {
		log.Error("init task executor error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("initialized task executor successfully")

	// 初始化节点服务
	if err := services.InitNodeService(); err != nil {
		log.Error("init node service error:" + err.Error())
		panic(err)
	}
	log.Info("initialized node service successfully")

	// 初始化爬虫服务
	if err := services.InitSpiderService(); err != nil {
		log.Error("init spider service error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("initialized spider service successfully")

	// 初始化RPC服务
	if err := services.InitRpcService(); err != nil {
		log.Error("init rpc service error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("initialized rpc service successfully")

	// 以下为主节点服务
	if model.IsMaster() {
		// 中间件
		app.Use(middlewares.CORSMiddleware())
		anonymousGroup := app.Group("/")
		{
			anonymousGroup.POST("/login", routes.Login)       // 用户登录
			anonymousGroup.PUT("/users", routes.PutUser)      // 添加用户
			anonymousGroup.GET("/setting", routes.GetSetting) // 获取配置信息
			// release版本
			anonymousGroup.GET("/version", routes.GetVersion)               // 获取发布的版本
			anonymousGroup.GET("/releases/latest", routes.GetLatestRelease) // 获取最近发布的版本
		}
		authGroup := app.Group("/", middlewares.AuthorizationMiddleware())
		{
			// 节点
			{
				authGroup.GET("/nodes", routes.GetNodeList)                            // 节点列表
				authGroup.GET("/nodes/:id", routes.GetNode)                            // 节点详情
				authGroup.POST("/nodes/:id", routes.PostNode)                          // 修改节点
				authGroup.GET("/nodes/:id/tasks", routes.GetNodeTaskList)              // 节点任务列表
				authGroup.GET("/nodes/:id/system", routes.GetSystemInfo)               // 节点任务列表
				authGroup.DELETE("/nodes/:id", routes.DeleteNode)                      // 删除节点
				authGroup.GET("/nodes/:id/langs", routes.GetLangList)                  // 节点语言环境列表
				authGroup.GET("/nodes/:id/deps", routes.GetDepList)                    // 节点第三方依赖列表
				authGroup.GET("/nodes/:id/deps/installed", routes.GetInstalledDepList) // 节点已安装第三方依赖列表
				authGroup.POST("/nodes/:id/deps/install", routes.InstallDep)           // 节点安装依赖
				authGroup.POST("/nodes/:id/deps/uninstall", routes.UninstallDep)       // 节点卸载依赖
				authGroup.POST("/nodes/:id/langs/install", routes.InstallLang)         // 节点安装语言
			}
			// 爬虫
			{
				authGroup.GET("/spiders", routes.GetSpiderList)                                 // 爬虫列表
				authGroup.GET("/spiders/:id", routes.GetSpider)                                 // 爬虫详情
				authGroup.PUT("/spiders", routes.PutSpider)                                     // 添加爬虫
				authGroup.POST("/spiders", routes.UploadSpider)                                 // 上传爬虫
				authGroup.POST("/spiders/:id", routes.PostSpider)                               // 修改爬虫
				authGroup.POST("/spiders/:id/publish", routes.PublishSpider)                    // 发布爬虫
				authGroup.POST("/spiders/:id/upload", routes.UploadSpiderFromId)                // 上传爬虫（ID）
				authGroup.DELETE("/spiders/:id", routes.DeleteSpider)                           // 删除爬虫
				authGroup.GET("/spiders/:id/tasks", routes.GetSpiderTasks)                      // 爬虫任务列表
				authGroup.GET("/spiders/:id/file/tree", routes.GetSpiderFileTree)               // 爬虫文件目录树读取
				authGroup.GET("/spiders/:id/file", routes.GetSpiderFile)                        // 爬虫文件读取
				authGroup.POST("/spiders/:id/file", routes.PostSpiderFile)                      // 爬虫文件更改
				authGroup.PUT("/spiders/:id/file", routes.PutSpiderFile)                        // 爬虫文件创建
				authGroup.PUT("/spiders/:id/dir", routes.PutSpiderDir)                          // 爬虫目录创建
				authGroup.DELETE("/spiders/:id/file", routes.DeleteSpiderFile)                  // 爬虫文件删除
				authGroup.POST("/spiders/:id/file/rename", routes.RenameSpiderFile)             // 爬虫文件重命名
				authGroup.GET("/spiders/:id/dir", routes.GetSpiderDir)                          // 爬虫目录
				authGroup.GET("/spiders/:id/stats", routes.GetSpiderStats)                      // 爬虫统计数据
				authGroup.GET("/spiders/:id/schedules", routes.GetSpiderSchedules)              // 爬虫定时任务
				authGroup.GET("/spiders/:id/scrapy/spiders", routes.GetSpiderScrapySpiders)     // Scrapy 爬虫名称列表
				authGroup.PUT("/spiders/:id/scrapy/spiders", routes.PutSpiderScrapySpiders)     // Scrapy 爬虫创建爬虫
				authGroup.GET("/spiders/:id/scrapy/settings", routes.GetSpiderScrapySettings)   // Scrapy 爬虫设置
				authGroup.POST("/spiders/:id/scrapy/settings", routes.PostSpiderScrapySettings) // Scrapy 爬虫修改设置
				authGroup.GET("/spiders/:id/scrapy/items", routes.GetSpiderScrapyItems)         // Scrapy 爬虫 items
				authGroup.POST("/spiders/:id/scrapy/items", routes.PostSpiderScrapyItems)       // Scrapy 爬虫修改 items
				authGroup.GET("/spiders/:id/scrapy/pipelines", routes.GetSpiderScrapyPipelines)         // Scrapy 爬虫 pipelines
				authGroup.POST("/spiders/:id/git/sync", routes.PostSpiderSyncGit)               // 爬虫 Git 同步
				authGroup.POST("/spiders/:id/git/reset", routes.PostSpiderResetGit)             // 爬虫 Git 重置
			}
			// 可配置爬虫
			{
				authGroup.GET("/config_spiders/:id/config", routes.GetConfigSpiderConfig)           // 获取可配置爬虫配置
				authGroup.POST("/config_spiders/:id/config", routes.PostConfigSpiderConfig)         // 更改可配置爬虫配置
				authGroup.PUT("/config_spiders", routes.PutConfigSpider)                            // 添加可配置爬虫
				authGroup.POST("/config_spiders/:id", routes.PostConfigSpider)                      // 修改可配置爬虫
				authGroup.POST("/config_spiders/:id/upload", routes.UploadConfigSpider)             // 上传可配置爬虫
				authGroup.POST("/config_spiders/:id/spiderfile", routes.PostConfigSpiderSpiderfile) // 上传可配置爬虫
				authGroup.GET("/config_spiders_templates", routes.GetConfigSpiderTemplateList)      // 获取可配置爬虫模版列表
			}
			// 任务
			{
				authGroup.GET("/tasks", routes.GetTaskList)                                 // 任务列表
				authGroup.GET("/tasks/:id", routes.GetTask)                                 // 任务详情
				authGroup.PUT("/tasks", routes.PutTask)                                     // 派发任务
				authGroup.DELETE("/tasks/:id", routes.DeleteTask)                           // 删除任务
				authGroup.DELETE("/tasks_multiple", routes.DeleteMultipleTask)              // 删除多个任务
				authGroup.DELETE("/tasks_by_status", routes.DeleteTaskByStatus)             //删除指定状态的任务
				authGroup.POST("/tasks/:id/cancel", routes.CancelTask)                      // 取消任务
				authGroup.GET("/tasks/:id/log", routes.GetTaskLog)                          // 任务日志
				authGroup.GET("/tasks/:id/results", routes.GetTaskResults)                  // 任务结果
				authGroup.GET("/tasks/:id/results/download", routes.DownloadTaskResultsCsv) // 下载任务结果
			}
			// 定时任务
			{
				authGroup.GET("/schedules", routes.GetScheduleList)              // 定时任务列表
				authGroup.GET("/schedules/:id", routes.GetSchedule)              // 定时任务详情
				authGroup.PUT("/schedules", routes.PutSchedule)                  // 创建定时任务
				authGroup.POST("/schedules/:id", routes.PostSchedule)            // 修改定时任务
				authGroup.DELETE("/schedules/:id", routes.DeleteSchedule)        // 删除定时任务
				authGroup.POST("/schedules/:id/disable", routes.DisableSchedule) // 禁用定时任务
				authGroup.POST("/schedules/:id/enable", routes.EnableSchedule)   // 启用定时任务
			}
			// 用户
			{
				authGroup.GET("/users", routes.GetUserList)       // 用户列表
				authGroup.GET("/users/:id", routes.GetUser)       // 用户详情
				authGroup.POST("/users/:id", routes.PostUser)     // 更改用户
				authGroup.DELETE("/users/:id", routes.DeleteUser) // 删除用户
				authGroup.GET("/me", routes.GetMe)                // 获取自己账户
				authGroup.POST("/me", routes.PostMe)              // 修改自己账户
			}
			// 系统
			{
				authGroup.GET("/system/deps/:lang", routes.GetAllDepList)             // 节点所有第三方依赖列表
				authGroup.GET("/system/deps/:lang/:dep_name/json", routes.GetDepJson) // 节点第三方依赖JSON
			}
			// 全局变量
			{
				authGroup.GET("/variables", routes.GetVariableList)      // 列表
				authGroup.PUT("/variable", routes.PutVariable)           // 新增
				authGroup.POST("/variable/:id", routes.PostVariable)     //修改
				authGroup.DELETE("/variable/:id", routes.DeleteVariable) //删除
			}
			// 项目
			{
				authGroup.GET("/projects", routes.GetProjectList)       // 列表
				authGroup.GET("/projects/tags", routes.GetProjectTags)  // 项目标签
				authGroup.PUT("/projects", routes.PutProject)           //修改
				authGroup.POST("/projects/:id", routes.PostProject)     // 新增
				authGroup.DELETE("/projects/:id", routes.DeleteProject) //删除
			}
			// 统计数据
			authGroup.GET("/stats/home", routes.GetHomeStats) // 首页统计数据
			// 文件
			authGroup.GET("/file", routes.GetFile) // 获取文件
			// Git
			authGroup.GET("/git/branches", routes.GetGitBranches)       // 获取 Git 分支
			authGroup.GET("/git/public-key", routes.GetGitSshPublicKey) // 获取 SSH 公钥
		}

	}

	// 路由ping
	app.GET("/ping", routes.Ping)

	// 运行服务器
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	address := net.JoinHostPort(host, port)
	srv := &http.Server{
		Handler: app,
		Addr:    address,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Error("run server error:" + err.Error())
			} else {
				log.Info("server graceful down")
			}
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx2, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx2); err != nil {
		log.Error("run server error:" + err.Error())
	}
}

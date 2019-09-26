package main

import (
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
	"runtime/debug"
)

func main() {
	binding.Validator = new(validate_bridge.DefaultValidator)
	app := gin.Default()

	// 初始化配置
	if err := config.InitConfig(""); err != nil {
		log.Error("init config error:" + err.Error())
		panic(err)
	}
	log.Info("初始化配置成功")

	// 初始化日志设置
	logLevel := viper.GetString("log.level")
	if logLevel != "" {
		log.SetLevelFromString(logLevel)
	}
	log.Info("初始化日志设置成功")

	if viper.GetString("log.isDeletePeriodically") == "Y" {
		err := services.InitDeleteLogPeriodically()
		if err != nil {
			log.Error("Init DeletePeriodically Failed")
			panic(err)
		}
		log.Info("初始化定期清理日志配置成功")
	}

	// 初始化Mongodb数据库
	if err := database.InitMongo(); err != nil {
		log.Error("init mongodb error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("初始化Mongodb数据库成功")

	// 初始化Redis数据库
	if err := database.InitRedis(); err != nil {
		log.Error("init redis error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("初始化Redis数据库成功")

	if model.IsMaster() {
		// 初始化定时任务
		if err := services.InitScheduler(); err != nil {
			log.Error("init scheduler error:" + err.Error())
			debug.PrintStack()
			panic(err)
		}
		log.Info("初始化定时任务成功")
	}

	// 初始化任务执行器
	if err := services.InitTaskExecutor(); err != nil {
		log.Error("init task executor error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("初始化任务执行器成功")

	// 初始化节点服务
	if err := services.InitNodeService(); err != nil {
		log.Error("init node service error:" + err.Error())
		panic(err)
	}
	log.Info("初始化节点配置成功")

	// 初始化爬虫服务
	if err := services.InitSpiderService(); err != nil {
		log.Error("init spider service error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("初始化爬虫服务成功")

	// 初始化用户服务
	if err := services.InitUserService(); err != nil {
		log.Error("init user service error:" + err.Error())
		debug.PrintStack()
		panic(err)
	}
	log.Info("初始化用户服务成功")

	// 以下为主节点服务
	if model.IsMaster() {
		// 中间件
		app.Use(middlewares.CORSMiddleware())
		//app.Use(middlewares.AuthorizationMiddleware())
		anonymousGroup := app.Group("/")
		{
			anonymousGroup.POST("/login", routes.Login)  // 用户登录
			anonymousGroup.PUT("/users", routes.PutUser) // 添加用户

		}
		authGroup := app.Group("/", middlewares.AuthorizationMiddleware())
		{
			// 路由
			// 节点
			authGroup.GET("/nodes", routes.GetNodeList)               // 节点列表
			authGroup.GET("/nodes/:id", routes.GetNode)               // 节点详情
			authGroup.POST("/nodes/:id", routes.PostNode)             // 修改节点
			authGroup.GET("/nodes/:id/tasks", routes.GetNodeTaskList) // 节点任务列表
			authGroup.GET("/nodes/:id/system", routes.GetSystemInfo)  // 节点任务列表
			authGroup.DELETE("/nodes/:id", routes.DeleteNode)         // 删除节点
			// 爬虫
			authGroup.GET("/spiders", routes.GetSpiderList)              // 爬虫列表
			authGroup.GET("/spiders/:id", routes.GetSpider)              // 爬虫详情
			authGroup.POST("/spiders", routes.PutSpider)                 // 上传爬虫
			authGroup.POST("/spiders/:id", routes.PostSpider)            // 修改爬虫
			authGroup.POST("/spiders/:id/publish", routes.PublishSpider) // 发布爬虫
			authGroup.DELETE("/spiders/:id", routes.DeleteSpider)        // 删除爬虫
			authGroup.GET("/spiders/:id/tasks", routes.GetSpiderTasks)   // 爬虫任务列表
			authGroup.GET("/spiders/:id/file", routes.GetSpiderFile)     // 爬虫文件读取
			authGroup.POST("/spiders/:id/file", routes.PostSpiderFile)   // 爬虫目录写入
			authGroup.GET("/spiders/:id/dir", routes.GetSpiderDir)       // 爬虫目录
			authGroup.GET("/spiders/:id/stats", routes.GetSpiderStats)   // 爬虫统计数据
			// 任务
			authGroup.GET("/tasks", routes.GetTaskList)                                 // 任务列表
			authGroup.GET("/tasks/:id", routes.GetTask)                                 // 任务详情
			authGroup.PUT("/tasks", routes.PutTask)                                     // 派发任务
			authGroup.DELETE("/tasks/:id", routes.DeleteTask)                           // 删除任务
			authGroup.POST("/tasks/:id/cancel", routes.CancelTask)                      // 取消任务
			authGroup.GET("/tasks/:id/log", routes.GetTaskLog)                          // 任务日志
			authGroup.GET("/tasks/:id/results", routes.GetTaskResults)                  // 任务结果
			authGroup.GET("/tasks/:id/results/download", routes.DownloadTaskResultsCsv) // 下载任务结果
			// 定时任务
			authGroup.GET("/schedules", routes.GetScheduleList)       // 定时任务列表
			authGroup.GET("/schedules/:id", routes.GetSchedule)       // 定时任务详情
			authGroup.PUT("/schedules", routes.PutSchedule)           // 创建定时任务
			authGroup.POST("/schedules/:id", routes.PostSchedule)     // 修改定时任务
			authGroup.DELETE("/schedules/:id", routes.DeleteSchedule) // 删除定时任务
			// 统计数据
			authGroup.GET("/stats/home", routes.GetHomeStats) // 首页统计数据
			// 用户
			authGroup.GET("/users", routes.GetUserList)       // 用户列表
			authGroup.GET("/users/:id", routes.GetUser)       // 用户详情
			authGroup.POST("/users/:id", routes.PostUser)     // 更改用户
			authGroup.DELETE("/users/:id", routes.DeleteUser) // 删除用户
			authGroup.GET("/me", routes.GetMe)                // 获取自己账户
		}

	}

	// 路由ping
	app.GET("/ping", routes.Ping)

	// 运行服务器
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	if err := app.Run(host + ":" + port); err != nil {
		log.Error("run server error:" + err.Error())
		panic(err)
	}
}

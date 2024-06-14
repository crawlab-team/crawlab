package controllers

import (
	"github.com/crawlab-team/crawlab/core/middlewares"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RouterGroups struct {
	AuthGroup      *gin.RouterGroup
	AnonymousGroup *gin.RouterGroup
}

func NewRouterGroups(app *gin.Engine) (groups *RouterGroups) {
	return &RouterGroups{
		AuthGroup:      app.Group("/", middlewares.AuthorizationMiddlewareV2()),
		AnonymousGroup: app.Group("/"),
	}
}

func RegisterController[T any](group *gin.RouterGroup, basePath string, ctr *BaseControllerV2[T]) {
	actionPaths := make(map[string]bool)
	for _, action := range ctr.actions {
		group.Handle(action.Method, basePath+action.Path, action.HandlerFunc)
		path := basePath + action.Path
		key := action.Method + " - " + path
		actionPaths[key] = true
	}
	registerBuiltinHandler(group, http.MethodGet, basePath+"", ctr.GetList, actionPaths)
	registerBuiltinHandler(group, http.MethodGet, basePath+"/:id", ctr.GetById, actionPaths)
	registerBuiltinHandler(group, http.MethodPost, basePath+"", ctr.Post, actionPaths)
	registerBuiltinHandler(group, http.MethodPut, basePath+"/:id", ctr.PutById, actionPaths)
	registerBuiltinHandler(group, http.MethodPatch, basePath+"", ctr.PatchList, actionPaths)
	registerBuiltinHandler(group, http.MethodDelete, basePath+"/:id", ctr.DeleteById, actionPaths)
	registerBuiltinHandler(group, http.MethodDelete, basePath+"", ctr.DeleteList, actionPaths)
}

func RegisterActions(group *gin.RouterGroup, basePath string, actions []Action) {
	for _, action := range actions {
		group.Handle(action.Method, basePath+action.Path, action.HandlerFunc)
	}
}

func registerBuiltinHandler(group *gin.RouterGroup, method, path string, handlerFunc gin.HandlerFunc, existingActionPaths map[string]bool) {
	key := method + " - " + path
	_, ok := existingActionPaths[key]
	if ok {
		return
	}
	group.Handle(method, path, handlerFunc)
}

func InitRoutes(app *gin.Engine) (err error) {
	// routes groups
	groups := NewRouterGroups(app)

	RegisterController(groups.AuthGroup, "/data/collections", NewControllerV2[models.DataCollectionV2]())
	RegisterController(groups.AuthGroup, "/data-sources", NewControllerV2[models.DataSourceV2]())
	RegisterController(groups.AuthGroup, "/environments", NewControllerV2[models.EnvironmentV2]())
	RegisterController(groups.AuthGroup, "/gits", NewControllerV2[models.GitV2]())
	RegisterController(groups.AuthGroup, "/nodes", NewControllerV2[models.NodeV2]())
	RegisterController(groups.AuthGroup, "/notifications/settings", NewControllerV2[models.SettingV2]())
	RegisterController(groups.AuthGroup, "/permissions", NewControllerV2[models.PermissionV2]())
	RegisterController(groups.AuthGroup, "/projects", NewControllerV2[models.ProjectV2]())
	RegisterController(groups.AuthGroup, "/roles", NewControllerV2[models.RoleV2]())
	RegisterController(groups.AuthGroup, "/schedules", NewControllerV2[models.ScheduleV2](
		Action{
			Method:      http.MethodPost,
			Path:        "",
			HandlerFunc: PostSchedule,
		},
		Action{
			Method:      http.MethodPut,
			Path:        "/:id",
			HandlerFunc: PutScheduleById,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/enable",
			HandlerFunc: PostScheduleEnable,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/disable",
			HandlerFunc: PostScheduleDisable,
		},
	))
	RegisterController(groups.AuthGroup, "/spiders", NewControllerV2[models.SpiderV2](
		Action{
			Method:      http.MethodGet,
			Path:        "/:id",
			HandlerFunc: GetSpiderById,
		},
		Action{
			Method:      http.MethodGet,
			Path:        "",
			HandlerFunc: GetSpiderList,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "",
			HandlerFunc: PostSpider,
		},
		Action{
			Method:      http.MethodPut,
			Path:        "/:id",
			HandlerFunc: PutSpiderById,
		},
		Action{
			Method:      http.MethodDelete,
			Path:        "/:id",
			HandlerFunc: DeleteSpiderById,
		},
		Action{
			Method:      http.MethodDelete,
			Path:        "",
			HandlerFunc: DeleteSpiderList,
		},
		Action{
			Method:      http.MethodGet,
			Path:        "/:id/files/list",
			HandlerFunc: GetSpiderListDir,
		},
		Action{
			Method:      http.MethodGet,
			Path:        "/:id/files/get",
			HandlerFunc: GetSpiderFile,
		},
		Action{
			Method:      http.MethodGet,
			Path:        "/:id/files/info",
			HandlerFunc: GetSpiderFileInfo,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/files/save",
			HandlerFunc: PostSpiderSaveFile,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/files/save/batch",
			HandlerFunc: PostSpiderSaveFiles,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/files/save/dir",
			HandlerFunc: PostSpiderSaveDir,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/files/rename",
			HandlerFunc: PostSpiderRenameFile,
		},
		Action{
			Method:      http.MethodDelete,
			Path:        "/:id/files",
			HandlerFunc: DeleteSpiderFile,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/files/copy",
			HandlerFunc: PostSpiderCopyFile,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/files/export",
			HandlerFunc: PostSpiderExport,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/run",
			HandlerFunc: PostSpiderRun,
		},
		Action{
			Method:      http.MethodGet,
			Path:        "/:id/git",
			HandlerFunc: GetSpiderGit,
		},
		Action{
			Method:      http.MethodGet,
			Path:        "/:id/git/remote-refs",
			HandlerFunc: GetSpiderGitRemoteRefs,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/git/checkout",
			HandlerFunc: PostSpiderGitCheckout,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/git/pull",
			HandlerFunc: PostSpiderGitPull,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/git/commit",
			HandlerFunc: PostSpiderGitCommit,
		},
		Action{
			Method:      http.MethodGet,
			Path:        "/:id/data-source",
			HandlerFunc: GetSpiderDataSource,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/data-source/:ds_id",
			HandlerFunc: PostSpiderDataSource,
		},
	))
	RegisterController(groups.AuthGroup, "/tasks", NewControllerV2[models.TaskV2](
		Action{
			Method:      http.MethodGet,
			Path:        "/:id",
			HandlerFunc: GetTaskById,
		},
		Action{
			Method:      http.MethodGet,
			Path:        "",
			HandlerFunc: GetTaskList,
		},
		Action{
			Method:      http.MethodDelete,
			Path:        "/:id",
			HandlerFunc: DeleteTaskById,
		},
		Action{
			Method:      http.MethodDelete,
			Path:        "",
			HandlerFunc: DeleteList,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/run",
			HandlerFunc: PostTaskRun,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/restart",
			HandlerFunc: PostTaskRestart,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/cancel",
			HandlerFunc: PostTaskCancel,
		},
		Action{
			Method:      http.MethodGet,
			Path:        "/:id/logs",
			HandlerFunc: GetTaskLogs,
		},
		Action{
			Method:      http.MethodGet,
			Path:        "/:id/data",
			HandlerFunc: GetTaskData,
		},
	))
	RegisterController(groups.AuthGroup, "/tokens", NewControllerV2[models.TokenV2](
		Action{
			Method:      http.MethodPost,
			Path:        "",
			HandlerFunc: PostToken,
		},
	))
	RegisterController(groups.AuthGroup, "/users", NewControllerV2[models.UserV2](
		Action{
			Method:      http.MethodPost,
			Path:        "",
			HandlerFunc: PostUser,
		},
		Action{
			Method:      http.MethodPost,
			Path:        "/:id/change-password",
			HandlerFunc: PostUserChangePassword,
		},
		Action{
			Method:      http.MethodGet,
			Path:        "/me",
			HandlerFunc: GetUserMe,
		},
		Action{
			Method:      http.MethodPut,
			Path:        "/me",
			HandlerFunc: PutUserById,
		},
	))

	RegisterActions(groups.AuthGroup, "/results", []Action{
		{
			Method:      http.MethodGet,
			Path:        "/:id",
			HandlerFunc: GetResultList,
		},
	})
	RegisterActions(groups.AuthGroup, "/export", []Action{
		{
			Method:      http.MethodPost,
			Path:        "/:type",
			HandlerFunc: PostExport,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:type/:id",
			HandlerFunc: GetExport,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:type/:id/download",
			HandlerFunc: GetExportDownload,
		},
	})
	RegisterActions(groups.AuthGroup, "/filters", []Action{
		{
			Method:      http.MethodGet,
			Path:        "/:col",
			HandlerFunc: GetFilterColFieldOptions,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:col/:value",
			HandlerFunc: GetFilterColFieldOptions,
		},
		{
			Method:      http.MethodGet,
			Path:        "/:col/:value/:label",
			HandlerFunc: GetFilterColFieldOptions,
		},
	})
	RegisterActions(groups.AuthGroup, "/settings", []Action{
		{
			Method:      http.MethodGet,
			Path:        "/:id",
			HandlerFunc: GetSetting,
		},
		{
			Method:      http.MethodPut,
			Path:        "/:id",
			HandlerFunc: PutSetting,
		},
	})
	RegisterActions(groups.AuthGroup, "/stats", []Action{
		{
			Method:      http.MethodGet,
			Path:        "/overview",
			HandlerFunc: GetStatsOverview,
		},
		{
			Method:      http.MethodGet,
			Path:        "/daily",
			HandlerFunc: GetStatsDaily,
		},
		{
			Method:      http.MethodGet,
			Path:        "/tasks",
			HandlerFunc: GetStatsTasks,
		},
	})

	RegisterActions(groups.AnonymousGroup, "/system-info", []Action{
		{
			Path:        "",
			Method:      http.MethodGet,
			HandlerFunc: GetSystemInfo,
		},
	})
	RegisterActions(groups.AnonymousGroup, "/version", []Action{
		{
			Method:      http.MethodGet,
			Path:        "",
			HandlerFunc: GetVersion,
		},
	})
	RegisterActions(groups.AnonymousGroup, "/", []Action{
		{
			Method:      http.MethodPost,
			Path:        "/login",
			HandlerFunc: PostLogin,
		},
		{
			Method:      http.MethodPost,
			Path:        "/logout",
			HandlerFunc: PostLogout,
		},
	})

	return nil
}

package controllers

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/service"
)

func InitControllers() (err error) {
	modelSvc, err := service.GetService()
	if err != nil {
		return err
	}

	NodeController = newNodeController()
	ProjectController = newProjectController()
	SpiderController = newSpiderController()
	TaskController = newTaskController()
	UserController = newUserController()
	TagController = NewListControllerDelegate(ControllerIdTag, modelSvc.GetBaseService(interfaces.ModelIdTag))
	SettingController = newSettingController()
	LoginController = NewActionControllerDelegate(ControllerIdLogin, getLoginActions())
	DataCollectionController = newDataCollectionController()
	ResultController = NewActionControllerDelegate(ControllerIdResult, getResultActions())
	ScheduleController = newScheduleController()
	StatsController = NewActionControllerDelegate(ControllerIdStats, getStatsActions())
	TokenController = newTokenController()
	FilerController = NewActionControllerDelegate(ControllerIdFiler, getFilerActions())
	GitController = NewListControllerDelegate(ControllerIdGit, modelSvc.GetBaseService(interfaces.ModelIdGit))
	VersionController = NewActionControllerDelegate(ControllerIdVersion, getVersionActions())
	SystemInfoController = NewActionControllerDelegate(ControllerIdSystemInfo, getSystemInfoActions())
	DemoController = NewActionControllerDelegate(ControllerIdDemo, getDemoActions())
	RoleController = NewListControllerDelegate(ControllerIdRole, modelSvc.GetBaseService(interfaces.ModelIdRole))
	PermissionController = NewListControllerDelegate(ControllerIdPermission, modelSvc.GetBaseService(interfaces.ModelIdPermission))
	ExportController = NewActionControllerDelegate(ControllerIdExport, getExportActions())
	NotificationController = NewActionControllerDelegate(ControllerIdNotification, getNotificationActions())
	FilterController = NewActionControllerDelegate(ControllerIdFilter, getFilterActions())
	SyncController = NewActionControllerDelegate(ControllerIdSync, getSyncActions())
	DataSourceController = newDataSourceController()
	EnvironmentController = newEnvironmentController()

	return nil
}

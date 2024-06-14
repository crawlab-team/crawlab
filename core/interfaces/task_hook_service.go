package interfaces

type TaskHookService interface {
	PreActions(Task, Spider, FsServiceV2, TaskHandlerService) (err error)
	PostActions(Task, Spider, FsServiceV2, TaskHandlerService) (err error)
}

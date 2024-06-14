package interfaces

type ServiceCrudOptions struct {
	IsAbsolute         bool // whether the path is absolute
	OnlyFromWorkspace  bool // whether only sync from workspace
	NotSyncToWorkspace bool // whether not sync to workspace
}

type ServiceCrudOption func(o *ServiceCrudOptions)

func WithOnlyFromWorkspace() ServiceCrudOption {
	return func(o *ServiceCrudOptions) {
		o.OnlyFromWorkspace = true
	}
}

func WithNotSyncToWorkspace() ServiceCrudOption {
	return func(o *ServiceCrudOptions) {
		o.NotSyncToWorkspace = true
	}
}

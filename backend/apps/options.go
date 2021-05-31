package apps

type MasterOption func(app MasterApp)

func WithConfigPath(path string) MasterOption {
	return func(app MasterApp) {
		app.SetConfigPath(path)
	}
}

func WithRunOnMaster(ok bool) MasterOption {
	return func(app MasterApp) {
		app.SetRunOnMaster(ok)
	}
}
